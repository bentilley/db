package db_test

import (
	"testing"

	"github.com/bentilley/db/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testDB := db.Database{
		Host:     "localhost",
		Port:     "5432",
		Database: "tester",
	}
	tests := []struct {
		name       string
		yamlConfig string
		want       *db.Config
	}{
		{
			name: "simple",
			yamlConfig: `
databases:
  test-db:
    host: localhost
    port: 5432
    database: tester
connections:
  - database: test-db
    user: testuser
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Connections: []db.Connection{
					{
						DatabaseName: "test-db",
						Database:     &testDB,
						User:         "testuser",
					},
				},
			},
		},
		{
			name: "text password",
			yamlConfig: `
databases:
  test-db:
    host: localhost
    port: 5432
    database: tester
connections:
  - database: test-db
    user: testuser
    password: testpassword
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Connections: []db.Connection{
					{
						DatabaseName: "test-db",
						Database:     &testDB,
						User:         "testuser",
						Password: db.Password{
							Config: &db.PlainTextPassword{Value: "testpassword"},
						},
					},
				},
			},
		},
		{
			name: "pass password",
			yamlConfig: `
databases:
  test-db:
    host: localhost
    port: 5432
    database: tester
connections:
  - database: test-db
    user: testuser
    password:
        type: pass
        path: db/testuser/password
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Connections: []db.Connection{
					{
						DatabaseName: "test-db",
						Database:     &testDB,
						User:         "testuser",
						Password: db.Password{
							Config: &db.PassPassword{Path: "db/testuser/password"},
						},
					},
				},
			},
		},
		{
			name: "env password",
			yamlConfig: `
databases:
  test-db:
    host: localhost
    port: 5432
    database: tester
connections:
  - database: test-db
    user: testuser
    password:
        type: env
        var: DB_PASSWORD
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Connections: []db.Connection{
					{
						DatabaseName: "test-db",
						Database:     &testDB,
						User:         "testuser",
						Password: db.Password{
							Config: &db.EnvPassword{Var: "DB_PASSWORD"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.ParseConfig([]byte(tt.yamlConfig))
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	got, err := db.LoadConfig("testdata/example-config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	localDev := db.Database{
		Host: "localhost",
		Port: "5432",
	}
	want := &db.Config{
		Databases: map[string]db.Database{
			"local-dev": localDev,
		},
		Connections: []db.Connection{
			{
				DatabaseName: "local-dev",
				Database:     &localDev,
				User:         "testuser",
				Password: db.Password{Config: &db.PassPassword{
					Path: "db/testuser/password",
				}},
			},
			{
				DatabaseName: "local-dev",
				Database:     &localDev,
				User:         "anotheruser",
				Password: db.Password{Config: &db.EnvPassword{
					Var: "DB_PASSWORD",
				}},
			},
		},
	}

	assert.Equal(t, want, got)
}
