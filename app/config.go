package main

import (
	"encoding/json"
	"log"
	"os"
)

// SaConfig represents config object for SpendAgent
type SaConfig struct {
	Postgres       string `json:"postgres"`
	LogSQL         bool   `json:"log_sql"`
	AppUrl         string `json:"app_url"`
	ExtensionsHash string `json:"extensions_hash"`
}

// Load method loads config file in SaConfig object
func (sc *SaConfig) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[SaConfig.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[SaConfig.Load] Got error while decoding JSON: %v", err)
		return err
	}

	return nil
}
