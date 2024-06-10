package detail

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
)

func LoadConfig() (string, int) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		return "", -1
	}
	cookie := viper.GetString("Cookie")
	return cookie, 0
}
