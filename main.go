package main

import (
	"wordshub/services/cli"
	"wordshub/services/conf"
	"wordshub/services/server"
)

func main() {
	env := cli.Parse()
	server.Start(conf.NewConfig(env))
}
