package resolver_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpipe/internal/resolver"
)

// mockClient satisfies the interface used by Resolver for testing.
type mockClient struct {
	secrets map[string]string
	err     error
}

func (m *mockClient) GetSecret(path, field string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	key := path + "#" + field
	val, ok := m.secrets[key]
	if !ok {
		return "", errors.New("secret not found: " + key)
	}
	return val, nil
}

func TestResolve_Success(t *testing.T) {
	refs := []resolver.SecretRef{
		{EnvKey: "DB_PASS", Path: "secret/db", Field: "password"},
	}
	// Integration path tested via real client; unit test skipped without mock injection.
	_ = refs
	t.Skip("requires mock injection refactor")
}

func TestResolve_Error(t *testing.T) {
	_ = &mockClient{err: errors.New("vault unavailable")}
	t.Skip("requires mock injection refactor")
}
