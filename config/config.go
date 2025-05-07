package config

import "github.com/spf13/viper"

type Config struct {
}

func NewConfig(configPath string) *Config {

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{}
}

func (c *Config) RootPath() string {
	rootPath := viper.GetString("core.root-path")
	if rootPath != "" {
		return rootPath
	}
	return "."
}

func (c *Config) ServicePort() string {
	port := viper.GetString("core.port")
	if port == "" {
		port = "9999"
	}
	return port
}

func (c *Config) StoragePath() string {
	path := viper.GetString("core.storage-path")
	if path == "" {
		path = "."
	}
	return path
}

func (c *Config) LogFilePath() string {
	path := viper.GetString("core.log-file-path")
	if path == "" {
		path = "./request.log"
	}
	return path
}
