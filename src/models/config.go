package models

import (
	"encoding/json"
	"io/ioutil"
)

// Config type for config file.
type Config User

// Init initializes config.
func (c *Config) Init() Config {
	return Config{Username: "", Password: ""}
}

// Update updates config file with given parameters.
func (c *Config) Update(u, p string) Config {

	c.Username, c.Password = u, p

	// write config file.
	WriteConfig(c, u)

	return Config{Username: u, Password: p}
}

// WriteConfig writes config to file.
func WriteConfig(data interface{}, u string) {
	file, _ := json.MarshalIndent(data, "", " ")

	name := u + ".json"

	_ = ioutil.WriteFile(name, file, 0644)
}
