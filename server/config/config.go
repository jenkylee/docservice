package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type Configuration interface {
	GetBoolean(path string, defaultVal ...bool) bool
	GetString(path string, defaultVal ...string) string
	GetTimeDuration(path string, defaultVal ...time.Duration) time.Duration
	String() string
}


type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}
	
	// 初始化配置文件
	if err := c.initConfig(); err != nil{
		return err
	}

	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// 如果指定配置文件，则解析指定配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}

	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")
	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// 监听配置文件是否修改， 用于热更新
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})
}

func (c *Config) GetString(key string, defaultVal ...string) string {
	v := viper.GetString(key)
	if v == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}

	return v
}

func (c *Config) GetBoolean(key string, defaultVal ...bool) bool {
	return viper.GetBool(key)
}

func (c *Config) GetTimeDuration(key string, defaultVal ...time.Duration) time.Duration {
	v := viper.GetDuration(key)
	if v == 0 {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return v * time.Second
}


