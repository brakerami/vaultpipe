package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpipe/internal/audit"
	"vaultpipe/internal/config"
	"vaultpipe/internal/env"
	"vaultpipe/internal/process"
	"vaultpipe/internal/resolver"
	"vaultpipe/internal/vault"
)

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified: use -- <command> [args...]")
	}

	cfg, err := config.LoadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(vault.Config{
		Address: vaultAddr,
		Token:   vaultToken,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	logger := audit.NewLogger(os.Stderr)

	res := resolver.New(client, logger)
	refs := resolver.RefsFromConfig(cfg)

	secrets, err := res.Resolve(cmd.Context(), refs)
	if err != nil {
		return fmt.Errorf("resolve secrets: %w", err)
	}

	injected := env.BuildEnv(os.Environ(), secrets)

	runner, err := process.NewRunner(args)
	if err != nil {
		return fmt.Errorf("create runner: %w", err)
	}

	return runner.Run(cmd.Context(), injected)
}
