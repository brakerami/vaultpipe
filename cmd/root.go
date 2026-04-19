package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	vaultAddr  string
	vaultToken string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpipe [flags] -- <command> [args...]",
	Short: "Inject Vault secrets into a process environment",
	Long: `vaultpipe fetches secrets from HashiCorp Vault and injects them
as environment variables into a child process without writing to disk.`,
	Example: `  vaultpipe --config secrets.yaml -- env
  vaultpipe --vault-addr https://vault:8200 --config secrets.yaml -- ./app`,
	RunE: runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to secrets config file (required)")
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	_ = rootCmd.MarkPersistentFlagRequired("config")
}
