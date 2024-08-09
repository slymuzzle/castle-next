env "local" {
  // Declare where the schema definition resides.
  src = "ent://ent/schema"

  // Define the URL of the database which is managed
  // in this environment.
  url = "postgres://admin:123456@localhost:20082/castle?search_path=public&sslmode=disable"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "postgres://admin:123456@localhost:20082/atlas?search_path=public&sslmode=disable"

  migration {
    // URL where the migration directory resides.
    dir = "file://migrations"
  }
}

env "docker" {
  // Declare where the schema definition resides.
  src = "ent://ent/schema"

  // Define the URL of the database which is managed
  // in this environment.
  url = "postgres://admin:123456@postgres:5432/castle?search_path=public&sslmode=disable"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "postgres://admin:123456@postgres:5432/atlas?search_path=public&sslmode=disable"

  migration {
    // URL where the migration directory resides.
    dir = "file://migrations"
  }
}
