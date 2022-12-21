package config

import (
	"bytes"
	"encoding/json"
	"git.xantus.network/apt-gemxyz/pkg/database"
	"git.xantus.network/apt-gemxyz/pkg/log"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Log              log.Config           `json:"log" mapstructure:"log"`
	MySQL            database.MySQLConfig `json:"mysql" mapstructure:"mysql"`
	MigrationsFolder string               `json:"migrations_folder" mapstructure:"migrations_folder"`
}

func Load() (*Config, error) {

	// You should set default config value here
	c := &Config{
		MySQL: database.MySQLConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			Database: "apt_gemxyz",
			Username: "root",
			Password: "root",
			Debug:    true,
			Options:  "?parseTime=true",
		},
		Log: log.Config{
			Level:  "debug",
			Format: "text",
		},
		MigrationsFolder: "file://migrations",
	}

	// --- hacking to load reflect structure config into env ----//
	viper.SetConfigType("json")
	configBuffer, err := json.Marshal(c)

	if err != nil {
		return nil, err
	}

	if err := viper.ReadConfig(bytes.NewBuffer(configBuffer)); err != nil {
		panic(err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// -- end of hacking --//
	viper.AutomaticEnv()
	err = viper.Unmarshal(c)
	return c, err
}
