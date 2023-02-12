package db_test

import (
	"testing"

	"github.com/bentilley/db/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testDB := db.Database{Config: &db.Postgres{
		Host:     "localhost",
		Port:     "5432",
		Database: "tester",
	}}
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
    type: postgres
    host: localhost
    port: 5432
    database: tester
sessions:
  - database: test-db
    user: testuser
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Sessions: []db.Session{
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
    type: postgres
    host: localhost
    port: 5432
    database: tester
sessions:
  - database: test-db
    user: testuser
    password: testpassword
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Sessions: []db.Session{
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
    type: postgres
    host: localhost
    port: 5432
    database: tester
sessions:
  - database: test-db
    user: testuser
    password:
        type: pass
        path: db/testuser/password
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Sessions: []db.Session{
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
    type: postgres
    host: localhost
    port: 5432
    database: tester
sessions:
  - database: test-db
    user: testuser
    password:
        type: env
        var: DB_PASSWORD
`,
			want: &db.Config{
				Databases: map[string]db.Database{"test-db": testDB},
				Sessions: []db.Session{
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

	localDev := db.Database{Config: &db.Postgres{
		Description: "Local Database",
		Host:        "localhost",
		Port:        "5432",
	}}
	want := &db.Config{
		Databases: map[string]db.Database{
			"local-dev": localDev,
		},
		Sessions: []db.Session{
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

func TestSessionString(t *testing.T) {
	tests := []struct {
		name    string
		session db.Session
		want    string
	}{
		{
			name: "no user",
			session: db.Session{
				DatabaseName: "test-db",
				Database: &db.Database{Config: &db.Postgres{
					Host:     "localhost",
					Port:     "5432",
					Database: "tester",
				}},
			},
			want: "postgres://localhost:5432/tester",
		},
		{
			name: "just user",
			session: db.Session{
				DatabaseName: "test-db",
				Database: &db.Database{Config: &db.Postgres{
					Host:     "localhost",
					Port:     "5432",
					Database: "tester",
				}},
				User: "someuser",
			},
			want: "postgres://someuser@localhost:5432/tester",
		},
		{
			name: "user and password",
			session: db.Session{
				DatabaseName: "test-db",
				Database: &db.Database{Config: &db.Postgres{
					Host:     "some.host",
					Port:     "5432",
					Database: "tester",
				}},
				User: "someuser",
				Password: db.Password{Config: &db.PlainTextPassword{
					Value: "somepassword",
				}},
			},
			want: "postgres://someuser:***@some.host:5432/tester",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.session.String())
		})
	}
}

func TestURI(t *testing.T) {
	tests := []struct {
		name    string
		session db.Session
		want    string
	}{
		{
			name: "no user",
			session: db.Session{
				DatabaseName: "test-db",
				Database: &db.Database{Config: &db.Postgres{
					Host:     "localhost",
					Port:     "5432",
					Database: "tester",
				}},
			},
			want: "postgres://localhost:5432/tester",
		},
		{
			name: "just user",
			session: db.Session{
				DatabaseName: "test-db",
				Database: &db.Database{Config: &db.Postgres{
					Host:     "localhost",
					Port:     "5432",
					Database: "tester",
				}},
				User: "someuser",
			},
			want: "postgres://someuser@localhost:5432/tester",
		},
		{
			name: "user and password",
			session: db.Session{
				DatabaseName: "test-db",
				Database: &db.Database{Config: &db.Postgres{
					Host:     "some.host",
					Port:     "5432",
					Database: "tester",
				}},
				User: "someuser",
				Password: db.Password{Config: &db.PlainTextPassword{
					Value: "somepassword",
				}},
			},
			want: "postgres://someuser:somepassword@some.host:5432/tester",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := tt.session.URI()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.want, uri)
		})
	}
}
