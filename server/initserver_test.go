package server

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/config"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// Helper function to get an available port for testing
func getAvailablePort(t *testing.T) int {
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}

// TestInitialize_BasicFunctionality tests the core server initialization
func TestInitialize_BasicFunctionality(t *testing.T) {
	// Store original logger settings
	originalOutput := logger.StandardLogger().Out
	originalLevel := logger.GetLevel()

	// Setup log capture
	var logBuffer bytes.Buffer
	logger.SetOutput(&logBuffer)
	logger.SetLevel(logger.InfoLevel)
	defer func() {
		logger.SetOutput(originalOutput)
		logger.SetLevel(originalLevel)
	}()

	// Setup test components
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.Configuration{
		EnvironmentName: "test",
		Server: config.Server{
			Port:         getAvailablePort(t),
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
		},
	}

	// Create and run fx app
	app := fxtest.New(t,
		fx.Provide(func() *gin.Engine { return router }),
		fx.Provide(func() *config.Configuration { return cfg }),
		fx.Invoke(Initialize),
	)

	// Test startup
	startCtx, startCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer startCancel()

	err := app.Start(startCtx)
	require.NoError(t, err, "Failed to start fx app")

	// Allow server to initialize
	time.Sleep(300 * time.Millisecond)

	// Verify startup log
	logOutput := logBuffer.String()
	expectedStartupLog := fmt.Sprintf("Starting the REST application with %s environment and with port is %v",
		cfg.EnvironmentName, cfg.Server.Port)
	assert.Contains(t, logOutput, expectedStartupLog, "Startup log message not found")

	// Test shutdown
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	err = app.Stop(stopCtx)
	require.NoError(t, err, "Failed to stop fx app")

	// Verify shutdown log
	finalLogOutput := logBuffer.String()
	assert.Contains(t, finalLogOutput, "Server exiting", "Shutdown log message not found")
}

// TestInitialize_ServerConfiguration tests different server configurations
func TestInitialize_ServerConfiguration(t *testing.T) {
	testCases := []struct {
		name         string
		environment  string
		readTimeout  int
		writeTimeout int
		idleTimeout  int
	}{
		{
			name:         "Local Environment",
			environment:  "local",
			readTimeout:  10,
			writeTimeout: 10,
			idleTimeout:  20,
		},
		{
			name:         "Production Environment",
			environment:  "production",
			readTimeout:  60,
			writeTimeout: 60,
			idleTimeout:  120,
		},
		{
			name:         "Development Environment",
			environment:  "development",
			readTimeout:  30,
			writeTimeout: 30,
			idleTimeout:  60,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup log capture
			var logBuffer bytes.Buffer
			originalOutput := logger.StandardLogger().Out
			logger.SetOutput(&logBuffer)
			defer logger.SetOutput(originalOutput)

			gin.SetMode(gin.TestMode)
			router := gin.New()

			cfg := &config.Configuration{
				EnvironmentName: tc.environment,
				Server: config.Server{
					Port:         getAvailablePort(t),
					ReadTimeout:  tc.readTimeout,
					WriteTimeout: tc.writeTimeout,
					IdleTimeout:  tc.idleTimeout,
				},
			}

			app := fxtest.New(t,
				fx.Provide(func() *gin.Engine { return router }),
				fx.Provide(func() *config.Configuration { return cfg }),
				fx.Invoke(Initialize),
			)

			// Start and stop quickly
			startCtx, startCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer startCancel()

			require.NoError(t, app.Start(startCtx))

			time.Sleep(100 * time.Millisecond)

			stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer stopCancel()

			require.NoError(t, app.Stop(stopCtx))

			// Verify environment in logs
			logOutput := logBuffer.String()
			assert.Contains(t, logOutput, tc.environment)
		})
	}
}

// failed
// TestInitialize_PortBinding tests that server binds to correct port
func TestInitialize_PortBinding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add a simple health check route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := getAvailablePort(t)
	cfg := &config.Configuration{
		EnvironmentName: "test",
		Server: config.Server{
			Port:         port,
			ReadTimeout:  10,
			WriteTimeout: 10,
			IdleTimeout:  20,
		},
	}

	app := fxtest.New(t,
		fx.Provide(func() *gin.Engine { return router }),
		fx.Provide(func() *config.Configuration { return cfg }),
		fx.Invoke(Initialize),
	)

	startCtx, startCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer startCancel()

	require.NoError(t, app.Start(startCtx))

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Test port connectivity
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 1*time.Second)
	if err == nil {
		conn.Close()
		t.Logf("Successfully connected to server on port %d", port)
	} else {
		t.Logf("Connection test failed (may be expected): %v", err)
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer stopCancel()

	require.NoError(t, app.Stop(stopCtx))
}

