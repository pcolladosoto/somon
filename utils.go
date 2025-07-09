//go:build !patch

package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func applyConfiguration() error {
	level := slog.LevelInfo
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		switch strings.ToLower(logLevel) {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: false})))

	return nil
}

func getEnvVar(envVar string) (string, error) {
	val := os.Getenv(envVar)
	if val == "" {
		return "", fmt.Errorf("%s environment variable is not defined", envVar)
	}
	return val, nil
}
