package conf_reader

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	CacheCapacity int64 `yaml:"cache_capacity"`
	CacheDuration int64 `yaml:"cache_duration_minutes"`
	MaxRetries    int64 `yaml:"max_retries"`
	DialTimeout   int64 `yaml:"dial_timeout"`
	Timeout       int64 `yaml:"timeout"`
}

func (c *Conf) GetConf(config_path string) *Conf {

	yamlFile, err := os.ReadFile(config_path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
