package config

import (
	"os"
	"path/filepath"

	"sentinel/pkg/logger"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type Config struct {
	App     App     `mapstructure:"app"`
	Auth    Auth    `mapstructure:"auth"`
	DB      DB      `mapstructure:"db"`
	Cache   Cache   `mapstructure:"cache"`
	Broker  Broker  `mapstructure:"broker"`
	Cookie  Cookie  `mapstructure:"cookie"`
	Session Session `mapstructure:"session"`
	Metric  Metric  `mapstructure:"metric"`
	Jaeger  Jaeger  `mapstructure:"jaeger"`
	Mail    Mail    `mapstructure:"mail"`
}

type App struct {
	Mode    string `mapstructure:"mode"`
	Port    string `mapstructure:"port"`
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
	Expire  int    `mapstructure:"expire"`
}

type Auth struct {
	JwtKey string `mapstructure:"jwt_pub"`
}

type DB struct {
	Active bool   `mapstructure:"active"`
	Url    string `mapstructure:"url"`
}

type Cache struct {
	Active bool   `mapstructure:"active"`
	Url    string `mapstructure:"url"`
}

type Broker struct {
	Url           string `mapstructure:"url"`
	ConsumerGroup string `mapstructure:"consumer_group"`
	Topic         string `mapstructure:"topic"`
}

type Cookie struct {
	Name     string `mapstructure:"name"`
	MaxAge   int    `mapstructure:"max_age"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"http_only"`
}

type Session struct {
	Name   string `mapstructure:"name"`
	Prefix string `mapstructure:"prefix"`
	Expire int    `mapstructure:"expire"`
}

type Metric struct {
	Url     string `mapstructure:"url"`
	Service string `mapstructure:"service"`
}

type Jaeger struct {
	Host        string `mapstructure:"host"`
	ServiceName string `mapstructure:"service_name"`
	LogSpans    bool   `mapstructure:"log_spans"`
}

type Mail struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	FromName string `mapstructure:"from_name"`
	FromMail string `mapstructure:"from_mail"`
	To       string `mapstructure:"to"`
	CC       string `mapstructure:"cc"`
}

var C Config

func ReadConfig(processCwdir string) {
	Config := &C
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(processCwdir, "config"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Cannot read config file:", err)
		logger.CLogger.Info("Cannot read config file:", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		logger.CLogger.Info("Cannot unmarshal config file:", err)
		os.Exit(1)
	}

	spew.Dump(C)
}
