package db

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases   map[string]Database
	Connections []Connection
}

func (c *Config) URIs() ([]string, error) {
	uris := []string{}
	for _, connection := range c.Connections {
		uri, err := connection.URI()
		if err != nil {
			return nil, fmt.Errorf("get uri: %w", err)
		}
		uris = append(uris, uri)
	}
	return uris, nil
}

type Database struct {
	Description string
	Host        string
	Port        string
	Database    string
}

type Connection struct {
	Description  string
	DatabaseName string    `yaml:"database"`
	Database     *Database `yaml:"-"`
	User         string
	Password     Password
}

func (c *Connection) URI() (string, error) {
	userSlug, err := c.UserSlug()
	if err != nil {
		return "", fmt.Errorf("get user slug: %w", err)
	}
	return fmt.Sprintf(
		"postgres://%s@%s:%s/%s",
		userSlug,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	), nil
}

func (c *Connection) UserSlug() (string, error) {
	if c.User == "" {
		return "", nil
	}
	password, err := c.Password.GetPassword()
	if err != nil {
		return "", fmt.Errorf("get password: %w", err)
	}
	if password == "" {
		return c.User, nil
	} else {
		return fmt.Sprintf("%s:%s", c.User, password), nil
	}
}

type Password struct {
	Config PasswordConfig
}

func (a *Password) GetPassword() (string, error) {
	if a.Config == nil {
		return "", nil
	}
	return a.Config.GetPassword()
}

func (a *Password) UnmarshalYAML(value *yaml.Node) error {
	// handle the simple case where the password is in plaintext in the yaml.
	if value.Tag == "!!str" {
		var password PlainTextPassword
		a.Config = &password
		return value.Decode(&password.Value)
	}

	typeIndex := -1
	for i, node := range value.Content {
		if node.Tag == "!!str" && node.Value == "type" {
			typeIndex = i
			break
		}
	}
	if typeIndex == -1 || typeIndex+1 > len(value.Content) {
		return fmt.Errorf("could not find type node")
	}
	passwordType := value.Content[typeIndex+1].Value
	switch passwordType {
	case "pass":
		a.Config = &PassPassword{}
	case "env":
		a.Config = &EnvPassword{}
	}
	return value.Decode(a.Config)
}

type PasswordConfig interface {
	GetPassword() (string, error)
}

func ParseConfig(configYaml []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(configYaml, config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}
	for i, connection := range config.Connections {
		database, ok := config.Databases[connection.DatabaseName]
		if !ok {
			return nil, fmt.Errorf("could not find database %s", connection.DatabaseName)
		}
		config.Connections[i].Database = &database
	}
	return config, nil
}

func LoadConfig(file string) (*Config, error) {
	configYaml, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	config, err := ParseConfig(configYaml)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return config, nil
}
