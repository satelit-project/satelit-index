# satelit-index

Searches for data to be scraped from external sources. 

## Dependencies

- Go 1.14

All additional project dependencies can be installed via [`tools/deps.fish`](tools/deps.fish). See `--help` to find out how.

### DB

After first Postgres container run you need to migrate it's schema. This can be done by using `goose`:
```bash
goose postgres "<url>" up
```

To generate SQL queries and native types you can use `sqlc`:
```bash
sqlc generate
```
