package config

import (
	"io/ioutil"
	"log"
	"os"

	"fmt"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Queue    struct {
		Name             string `yaml:"name"`
		Vhost            string `yaml:"vhost"`
		Durable          bool   `yaml:"durable"`
		DeleteWhenUnused bool   `yaml:"delete_when_unused"`
		Exclusive        bool   `yaml:"exclusive"`
		NoWait           bool   `yaml:"no_wait"`
	} `yaml:"queue"`
}

// ReadFromConfig takes the path of a YAML config file and returns a Config struct
func ReadFromConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	c := Config{}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return c
}

func (c Config) GetURL() string {
	url := fmt.Sprintf("amqp://%s:%v@%s/%s", c.Username, c.Password, c.Host, c.Queue.Vhost)
	if c.Port != "" {
		url = fmt.Sprintf("amqp://%s:%v@%s/%s/:%s", c.Username, c.Password, c.Host, c.Queue.Vhost, c.Port)
	}
	return url
}
