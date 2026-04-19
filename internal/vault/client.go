package vault

import (
	"errors"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	api *vaultapi.Client
}

// Config holds explicit Vault connection settings.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a Vault client from explicit config or environment variables.
func NewClient(cfg *Config) (*Client, error) {
	address := cfg.Address
	token := cfg.Token

	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if address == "" {
		return nil, errors.New("vault address is required (set VAULT_ADDR or provide explicitly)")
	}
	if token == "" {
		return nil, errors.New("vault token is required (set VAULT_TOKEN or provide explicitly)")
	}

	apiCfg := vaultapi.DefaultConfig()
	apiCfg.Address = address

	c, err := vaultapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}
	c.SetToken(token)
	return &Client{api: c}, nil
}

// GetSecret retrieves a single field from a KV v2 secret path.
func (c *Client) GetSecret(path, field string) (string, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("reading vault path %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no data at vault path %q", path)
	}
	// KV v2 wraps data under "data" key.
	data := secret.Data
	if inner, ok := data["data"]; ok {
		if m, ok := inner.(map[string]interface{}); ok {
			data = m
		}
	}
	val, ok := data[field]
	if !ok {
		return "", fmt.Errorf("field %q not found at path %q", field, path)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("field %q at path %q is not a string", field, path)
	}
	return str, nil
}
