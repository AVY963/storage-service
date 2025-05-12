package config

import (
	"tages/internal/repository/pg"

	"github.com/spf13/viper"
)

type Config struct {
	PG   *pg.Config `mapstructure:"db"`
	HTTP *HTTP      `mapstructure:"http"`
	App  *App       `mapstructure:"app"`
	JWT  *JWT       `mapstructure:"jwt"`
}

type HTTP struct {
	Port            string `mapstructure:"port"`
	ReadTimeout     int    `mapstructure:"readTimeout"`
	WriteTimeout    int    `mapstructure:"writeTimeout"`
	ShutdownTimeout int    `mapstructure:"shutdownTimeout"`
}

type App struct {
	UploadLimiterConcurrency int    `mapstructure:"uploadLimiterConcurrency"`
	ListLimiterConcurrency   int    `mapstructure:"listLimiterConcurrency"`
	UploadDir                string `mapstructure:"uploadDir"`
}

type JWT struct {
	AccessTokenExpiration  int    `mapstructure:"accessTokenExpiration"`  // в минутах
	RefreshTokenExpiration int    `mapstructure:"refreshTokenExpiration"` // в часах
	AccessTokenSecret      string `mapstructure:"accessTokenSecret"`
	RefreshTokenSecret     string `mapstructure:"refreshTokenSecret"`
}

func InitConfig(path string) (*Config, error) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Значения по умолчанию для HTTP
	if cfg.HTTP == nil {
		cfg.HTTP = &HTTP{
			Port:            "8080",
			ReadTimeout:     15,
			WriteTimeout:    15,
			ShutdownTimeout: 5,
		}
	}

	// Значения по умолчанию для JWT
	if cfg.JWT == nil {
		cfg.JWT = &JWT{
			AccessTokenExpiration:  15,     // 15 минут по умолчанию
			RefreshTokenExpiration: 24 * 7, // 7 дней по умолчанию
			AccessTokenSecret:      "access_secret_key_change_in_production",
			RefreshTokenSecret:     "refresh_secret_key_change_in_production",
		}
	}

	return &cfg, nil
}
