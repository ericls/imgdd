# Migrations

This folder contains migrations written in .sql files. The migration is ultimately managed by the [migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4@v4.16.2) package.


Migrations are the source of truth of database structures. "Models" will be generated based on the tables.

### Adding new migrations
1. Try to provide both an up and a down migration
1. Table names should be `<singular_name>_table`. This does not mean that I think it's a good naming convention for tables, this is just to stop cointributors from arguing or thinking about this.
1. The primary key of an column should be called `id`, it should not be prepended with table/entity names. `<entity>_id` is reserved for references.