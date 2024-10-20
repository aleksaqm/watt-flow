package server

import "watt-flow/config"

func Init() {
	serverConfig := config.GetConfig()
	r := NewRouter()
	r.Run(":" + serverConfig.GetString("server.port"))
}
