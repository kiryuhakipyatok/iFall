package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App         AppConfig         `mapstructure:"app"`
	Storage     StorageConfig     `mapstructure:"storage"`
	Server      ServerConfig      `mapstructure:"server"`
	Email       EmailConfig       `mapstructure:"email"`
	ApiClient   ApiClientConfig   `mapstructure:"apiClient"`
	Scheduler   SchedulerConfig   `mapstructure:"scheduler"`
	IPhones     IPhonesConfig     `mapstructure:"iphones"`
	TelegramBot TelegramBotConfig `mapstructure:"telegramBot"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Version string `mapstructure:"version"`
	LogPath string `mapstructure:"logPath"`
}

type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	RequestTimeout time.Duration `mapstructure:"requestTimeout"`
	CloseTimeout   time.Duration `mapstructure:"closeTimeout"`
}

type StorageConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	Database       string        `mapstructure:"database"`
	SSLMode        string        `mapstructure:"sslMode"`
	Timezone       string        `mapstructure:"timezone"`
	ConnectTimeout time.Duration `mapstructure:"connectTimeout"`
	PingTimeout    time.Duration `mapstructure:"pingTimeout"`
	AmountOfConns  int32         `mapstructure:"amountOfConns"`
}

type EmailConfig struct {
	Name              string `mapstructure:"name"`
	Password          string `mapstructure:"password"`
	Address           string `mapstructure:"address"`
	SmtpAddress       string `mapstructure:"smtpAddress"`
	SmtpServerAddress string `mapstructure:"smtpServerAddress"`
}

type ApiClientConfig struct {
	BaseURL string        `mapstructure:"baseURL"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type SchedulerConfig struct {
	Hour   int `mapstructure:"hour"`
	Minute int `mapstructure:"minute"`
}

type IPhonesConfig struct {
	Black   string        `mapstructure:"black"`
	White   string        `mapstructure:"white"`
	Green   string        `mapstructure:"green"`
	Pink    string        `mapstructure:"pink"`
	Blue    string        `mapstructure:"blue"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type TelegramBotConfig struct {
	Token   string        `mapstructure:"token"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func MustLoad(path string) *Config {
	if path == "" {
		panic("config path is empty")
	}
	filename := filepath.Join(path, "config.yaml")
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
	data = []byte(os.ExpandEnv(string(data)))
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	return cfg
}
