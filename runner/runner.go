package runner

import (
	"fmt"

	"github.com/chasehampton/gom/handlers"
	"github.com/chasehampton/gom/models"
)

func Run(act models.Action, bh *handlers.BaseHandler) error {
	upload := act.IsUpload.Bool
	handler, err := getHandler(&act, bh)
	if err != nil {
		return err
	}
	if upload {
		err = handler.UploadFiles(act)
	} else {
		err = handler.DownloadFiles(act)
	}
	return err
}

func getHandler(act *models.Action, bh *handlers.BaseHandler) (handlers.Handler, error) {
	switch protocol := act.Connection.Protocol.ProtocolID; protocol {
	case 4:
		return &handlers.HttpHandler{BaseHandler: bh}, nil
	case 3:
		return handlers.GetS3HandlerFromAction(act, bh)
	default:
		return nil, fmt.Errorf("Unsupported protocol: %v", protocol)
	}
}
