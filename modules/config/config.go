package config

import (
	"os"
	"sync"

	"github.com/lafriakh/kira/modules/log"
)

const delimiter = "="

// Config - This will store all the configurations we have.
type Config struct {
	lock   sync.RWMutex
	bucket map[string]any
}

// New - return kon instance
func New() *Config {
	return &Config{
		lock:   sync.RWMutex{},
		bucket: make(map[string]any),
	}
}

// Get - key value.
func (k *Config) Get(key string) any {
	k.lock.RLock()
	defer k.lock.RUnlock()

	if k.bucket[key] != nil {
		return k.bucket[key]
	}

	return nil
}

// Set ...
func (k *Config) Set(key string, value any) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.bucket[key] = value
}

// NewFromFile - read a file and return the configuration.
func NewFromFile(files ...string) *Config {
	kon := New()

	for _, file := range files {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			data, err := os.Open(file)
			if err != nil {
				log.Fatal(err)
			}
			defer data.Close()

			// Parse the file.
			conf, err := Parse(data)
			if err != nil {
				log.Panic(err)
			}

			// Set the configs.
			for key, value := range conf {
				kon.Set(key, value)
			}
		}
	}

	return kon
}

// NewFromMap return configuration from kira.Map type
func NewFromMap(m ...map[string]any) *Config {
	kon := New()
	for _, mapV := range m {
		for k, v := range mapV {
			kon.Set(k, v)
		}
	}
	return kon
}
