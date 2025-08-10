package tp_logger

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.SugaredLogger

// Config holds logger configuration
type Config struct {
	ServiceName string // Required: "scraper", "core", "matcher"
	LogFile     string // Required: "/app/logs/file.log"
	Environment string // Optional: defaults to "dev"
	Version     string // Optional: defaults to "1.0.0"
	Console     bool   // Optional: enable console output - defaults to true
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
	var outputs []string
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

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	globalLogger = logger.Sugar()
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
	globalLogger.Info(args...)
}

func Printf(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Println(args ...interface{}) {
	globalLogger.Info(args...)
}

// Fatal functions
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Panic functions
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	globalLogger.Panic(args...)
}

// Error functions
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Warn functions
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Info functions
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Debug functions
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Structured logging functions
func InfoStruct(msg string, keysAndValues ...interface{}) {
	globalLogger.Infow(msg, keysAndValues...)
}

func ErrorStruct(msg string, keysAndValues ...interface{}) {
	globalLogger.Errorw(msg, keysAndValues...)
}

func DebugStruct(msg string, keysAndValues ...interface{}) {
	globalLogger.Debugw(msg, keysAndValues...)
}

func WarnStruct(msg string, keysAndValues ...interface{}) {
	globalLogger.Warnw(msg, keysAndValues...)
}

// Context logging - creates logger with additional fields
func WithFields(keysAndValues ...interface{}) *zap.SugaredLogger {
	return globalLogger.With(keysAndValues...)
}

// Sync flushes any buffered log entries
func Sync() {
	if globalLogger != nil {
		globalLogger.Sync()
	}
}
