package config

const CODE_TTL = 300
const REPEAT_SENT_CODE_TTL = 60

type AppConfig struct {
	CodeTtl           int `yaml:"code_ttl"`
	RepeatSentCodeTtl int `yaml:"repeat_sent_code_ttl"`
}
