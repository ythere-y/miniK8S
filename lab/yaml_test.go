package lab

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type HTTP struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type MQ struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Debug struct {
	User []string `yaml:"user"`
	MQTT MQ       `yaml:"mqtt"`
	Http HTTP     `yaml:"http"`
}

type Config struct {
	Debug Debug `yaml:"config"`
}

func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	fmt.Println("conf: ", conf)
	return conf, nil
}

func main() {
	conf, err := ReadYamlConfig("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	byts, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("string:", string(byts))
}
