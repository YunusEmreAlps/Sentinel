package config

import (
	"os"
	"path/filepath"
	"sentinel/logger"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type config struct {
	App struct {
		Env       string `mapstructure:"env"`
		Version   string `mapstructure:"version"`
		Name      string `mapstructure:"name"`
		TargetApp string `mapstructure:"target_app"`
	} `mapstructure:"app"`

	DB struct {
		Type     string `mapstructure:"type"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"user"`
		Password string `mapstructure:"pass"`
		DBName   string `mapstructure:"db"`
		SSLMode  string `mapstructure:"ssl"`
	} `mapstructure:"db"`

	Mail struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		FromName string `mapstructure:"from_name"`
		FromMail string `mapstructure:"from_mail"`
	} `mapstructure:"mail"`
}

var C config

func ReadConfig(processCwdir string) {
	Config := &C
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(processCwdir, "config"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.CLogger.Error("INIT: Cannot read config file.")
	}

	if err := viper.Unmarshal(&Config); err != nil {
		logger.CLogger.Error("INIT: Cannot unmarshal config file.")
		os.Exit(1)
	}

	spew.Dump(C)
}
