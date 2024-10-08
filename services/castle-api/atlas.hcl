env "local" {
  // Declare where the schema definition resides.
  src = "ent://ent/schema"

  // Define the URL of the database which is managed
  // in this environment.
  url = "postgres://postgres:pass@localhost:5432/database?search_path=public&sslmode=disable"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    // URL where the migration directory resides.
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
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

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
