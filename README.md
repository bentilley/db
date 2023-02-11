# `db`

A simple cli for defining database URIs with support for multiple password
backends.

I created this tool to solve the problem of having to keep database URIs
containing passwords in plaintext. With `db` you define all your database
connections in a yaml file, including where to find the relevant password
information, and then `db` compiles the URIs on the fly and outputs the results
through a fuzzy finder of your choosing (default `fzf`).

## Sample Config

The default location for the config file is `~/.config/db/config.yaml`.
A sample config looks like this:

```yaml
databases:
  local-dev:
    host: localhost
    port: 5432

  another-db:
    host: some.host
    port: 15500

connections:
  - database: local-dev
    user: testuser
    password: plaintext

  - database: local-dev
    user: adminuser
    password:
      type: pass
      path: db/testuser/password

  - database: another-db
    user: anotheruser
    password:
      type: env
      var: DB_PASSWORD
```
