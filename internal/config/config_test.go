package config_test

import (
	"context"
	"errors"
	"testing"

	"github.com/USA-RedDragon/rtz-server/cmd"
	"github.com/USA-RedDragon/rtz-server/internal/config"
)

//nolint:golint,gochecknoglobals
var requiredFlags = []string{
	"--jwt.secret", "changeme",
	"--http.backend_url", "http://localhost:8081",
	"--mapbox.secret_token", "dummy",
	"--mapbox.public_token", "dummy",
}

func TestExampleConfig(t *testing.T) {
	t.Parallel()
	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	err := cmd.ParseFlags([]string{"--config", "../../config.example.yaml"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TesMissingOLTPEndpoint(t *testing.T) {
	t.Parallel()

	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	err := cmd.ParseFlags(append([]string{"--http.tracing.enabled", "true"}, requiredFlags...))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); !errors.Is(err, config.ErrOTLPEndpointRequired) {
		t.Errorf("unexpected error: %v", err)
	}

	err = cmd.ParseFlags(append([]string{"--http.tracing.enabled", "true", "--http.tracing.otlp_endpoint", "dummy"}, requiredFlags...))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err = config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMissingJWTSecret(t *testing.T) {
	t.Parallel()
	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	err := cmd.ParseFlags([]string{
		"--http.backend_url", "http://localhost:8081",
		"--mapbox.secret_token", "dummy",
		"--mapbox.public_token", "dummy",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); !errors.Is(err, config.ErrJWTSecretRequired) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMissingMapboxTokens(t *testing.T) {
	t.Parallel()
	baseCmd := cmd.NewCommand("testing", "deadbeef")
	baseCmd.SetContext(context.Background())
	baseFlags := []string{"--jwt.secret", "changeme", "--http.backend_url", "http://localhost:8081"}
	err := baseCmd.ParseFlags(append(baseFlags, []string{"--mapbox.secret_token", "dummy"}...))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err := config.LoadConfig(baseCmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); !errors.Is(err, config.ErrMapboxPublicTokenRequired) {
		t.Errorf("unexpected error: %v", err)
	}
	baseCmd = cmd.NewCommand("testing", "deadbeef")
	baseCmd.SetContext(context.Background())
	err = baseCmd.ParseFlags(append(baseFlags, []string{"--mapbox.public_token", "dummy"}...))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err = config.LoadConfig(baseCmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); !errors.Is(err, config.ErrMapboxSecretTokenRequired) {
		t.Errorf("unexpected error: %v", err)
	}
}

// Parallel tests are not allowed with t.Setenv
//
//nolint:golint,paralleltest
func TestEnvConfig(t *testing.T) {
	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	t.Setenv("HTTP__PORT", "8087")
	t.Setenv("HTTP__METRICS__PORT", "8088")
	t.Setenv("HTTP__METRICS__IPV4_HOST", "0.0.0.0")
	config, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if config.HTTP.Port != 8087 {
		t.Errorf("unexpected HTTP port: %d", config.HTTP.Port)
	}
	if config.HTTP.Metrics.Port != 8088 {
		t.Errorf("unexpected HTTP metrics port: %d", config.HTTP.Metrics.Port)
	}
	if config.HTTP.Metrics.IPV4Host != "0.0.0.0" {
		t.Errorf("unexpected HTTP metrics IPv4 host: %s", config.HTTP.Metrics.IPV4Host)
	}
}
