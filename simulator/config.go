package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ConfigEntry struct {
	LastMeasurement    time.Time
	DowntimeSimulation time.Duration
}

type Config struct {
	data     *ConfigEntry
	filename string
}

func (c *Config) LoadConfig() error {
	if _, err := os.Stat(c.filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %w", err)
	}

	fileData, err := os.ReadFile(c.filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}
	var data ConfigEntry
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}
	c.data = &data
	return nil
}

func (c *Config) SaveConfig() error {
	file, err := os.Create(c.filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	fileData, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("could not encode data: %w", err)
	}

	_, err = file.Write(fileData)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}
	log.Printf("wrote config file to %s", c.filename)

	return nil
}
