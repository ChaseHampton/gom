package vaultclient

import (
	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *vault.Client
}

func NewVaultClient(address, token string) (*VaultClient, error) {
	config := vault.DefaultConfig()
	config.Address = address
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return &VaultClient{client: client}, nil
}

func (vc *VaultClient) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := vc.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}
	return secret.Data, nil
}
