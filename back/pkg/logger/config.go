package logger

type Config struct {
	Level string `json:"loglevel" yaml:"loglevel" mapstructure:"loglevel"`
}
