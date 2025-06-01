package config

import (
	"strings"
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
	// Если есть файл конфигурации — читаем
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Чтение из файла
	_ = viper.ReadInConfig() // игнорируем ошибку, если файла нет

	// Поддержка переменных окружения
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Значения по умолчанию
	if cfg.HTTP == nil {
		cfg.HTTP = &HTTP{
			Port:            "8081",
			ReadTimeout:     15,
			WriteTimeout:    15,
			ShutdownTimeout: 5,
		}
	}
	if cfg.JWT == nil {
		cfg.JWT = &JWT{
			AccessTokenExpiration:  15,
			RefreshTokenExpiration: 24 * 7,
			AccessTokenSecret:      "access_secret_key_change_in_production",
			RefreshTokenSecret:     "refresh_secret_key_change_in_production",
		}
	}
	return &cfg, nil
}
