package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "../backend/config.yaml"

type ConfigYaml struct {
	LogLevel string `yaml:"log_level"`
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	Cache struct {
		Expiration string `yaml:"expiration"`
		Purge      string `yaml:"purge"`
	} `yaml:"cache"`
	ClientSession struct {
		JwtSignature string `yaml:"jwt_signature"`
		Expire       string `yaml:"expire"`
	} `yaml:"client_session"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Reservation struct {
		Expiration string `yaml:"expiration"`
	} `yaml:"reservation"`
	Mail struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Template struct {
			Subject   string `json:"body"`
			Body      string `yaml:"body"`
			BodyPlain string `yaml:"body_plain"`
		} `json:"template"`
	} `yaml:"mail"`
}

type CacheConfig struct {
	Expiration time.Duration
	Purge      time.Duration
}

var config ConfigYaml

func loadConfig() ConfigYaml {
	config := ConfigYaml{}

	yamlFile, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		panic(fmt.Sprintf("Error opening config-file: %v", err))
	}

	reader := bytes.NewReader(yamlFile)

	dec := yaml.NewDecoder(reader)
	dec.KnownFields(true)
	err = dec.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config-file: %v", err)
		os.Exit(1)
	}

	return config
}

func writeConfig() {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	// Can set default indent here on the encoder
	if err := enc.Encode(&config); err != nil {
		panic(err)
	} else {
		if err := os.WriteFile(CONFIG_PATH, buf.Bytes(), 0644); err != nil {
			panic(err)
		}
	}
}

func init() {
	config = loadConfig()
}
