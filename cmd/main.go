package main

import (
	"Test_derictory/configs"
	"Test_derictory/server"
	"github.com/spf13/viper"
	"log"
)

func main() {

	if err := configs.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error)
	}
}
