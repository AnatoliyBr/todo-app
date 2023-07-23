package app

import (
	"flag"
	"log"

	"github.com/AnatoliyBr/todo-app/internal/controller/apiserver"
	"github.com/BurntSushi/toml"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func Run() {

	// Controller
	flag.Parse()
	c := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, c)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.NewServer(c)
	if err := s.StartServer(); err != nil {
		log.Fatal(err)
	}
}
