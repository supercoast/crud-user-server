package helper

import (
	"log"

	"github.com/spf13/viper"
)

type ConfigManager struct {
	ConfigFileName string
	ConfigFilePath string
}

func NewConfigManager(configFileName, configFilePath string) *ConfigManager {
	return &ConfigManager{
		ConfigFileName: configFileName,
		ConfigFilePath: configFilePath,
	}
}

func LogError(err error) error {
	if err != nil {
		log.Println(err)
	}
	return err
}

func (c *ConfigManager) Init() {

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigName(c.ConfigFileName)
	viper.AddConfigPath(c.ConfigFilePath)
	viper.ReadInConfig()
}

func (c *ConfigManager) GetConfigValue(key string) string {
	return viper.GetString(key)

}

func (c *ConfigManager) GetEnvValue(key string) string {
	return viper.GetString(key)
}
