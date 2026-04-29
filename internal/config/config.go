package config

import (
	"fmt"

	"apigateway/internal/logger"
	"apigateway/internal/model"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var GlobalConfig model.Config

func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Watch config change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		if err := viper.Unmarshal(&GlobalConfig); err != nil {
			fmt.Printf("failed to unmarshal config on change: %v\n", err)
			return
		}
		updateDynamicSettings()
	})

	return nil
}

func updateDynamicSettings() {
	// Dynamically adjust log level
	if logger.Log != nil {
		level, err := zapcore.ParseLevel(GlobalConfig.Log.Level)
		if err == nil {
			if logger.AtomicLevel.Level() != level {
				logger.AtomicLevel.SetLevel(level)
				logger.Log.Info("Log level updated", zap.String("new_level", level.String()))
			}
		}
	}
}
