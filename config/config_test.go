package config

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/fx"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	full := filepath.Join(dir, name)
	if err := ioutil.WriteFile(full, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", full, err)
	}
}

func TestNewFxModule_LoadsConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	defaultYml := `
server:
  port: 8080
  readTimeout: 30
  writeTimeout: 40
  idleTimeout: 50
logger:
  base: logrus
  level: info
  format: text
  reportCaller: true
  enabled: true
  maxSize: 100
  maxAge: 7
  maxBackups: 3
  localTime: true
  compress: false
  logDir: /var/log/app
`
	writeFile(t, tmpDir, "server.yml", defaultYml)

	var cfg *Configuration
	app := fx.New(
		NewFxModule(tmpDir, ""),
		fx.Populate(&cfg),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("failed to start fx app: %v", err)
	}
	defer app.Stop(context.Background())

	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port=8080; got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 30 {
		t.Errorf("expected Server.ReadTimeout=30; got %d", cfg.Server.ReadTimeout)
	}
	if cfg.Logger.Format != "text" {
		t.Errorf("expected Logger.Format=\"text\"; got %q", cfg.Logger.Format)
	}
}

func TestNewFxModule_FallbackToEnv(t *testing.T) {
	// create temp dir and server.yml
	tmpDir, err := ioutil.TempDir("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fallbackYml := `
server:
  port: 8081
`
	writeFile(t, tmpDir, "server.yml", fallbackYml)

	os.Setenv("CONFIG_PATH", tmpDir)
	defer os.Unsetenv("CONFIG_PATH")

	var cfg *Configuration
	app := fx.New(
		NewFxModule("", ""),
		fx.Populate(&cfg),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("failed to start fx app: %v", err)
	}
	defer app.Stop(context.Background())

	if cfg.Server.Port != 8081 {
		t.Errorf("expected fallback Server.Port=8081; got %d", cfg.Server.Port)
	}
}

func TestNewFxModule_ErrorWhenMissingFile(t *testing.T) {
	var cfg *Configuration
	app := fx.New(
		NewFxModule("/does-not-exist", ""),
		fx.Populate(&cfg),
	)
	err := app.Start(context.Background())
	if err == nil {
		t.Fatal("expected error when config file is missing; got nil")
	}
}

func TestNewFxModule_OverrideYAML(t *testing.T) {
	// default + override in same temp dir
	tmpDir, err := ioutil.TempDir("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	defaultYml := `
server:
  port: 8080
`
	writeFile(t, tmpDir, "server.yml", defaultYml)

	overrideYml := `
server:
  port: 9090
`
	overridePath := filepath.Join(tmpDir, "override.yml")
	writeFile(t, tmpDir, "override.yml", overrideYml)

	var cfg *Configuration
	app := fx.New(
		NewFxModule(tmpDir, overridePath),
		fx.Populate(&cfg),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("failed to start fx app with override: %v", err)
	}
	defer app.Stop(context.Background())

	if cfg.Server.Port != 9090 {
		t.Errorf("expected overridden Server.Port=9090; got %d", cfg.Server.Port)
	}
}
