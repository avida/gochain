package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func ReadConf(path string) {
	viper.SetConfigName("db")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Println(viper.Get("postgres.Port"))
}

func DBConfig() (port int, user string, pass string) {
	port = viper.GetInt("postgres.Port")
	user = viper.GetString("postgres.User")
	pass = viper.GetString("postgres.Password")
	return
}
