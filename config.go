package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server   Server `yaml:"server"`
}

type Server struct {
	Token      string     `yaml:"token"`
	Banterms   []string   `yaml:"banterms"`
}


/**
Get config from yaml file
*/
func ConfigFromFile(filePath string) (*Config, error) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	conf := Config{}
	err = yaml.Unmarshal([]byte(dat), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
