package config

import (
	"os"
	"strconv"
	"time"
)

const defaultDb = 0
const defaultMaxRetries = 5
const defaultDialTimeout = 10
const defaultTimeout = 5
const defaultCodeTtl = 300
const defaultRepeatSentCodeTtl = 60

type StorageConfig struct {
	Addr        string        `yaml:"addr"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

type AppConfig struct {
	CodeTtl           int `yaml:"code_ttl"`
	RepeatSentCodeTtl int `yaml:"repeat_sent_code_ttl"`
	StorageConfig
}

func (cnf *AppConfig) LoadConfig() (*AppConfig, error) {
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbNumber = defaultDb
	}
	maxRetries, err := strconv.Atoi(os.Getenv("REDIS_MAX_RETRIES"))
	if err != nil {
		maxRetries = defaultMaxRetries
	}
	dialTimeout, err := strconv.Atoi(os.Getenv("REDIS_DIAL_TIMEOUT"))
	if err != nil {
		dialTimeout = defaultDialTimeout
	}
	timeout, err := strconv.Atoi(os.Getenv("REDIS_TIMEOUT"))
	if err != nil {
		timeout = defaultTimeout
	}

	codeTtl, err := strconv.Atoi(os.Getenv("CODE_TTL"))
	if err != nil {
		codeTtl = defaultCodeTtl
	}
	repeatSentCodeTtl, err := strconv.Atoi(os.Getenv("REPEAT_SENT_CODE_TTL"))
	if err != nil {
		repeatSentCodeTtl = defaultRepeatSentCodeTtl
	}

	appConfig := AppConfig{
		CodeTtl:           codeTtl,
		RepeatSentCodeTtl: repeatSentCodeTtl,
		StorageConfig: StorageConfig{
			Addr:        os.Getenv("REDIS_HOST"),
			DB:          dbNumber,
			MaxRetries:  maxRetries,
			DialTimeout: time.Duration(dialTimeout) * time.Second,
			Timeout:     time.Duration(timeout) * time.Second,
		},
	}

	return &appConfig, nil
}
