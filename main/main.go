package main

import (
	"github.com/jenpet/froxlor-dyndns/caller"
	"github.com/jenpet/froxlor-dyndns/config"
)

func main() {
	var coordinator caller.Coordinator
	coordinator.Init(config.ArgsConfig())
	coordinator.Run()
}
