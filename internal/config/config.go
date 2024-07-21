package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config структура для конфигурации приложения
type Config struct {
	Host            string    `toml:"host"`
	Port            int       `toml:"port"`
	Database        Database  `toml:"database"`
	ServerResources Resources `toml:"serverResources"`
}

// Database структура для конфигурации базы данных
type Database struct {
	DatabaseName string `toml:"databaseName"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	SSLMode      string `toml:"sslMode"`
}

// Resources структура для конфигурации ресурсов сервера
type Resources struct {
	CPU  uint   `toml:"cpu"`
	RAM  uint16 `toml:"ram"`
	Disk uint   `toml:"disk"`
}

// LoadConfig загружает конфигурацию из файла config.toml
func LoadConfig() (config Config, err error) {
	viper.SetConfigFile("./config.toml")
	//viper.AutomaticEnv()
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// GetDatabaseDSN возвращает строку подключения к базе данных
func (c *Config) GetDatabaseDSN() string {
	fmt.Println(c.Database.DatabaseName)
	fmt.Println(c.Database.SSLMode)
	fmt.Println(c.Database.Host)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.Username, c.Database.Password, c.Database.DatabaseName, c.Database.SSLMode)
}
