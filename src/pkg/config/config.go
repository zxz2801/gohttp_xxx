package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	mutex = sync.RWMutex{}
)

// HTTPServer :
type HTTPServer struct {
	Port           int `json:"port" yaml:"port"`
	ReadTimeout    int `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   int `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout    int `json:"idle_timeout" yaml:"idle_timeout"`
	MaxHeaderBytes int `json:"max_header_bytes" yaml:"max_header_bytes"`
}

// Log :
type Log struct {
	Base        string `json:"base" yaml:"base"`
	File        string `json:"file" yaml:"file"`
	RotateSize  int64  `json:"rotate_size" yaml:"rotate_size"`
	RotateCount int    `json:"rotate_count" yaml:"rotate_count"`
	Level       int64  `json:"level" yaml:"level"`
}

// Config :
type Config struct {
	HTTPServer `yaml:"http_server"`
	Log        `yaml:"log"`
}

var defaultConfig = Config{
	HTTPServer{
		Port:           11111,
		ReadTimeout:    60,
		WriteTimeout:   60,
		IdleTimeout:    60,
		MaxHeaderBytes: 20 * 1024 * 1024,
	},
	Log{
		Base:        "../log/",
		File:        "gohttp_xxx.log",
		RotateSize:  104857600,
		RotateCount: 5,
		Level:       3,
	},
}

func initFromYAML(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, &defaultConfig)
	if err != nil {
		return err
	}
	return nil
}

// InitGlobalConfig ...
func InitGlobalConfig(filename string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if strings.HasSuffix(filename, ".yml") {
		return initFromYAML(filename)
	}
	return errors.New("unknown file format")

}

// Global ...
func Global() *Config {
	mutex.RLock()
	defer mutex.RUnlock()

	return &defaultConfig
}
