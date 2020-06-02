package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration definition
type Configuration struct {
	Server struct {
		Hostname string
		Port     string
	}
	Database struct {
		Hostname string
		Port     int
		Username string
		Password string
		Name     string
	}
}

// Config - global config variable
var Config *Configuration

// Load - Reads configuration from YAML file and ENV vars
func Load(filename string) {
	// Check config filename not empty
	if len(filename) == 0 {
		log.Println("Empty config filename.")
		return
	}

	// Open config file
	configFile, err := os.Open(filename)
	if err != nil {
		log.Println("Cannot open config file.")
		return
	}
	defer configFile.Close()

	// Read config file
	bConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Println(err)
		return
	}
	// Unmarshal YAML content
	err = yaml.Unmarshal(bConfig, &Config)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
