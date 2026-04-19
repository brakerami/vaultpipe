package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
// If Address or Token are empty, it falls back to environment variables.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address not set: provide --vault-addr or VAULT_ADDR")
	}
	vcfg.Address = address

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token not set: provide --vault-token or VAULT_TOKEN")
	}
	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// ReadSecrets reads a KV v2 secret at the given path and returns key/value pairs.
func (c *Client) ReadSecrets(mountPath, secretPath string) (map[string]string, error) {
	kv := c.vc.KVv2(mountPath)
	secret, err := kv.Get(nil, secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q from mount %q: %w", secretPath, mountPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at %q", secretPath)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("secret key %q has non-string value", k)
		}
		result[k] = str
	}
	return result, nil
}
