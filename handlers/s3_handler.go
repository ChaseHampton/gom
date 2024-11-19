package handlers

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/chasehampton/gom/models"
)

type S3Handler struct {
	b   *BaseHandler
	s3c *s3.Client
}

type S3Target interface {
	Target
	GetBucket() string
}

type DLFile struct {
	Bucket string
	Object s3types.Object
}

type S3TargetImpl struct {
	TargetPath string
	Bucket     string
}

func (st S3TargetImpl) GetTarget() string {
	return st.TargetPath
}

func (st S3TargetImpl) GetBucket() string {
	return st.Bucket
}

func GetS3HandlerFromAction(act *models.Action, bh *BaseHandler) (*S3Handler, error) {
	creds, err := GetS3Creds(act.Connection.AuthDetail.VaultPath.String, bh)
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3Handler{b: bh, s3c: client}, nil
}

func GetS3Creds(vaultPath string, bh *BaseHandler) (aws.CredentialsProvider, error) {
	secret, err := bh.Vault.ReadSecret(vaultPath)
	if err != nil {
		return nil, err
	}
	data, ok := secret["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to extract data field from secret")
	}
	keyId, ok := data["AWS_ACCESS_KEY_ID"].(string)
	if !ok || keyId == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID not found in secret")
	}
	secretKey, ok := data["AWS_SECRET_ACCESS_KEY"].(string)
	if !ok || secretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY not found in secret")
	}

	credProv := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(keyId, secretKey, ""))
	return credProv, nil
}

func (h *S3Handler) TargetFromAction(act models.Action, targetName string) (Target, error) {
	p := normalizePrefix(act.RemotePath.String, act.Bucket.String)
	key, err := filepath.Rel(act.LocalPath.String, targetName)
	if err != nil {
		return nil, err
	}
	t := S3TargetImpl{TargetPath: path.Join(p, key), Bucket: act.Bucket.String}
	return t, nil
}

func (h *S3Handler) UploadFiles(act models.Action) error {
	files, err := h.b.GetUploadFiles(act.LocalPath.String)
	if err != nil {
		return err
	}
	for _, file := range files {
		// upload file
		reads, err := h.b.OpenFile(file)
		if err != nil {
			h.b.Logger.LogError(err, &act, "")
		}
		defer reads.Close()
		target, err := h.TargetFromAction(act, file)
		if err != nil {
			h.b.Logger.LogError(err, &act, "")
		}
		err = uploadFile(reads, target.(S3Target), h.s3c)
		if err != nil {
			h.b.Logger.LogError(err, &act, "")
		}
	}
	return nil
}

func (h *S3Handler) DownloadFiles(act models.Action) error {
	files, err := h.ListFiles(act)
	if err != nil {
		return err
	}
	for _, file := range files {
		// download file
		f, err := h.getDownloadFile(*file.(DLFile).Object.Key, &act)
		defer f.Close()
		if err != nil {
			h.b.Logger.LogError(err, &act, "")
		}
		if f == nil {
			continue
		}
		err = h.downloadFile(file.(DLFile), f)
	}
	return nil
}

func (h *S3Handler) DeleteFile() error {
	return nil
}

func (h *S3Handler) ListFiles(act models.Action) ([]interface{}, error) {
	fileList := []interface{}{}
	listRequest := &s3.ListObjectsV2Input{
		Bucket: aws.String(act.Bucket.String),
		Prefix: aws.String(normalizePrefix(act.RemotePath.String, act.Bucket.String)),
	}
	result, err := h.s3c.ListObjectsV2(context.TODO(), listRequest)
	if err != nil {
		return nil, err
	}
	for _, obj := range result.Contents {
		dfile := DLFile{Bucket: act.Bucket.String, Object: obj}
		fileList = append(fileList, dfile)
	}
	return fileList, nil
}

func normalizePrefix(path string, bucket string) string {
	normPath := strings.TrimPrefix(path, "/")
	if strings.HasPrefix(normPath, bucket+"/") {
		normPath = strings.TrimPrefix(normPath, bucket+"/")
	}
	return normPath
}

func uploadFile(file *os.File, target S3Target, s3c *s3.Client) error {
	_, err := s3c.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(target.(S3Target).GetBucket()),
		Key:    aws.String(target.(S3Target).GetTarget()),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *S3Handler) downloadFile(file DLFile, downloadFile *os.File) error {
	downloadRequest := &s3.GetObjectInput{
		Bucket: aws.String(file.Bucket),
		Key:    aws.String(*file.Object.Key),
	}
	result, err := h.s3c.GetObject(context.TODO(), downloadRequest)
	if err != nil {
		return err
	}
	defer result.Body.Close()
	_, err = downloadFile.ReadFrom(result.Body)
	if err != nil {
		return err
	}
	return nil
}

func (h *S3Handler) getDownloadFile(key string, act *models.Action) (*os.File, error) {
	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}
	if strings.HasPrefix(key, act.Bucket.String+"/") {
		key = strings.TrimPrefix(key, act.Bucket.String+"/")
	}
	rel, err := filepath.Rel(act.RemotePath.String, key)
	if err != nil {
		return nil, err
	}
	if rel == "." {
		return nil, nil
	}
	dfile := path.Join(act.LocalPath.String, rel)
	return h.b.OpenFile(dfile)
}
