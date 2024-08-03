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
	Certificates struct {
		Email  string `yaml:"email" env:"CERTIFICATES_EMAIL"`
		Server string `yaml:"server" env:"CERTIFICATES_SERVER"`
	} `yaml:"certificates"`
}

func Init() {
	var err error
	Config, err = config.GetConfig[Conf](config.Environment(os.Getenv("ENV")))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

}