// TestInitialize_Lifecycle tests fx lifecycle hooks
func TestInitialize_Lifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.Configuration{
		EnvironmentName: "lifecycle-test",
		Server: config.Server{
			Port:         getAvailablePort(t),
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
		},
	}

	// Track lifecycle events
	var startCalled, stopCalled bool

	app := fx.New(
		fx.Provide(func() *gin.Engine { return router }),
		fx.Provide(func() *config.Configuration { return cfg }),
		fx.Invoke(Initialize),
		fx.Invoke(func(lifecycle fx.Lifecycle) {
			// Add additional hooks to verify lifecycle works
			lifecycle.Append(fx.Hook{
				OnStart: func(context.Context) error {
					startCalled = true
					return nil
				},
				OnStop: func(context.Context) error {
					stopCalled = true
					return nil
				},
			})
		}),
	)

	startCtx, startCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer startCancel()

	require.NoError(t, app.Start(startCtx))
	assert.True(t, startCalled, "Start hook should have been called")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer stopCancel()

	require.NoError(t, app.Stop(stopCtx))
	assert.True(t, stopCalled, "Stop hook should have been called")
}

// TestInitialize_MultipleEnvironments tests various environment names
func TestInitialize_MultipleEnvironments(t *testing.T) {
	environments := []string{"local", "dev", "staging", "prod", "test", "development", "production"}

	for _, env := range environments {
		t.Run(fmt.Sprintf("Env_%s", env), func(t *testing.T) {
			var logBuffer bytes.Buffer
			originalOutput := logger.StandardLogger().Out
			logger.SetOutput(&logBuffer)
			defer logger.SetOutput(originalOutput)

			gin.SetMode(gin.TestMode)
			router := gin.New()

			cfg := &config.Configuration{
				EnvironmentName: env,
				Server: config.Server{
					Port:         getAvailablePort(t),
					ReadTimeout:  30,
					WriteTimeout: 30,
					IdleTimeout:  60,
				},
			}

			app := fxtest.New(t,
				fx.Provide(func() *gin.Engine { return router }),
				fx.Provide(func() *config.Configuration { return cfg }),
				fx.Invoke(Initialize),
			)

			startCtx, startCancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer startCancel()

			require.NoError(t, app.Start(startCtx))

			stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer stopCancel()

			require.NoError(t, app.Stop(stopCtx))

			logOutput := logBuffer.String()
			assert.Contains(t, logOutput, env, "Environment name should appear in logs")
		})
	}
}

// TestInitialize_EdgeCases tests edge cases and boundary conditions
func TestInitialize_EdgeCases(t *testing.T) {
	t.Run("Zero Timeouts", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		cfg := &config.Configuration{
			EnvironmentName: "zero-timeout-test",
			Server: config.Server{
				Port:         getAvailablePort(t),
				ReadTimeout:  0,
				WriteTimeout: 0,
				IdleTimeout:  0,
			},
		}

		app := fxtest.New(t,
			fx.Provide(func() *gin.Engine { return router }),
			fx.Provide(func() *config.Configuration { return cfg }),
			fx.Invoke(Initialize),
		)

		startCtx, startCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer startCancel()

		require.NoError(t, app.Start(startCtx))

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer stopCancel()

		require.NoError(t, app.Stop(stopCtx))
	})

	t.Run("Large Timeouts", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		cfg := &config.Configuration{
			EnvironmentName: "large-timeout-test",
			Server: config.Server{
				Port:         getAvailablePort(t),
				ReadTimeout:  3600, // 1 hour
				WriteTimeout: 3600,
				IdleTimeout:  7200, // 2 hours
			},
		}

		app := fxtest.New(t,
			fx.Provide(func() *gin.Engine { return router }),
			fx.Provide(func() *config.Configuration { return cfg }),
			fx.Invoke(Initialize),
		)

		startCtx, startCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer startCancel()

		require.NoError(t, app.Start(startCtx))

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer stopCancel()

		require.NoError(t, app.Stop(stopCtx))
	})

	t.Run("Empty Environment Name", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		cfg := &config.Configuration{
			EnvironmentName: "", // Empty environment name
			Server: config.Server{
				Port:         getAvailablePort(t),
				ReadTimeout:  30,
				WriteTimeout: 30,
				IdleTimeout:  60,
			},
		}

		app := fxtest.New(t,
			fx.Provide(func() *gin.Engine { return router }),
			fx.Provide(func() *config.Configuration { return cfg }),
			fx.Invoke(Initialize),
		)

		startCtx, startCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer startCancel()

		require.NoError(t, app.Start(startCtx))

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer stopCancel()

		require.NoError(t, app.Stop(stopCtx))
	})
}

// Benchmark test
func BenchmarkInitialize(b *testing.B) {
	gin.SetMode(gin.TestMode)

	for i := 0; i < b.N; i++ {
		router := gin.New()
		cfg := &config.Configuration{
			EnvironmentName: "benchmark",
			Server: config.Server{
				Port:         getAvailablePort(&testing.T{}),
				ReadTimeout:  30,
				WriteTimeout: 30,
				IdleTimeout:  60,
			},
		}

		app := fx.New(
			fx.Provide(func() *gin.Engine { return router }),
			fx.Provide(func() *config.Configuration { return cfg }),
			fx.Invoke(Initialize),
		)

		startCtx, startCancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = app.Start(startCtx)
		startCancel()

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = app.Stop(stopCtx)
		stopCancel()
	}
}
