package config

import (
	"log"

	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Drive GDriveConfig
	Local FSConfig
}

type GDriveConfig struct {
	Root     string
	Year     string
	Semester string
}

type FSConfig struct {
	SyncFolder string
	CredPath   string
}

func Init() (conf *Configurations) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return
}
