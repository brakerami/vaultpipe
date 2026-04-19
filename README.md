# vaultpipe

> CLI tool to inject secrets from Vault into process environments without writing to disk

---

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

---

## Usage

`vaultpipe` fetches secrets from HashiCorp Vault and injects them as environment variables into a subprocess — no temp files, no disk writes.

```bash
vaultpipe run --path secret/myapp -- ./myapp serve
```

The secrets stored at `secret/myapp` are injected directly into the environment of `./myapp serve`.

### Options

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path to read from |
| `--addr` | Vault server address (default: `$VAULT_ADDR`) |
| `--token` | Vault token (default: `$VAULT_TOKEN`) |
| `--prefix` | Optional prefix for injected env var names |

### Example

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=s.mytoken

vaultpipe run --path secret/data/prod/db -- python app.py
```

Secrets like `username` and `password` from Vault become `USERNAME` and `PASSWORD` in the child process environment.

---

## How It Works

1. Authenticates with Vault using the provided token or environment credentials
2. Reads the specified secret path
3. Exports key-value pairs as environment variables
4. Executes the given command with those variables injected — never touching disk

---

## License

MIT © [yourusername](https://github.com/yourusername)