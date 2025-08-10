package logger

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.SugaredLogger
var initialized bool

// Config holds logger configuration
type Config struct {
	ServiceName string // Optional: defaults to SERVICE_NAME env var
	LogFile     string // Optional: defaults to /app/logs/{service}.log
	Environment string // Optional: defaults to APP_ENV or "dev"
	Version     string // Optional: defaults to APP_VERSION or "1.0.0"
	Console     bool   // Optional: enable console output - defaults to true
}

// ensureInitialized initializes logger with defaults if not already done
func ensureInitialized() {
	if initialized {
		return
	}

	// Auto-initialize with defaults
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "app" // fallback
	}

	logFile := fmt.Sprintf("/app/logs/%s.log", serviceName)

	config := Config{
		ServiceName: serviceName,
		LogFile:     logFile,
		Console:     false,
	}

	if err := Init(config); err != nil {
		// If init fails, create minimal console-only logger
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.OutputPaths = []string{"stdout"}
		logger, _ := zapConfig.Build()
		globalLogger = logger.Sugar()
	}

	initialized = true
}

// generateTraceID creates a unique trace ID for this session
func generateTraceID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("trace_%d_%d", time.Now().Unix(), r.Intn(10000))
}

// getHostname returns the container hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// Init initializes the global logger with provided configuration
func Init(cfg Config) error {
	if cfg.ServiceName == "" {
		return fmt.Errorf("ServiceName is required")
	}
	if cfg.LogFile == "" {
		return fmt.Errorf("LogFile is required")
	}

	// Set defaults
	if cfg.Environment == "" {
		cfg.Environment = os.Getenv("APP_ENV")
		if cfg.Environment == "" {
			cfg.Environment = "dev"
		}
	}
	if cfg.Version == "" {
		cfg.Version = os.Getenv("APP_VERSION")
		if cfg.Version == "" {
			cfg.Version = "1.0.0"
		}
	}

	config := zap.NewProductionConfig()

	// Configure output paths
	outputs := []string{}
	if cfg.Console {
		outputs = append(outputs, "stdout")
	}
	outputs = append(outputs, cfg.LogFile)
	config.OutputPaths = outputs

	// Configure encoder for readable logs
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"

	// Development mode for dev environment
	if cfg.Environment == "dev" {
		config.Development = true
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	// Add default fields to ALL logs
	config.InitialFields = map[string]interface{}{
		"service":  cfg.ServiceName,
		"env":      cfg.Environment,
		"version":  cfg.Version,
		"trace_id": generateTraceID(),
		"host":     getHostname(),
	}

	// Ensure log directory exists
	logDir := filepath.Dir(cfg.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	globalLogger = logger.Sugar()
	initialized = true
	return nil
}

// MustInit initializes logger and panics on error
func MustInit(cfg Config) {
	if err := Init(cfg); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

// Drop-in replacement functions for standard log package

// Print functions
func Print(args ...interface{}) {
	ensureInitialized()
	globalLogger.Info(args...)
}

func Printf(format string, args ...interface{}) {
	ensureInitialized()
	globalLogger.Infof(format, args...)
}

func Println(args ...interface{}) {
	ensureInitialized()

	globalLogger.Info(args...)
}

// Fatal functions
func Fatal(args ...interface{}) {
	ensureInitialized()

	globalLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	ensureInitialized()

	globalLogger.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	ensureInitialized()

	globalLogger.Fatal(args...)
}

// Panic functions
func Panic(args ...interface{}) {
	ensureInitialized()

	globalLogger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	ensureInitialized()

	globalLogger.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	ensureInitialized()

	globalLogger.Panic(args...)
}

// Error functions
func Error(args ...interface{}) {
	ensureInitialized()

	globalLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	ensureInitialized()

	globalLogger.Errorf(format, args...)
}

// Warn functions
func Warn(args ...interface{}) {
	ensureInitialized()

	globalLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	ensureInitialized()

	globalLogger.Warnf(format, args...)
}

// Info functions
func Info(args ...interface{}) {
	ensureInitialized()

	globalLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	ensureInitialized()

	globalLogger.Infof(format, args...)
}

// Debug functions
func Debug(args ...interface{}) {
	ensureInitialized()

	globalLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	ensureInitialized()

	globalLogger.Debugf(format, args...)
}

// Structured logging functions
func InfoStruct(msg string, keysAndValues ...interface{}) {
	ensureInitialized()
	globalLogger.Infow(msg, keysAndValues...)
}

func ErrorStruct(msg string, keysAndValues ...interface{}) {
	ensureInitialized()

	globalLogger.Errorw(msg, keysAndValues...)
}

func DebugStruct(msg string, keysAndValues ...interface{}) {
	ensureInitialized()

	globalLogger.Debugw(msg, keysAndValues...)
}

func WarnStruct(msg string, keysAndValues ...interface{}) {
	ensureInitialized()

	globalLogger.Warnw(msg, keysAndValues...)
}

// Context logging - creates logger with additional fields
func WithFields(keysAndValues ...interface{}) *zap.SugaredLogger {
	return globalLogger.With(keysAndValues...)
}
