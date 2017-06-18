package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Port      string `yaml:"port"`
	Host      string `yaml:"host"`
	QueueName string `yaml:"queue_name"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
	return fmt.Sprintf("amqp://%v:%v@%v:%v/", c.Username, c.Password, c.Host, c.Port)
}
