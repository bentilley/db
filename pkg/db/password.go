package db

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type PlainTextPassword struct {
	Value string
}

func (p *PlainTextPassword) GetPassword() (string, error) {
	return p.Value, nil
}

type PassPassword struct {
	Path string
}

func (p *PassPassword) GetPassword() (string, error) {
	cmd := exec.Command("pass", "show", p.Path)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("run command: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

type EnvPassword struct {
	Var string
}

func (p *EnvPassword) GetPassword() (string, error) {
	if v, ok := os.LookupEnv(p.Var); ok {
		return v, nil
	} else {
		return "", fmt.Errorf("environment variable %q not set", p.Var)
	}
}
