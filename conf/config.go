package conf

import (
	"log"
	"os"

	"aidanwoods.dev/go-paseto"
	"github.com/risersh/config/config"
)

var Config *Conf
var PrivateKey paseto.V4AsymmetricSecretKey

type Conf struct {
	config.BaseConfig
	Port int `yaml:"port" env:"PORT" env-default:"8080"`
}

func Init() {
	var err error
	Config, err = config.GetConfig[Conf](config.Environment(os.Getenv("ENV")))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	PrivateKey, err = paseto.NewV4AsymmetricSecretKeyFromHex(Config.Sessions.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
}
