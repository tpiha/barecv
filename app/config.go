package main

import (
	"encoding/json"
	"log"
	"os"
)

// BareCVConfig represents config object for SpendAgent
type BareCVConfig struct {
	Postgres  string `json:"postgres"`
	LogSQL    bool   `json:"log_sql"`
	AppUrl    string `json:"app_url"`
	CVBase    string `json:"cv_base"`
	RubberBin string `json:"rubber_bin"`
}

// Load method loads config file in BareCVConfig object
func (sc *BareCVConfig) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[BareCVConfig.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[BareCVConfig.Load] Got error while decoding JSON: %v", err)
		return err
	}

	return nil
}
