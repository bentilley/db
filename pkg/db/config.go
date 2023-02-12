package db

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases map[string]Database
	Sessions  []Session
}

func (c *Config) SearchStrings() []string {
	searchDetails := [][]string{}
	for _, session := range c.Sessions {
		details := []string{session.String(), session.Description}
		searchDetails = append(searchDetails, details)
	}
	return ColumnFormat(searchDetails)
}

type Database struct {
	Config DatabaseConfig
}

func (d *Database) UnmarshalYAML(value *yaml.Node) error {
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
	case "postgres":
		d.Config = &Postgres{}
	}
	return value.Decode(d.Config)
}

func (d *Database) DatabaseSlug() string {
	return d.Config.DatabaseSlug()
}

type DatabaseConfig interface {
	DatabaseSlug() string
}

type Postgres struct {
	Description string
	Host        string
	Port        string
	Database    string
}

func (p *Postgres) DatabaseSlug() string {
	return fmt.Sprintf("%s:%s/%s", p.Host, p.Port, p.Database)
}

type Session struct {
	Description  string
	DatabaseName string    `yaml:"database"`
	Database     *Database `yaml:"-"`
	User         string
	Password     Password
}

func (s *Session) String() string {
	userSlug := s.PasswordlessUserSlug()
	databaseSlug := s.Database.DatabaseSlug()
	if userSlug == "" {
		return fmt.Sprintf("postgres://%s", databaseSlug)
	} else {
		return fmt.Sprintf("postgres://%s@%s", userSlug, databaseSlug)
	}
}

func (s *Session) URI() (string, error) {
	userSlug, err := s.UserSlug()
	if err != nil {
		return "", fmt.Errorf("get user slug: %w", err)
	}
	databaseSlug := s.Database.DatabaseSlug()
	if userSlug == "" {
		return fmt.Sprintf("postgres://%s", databaseSlug), nil
	} else {
		return fmt.Sprintf("postgres://%s@%s", userSlug, databaseSlug), nil
	}
}

func (s *Session) UserSlug() (string, error) {
	if s.User == "" {
		return "", nil
	}
	password, err := s.Password.GetPassword()
	if err != nil {
		return "", fmt.Errorf("get password: %w", err)
	}
	if password == "" {
		return s.User, nil
	} else {
		return fmt.Sprintf("%s:%s", s.User, password), nil
	}
}

func (s *Session) PasswordlessUserSlug() string {
	var password string
	if s.Password.HasPassword() {
		password = "***"
	}
	if password == "" {
		return s.User
	} else {
		return fmt.Sprintf("%s:%s", s.User, password)
	}
}

type Password struct {
	Config PasswordConfig
}

func (p *Password) HasPassword() bool {
	return p.Config != nil
}

func (p *Password) GetPassword() (string, error) {
	if p.Config == nil {
		return "", nil
	}
	return p.Config.GetPassword()
}

func (p *Password) UnmarshalYAML(value *yaml.Node) error {
	// handle the simple case where the password is in plaintext in the yaml.
	if value.Tag == "!!str" {
		var password PlainTextPassword
		p.Config = &password
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
		p.Config = &PassPassword{}
	case "env":
		p.Config = &EnvPassword{}
	}
	return value.Decode(p.Config)
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
	for i, session := range config.Sessions {
		database, ok := config.Databases[session.DatabaseName]
		if !ok {
			return nil, fmt.Errorf("could not find database %s", session.DatabaseName)
		}
		config.Sessions[i].Database = &database
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
