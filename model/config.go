package model

// Config represents the configuration for the Aerodatabox API
type Config struct {
	Host string `json:"host" mapstructure:"host"`
	Key  string `json:"key" mapstructure:"key"`
}
