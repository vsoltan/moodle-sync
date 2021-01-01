package config

import (
	"log"

	"github.com/spf13/viper"
)

// Configurations contains user-specific information for sync service
type Configurations struct {
	Drive GDriveConfig
	Local FSConfig
}

// GDriveConfig defines a struct for relevant Google Drive details
type GDriveConfig struct {
	Folders []string
}

// FSConfig defines a struct for local filesystem requirements for sync service
type FSConfig struct {
	SyncFolder string
	CredPath   string
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
}

// Parse reads the specified config file and decodes it into a struct for future use
func Parse() (conf *Configurations) {
	log.Println("Parsing config file...")

	initConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	log.Println("Success!")
	return
}
