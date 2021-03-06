package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
)

const (
	configPrefix = "seaside"
	Version      = "v1.1.1"
)

var (
	// will set by go tool link (-X option)
	clientID     string
	clientSecret string
)

type Credential struct {
	AccessToken string `json:"access_token"`
}

type Config struct {
	SeaOrigin         string     `split_words:"true" default:"https://c.linkage.community"`
	ClientID          string     `split_words:"true" required:"true"`
	ClientSecret      string     `split_words:"true" required:"true"`
	CredentialFile    string     `split_words:"true" default:"credential.json"`
	CurrentCredential Credential `ignored:"true"`
}

func Usage() {
	c := &Config{}
	envconfig.Usage(configPrefix, c)
}

func LoadConfig() (*Config, error) {
	if clientID != "" && clientSecret != "" {
		return &Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}, nil
	}

	// fallback
	c := &Config{}
	err := envconfig.Process(configPrefix, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) LoadCurrentCredential() error {
	bytes, err := ioutil.ReadFile(c.CredentialFile)
	if err != nil {
		return err
	}
	cr := &Credential{}
	if err := json.Unmarshal(bytes, cr); err != nil {
		return err
	}
	c.CurrentCredential = *cr
	return nil
}

func (c *Config) SaveCredential(crd *Credential) error {
	d, err := json.Marshal(crd)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.CredentialFile, d, 0644); err != nil {
		return err
	}
	return nil
}
