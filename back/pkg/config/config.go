package config

import (
	"errors"
	"flag"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"strings"
	"time"
)

const defaultConfigPath = "./config/config.yml"

var confPath string

func ViperInit() error {
	confPath = *flag.String("config", defaultConfigPath, "Path to the configuration file")
	flag.Parse()
	if strings.TrimSpace(confPath) == "" {
		log.Printf("use default config path: %s", defaultConfigPath)
		confPath = defaultConfigPath
	}
	viper.SetConfigType("yml")
	viper.SetConfigFile(confPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			GetFromEnv()
			return nil
		}
		return err
	}
	return nil
}

func GetFromEnv() {
	viper.BindEnv("app.addr", "APP_ADDR")
	viper.BindEnv("app.timeout", "APP_TIMEOUT")
	viper.BindEnv("app.log_level", "APP_LOG_LEVEL")

	viper.BindEnv("db.host", "POSTGRES_HOST")
	viper.BindEnv("db.port", "POSTGRES_PORT")
	viper.BindEnv("db.user", "POSTGRES_USER")
	viper.BindEnv("db.password", "POSTGRES_PASSWORD")
	viper.BindEnv("db.database", "POSTGRES_DB")
	viper.BindEnv("db.ssl", "POSTGRES_SSL")
	viper.BindEnv("db.debug", "POSTGRES_DEBUG")
	viper.BindEnv("db.max_idle", "POSTGRES_MAX_IDLE")
	viper.BindEnv("db.max_open", "POSTGRES_MAX_OPEN")
}

func ParseConfig() (*Config, error) {
	c := new(Config)
	viper.AllKeys()
	err := viper.UnmarshalExact(
		c,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.RecursiveStructToMapHookFunc(),
				mapstructure.StringToTimeDurationHookFunc(),
			),
		),
	)

	return c, err
}

type Config struct {
	App App `mapstructure:"app"`
	DB  DB  `mapstructure:"db"`
}

type App struct {
	Addr    string        `mapstructure:"addr"`
	Timeout time.Duration `mapstructure:"timeout"`
	Log     string        `mapstructure:"log_level"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSL      string `mapstructure:"ssl"`
	Debug    bool   `mapstructure:"debug"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
}
