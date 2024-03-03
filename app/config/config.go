package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

type Config struct {
	Server struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"redis"`
	Nats struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"nats"`
	Clickhouse struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"clickhouse"`
}

func NewConfig(path string) (*Config, error) {
	var conf Config
	err := configor.Load(&conf, path)
	return &conf, err
}

func (c *Config) DbConn() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=Asia/Tashkent",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name)
}

func (c *Config) ServerIP() string {
	return fmt.Sprintf("%s:%s", c.Server.IP, c.Server.Port)
}

func (c *Config) NastConn() string {
	return fmt.Sprintf("%s:%d", c.Nats.Host, c.Nats.Port)
}

func (c *Config) RedisConn() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) ClockhouseConn() string {
	return fmt.Sprintf("%s:%s", c.Clickhouse.Host, c.Clickhouse.Port)
}
