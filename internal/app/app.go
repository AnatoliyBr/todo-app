package app

import (
	"flag"
	"log"

	"github.com/AnatoliyBr/todo-app/internal/controller/apiserver"
	"github.com/AnatoliyBr/todo-app/internal/store"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func Run() {

	// PostgreSQL
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	configDB := store.NewConfig()
	_, err := store.NewDB(configDB)
	if err != nil {
		log.Fatal(err)
	}

	// Controller
	flag.Parse()
	configServer := apiserver.NewConfig()
	_, err = toml.DecodeFile(configPath, configServer)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.NewServer(configServer)
	if err := s.StartServer(); err != nil {
		log.Fatal(err)
	}
}
