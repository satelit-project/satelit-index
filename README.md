# satelit-index

Prepares index files for data import from external sources.

## Dependencies

- Go 1.13
- Docker
- Fish shell

All additional project dependencies can be installed via [`tools/db.fish`](tools/db.fish). See `--help` to find out how.

## Tools

[`tools`](tools) directory contains two additional scripts:

- [`tools/db.fish`](tools/db.fish) - can start/stop Postgres docker container. It also allows you to attach to container and interact with the DB via `psql`. The database data will be stored in `.mount` directory.

- [`tools/ngx.fish`](tools/db.fish) - can start/stop Nginx docker container. It will serve files from `static` directory in
the root project directory.

### DB

After first Postgres container run you need to migrate it's schema. This can be done by using `goose`:
```bash
goose postgres "<url>" up
```

To generate SQL queries and native types you can use `sqlc`:
```bash
sqlc generate
```
