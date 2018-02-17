package gmig

import (
	"encoding/json"
	"io/ioutil"
)

// ConfigFilename is for reading bucket info
const ConfigFilename = "gmig.json"

// Config holds gmig program config
type Config struct {
	Bucket      string `json:"bucket"`
	StateObject string `json:"state"`
	Verbose     bool   `json:"verbose"`
}

// LoadConfig reads from gmig.json
func LoadConfig() (Config, error) {
	data, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}
