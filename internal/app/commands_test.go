package app

import (
	"io"
	"os"
	"testing"

	"github.com/cainseing/drop-cli/internal/config"
)

func TestCreateGetCmd_RequiresArg(t *testing.T) {
	cmd := createGetCmd()
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when get command is called without a token")
	}
}

func TestCreatePurgeCmd_RequiresArg(t *testing.T) {
	cmd := createPurgeCmd()
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when purge command is called without a token")
	}
}

func TestIdentitySetCommand_SavesValidIdentity(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	cmd := CreateRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"identity", "set", "github", "alice"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected identity set command to succeed, got %v", err)
	}

	cfg := config.LoadUserConfig()
	if cfg.Username != "alice" {
		t.Fatalf("expected username alice, got %q", cfg.Username)
	}
	if cfg.Provider != "github" {
		t.Fatalf("expected provider github, got %q", cfg.Provider)
	}
}

func TestIdentitySetCommand_InvalidProvider(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	cmd := CreateRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"identity", "set", "bitbucket", "alex"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected invalid provider to return nil error, got %v", err)
	}

	cfgPath := config.GetConfigPath()
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Fatalf("expected config to not exist for unsupported provider, got %v", err)
	}
}

func TestCreateRootCmd_HasIdentitySubcommand(t *testing.T) {
	root := CreateRootCmd()

	if root == nil {
		t.Fatal("CreateRootCmd returned nil")
	}
	if root.Commands() == nil || len(root.Commands()) == 0 {
		t.Fatal("expected root command to have subcommands")
	}

	found, _, err := root.Find([]string{"identity"})
	if err != nil || found == nil {
		t.Fatalf("expected identity subcommand to exist, got err=%v", err)
	}
}

func TestCreateRootCmd_IdentitySetArgsValidation(t *testing.T) {
	cmd := CreateRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"identity", "set", "github"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when identity set command is called with too few args")
	}
}
