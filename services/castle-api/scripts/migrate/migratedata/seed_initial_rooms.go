package migratedata

import (
	"context"
	"fmt"
	"journeyhub/ent"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect/sql/schema"
)

// SeedUsers add the initial users to the database.
func SeedInitialRooms(dialect string, dir *migrate.LocalDir) error {
	w := &schema.DirWriter{Dir: dir}
	client := ent.NewClient(ent.Driver(schema.NewWriteDriver(dialect, w)))

	// The statement that generates the INSERT statement.
	err := client.Room.CreateBulk(
		client.Room.
			Create().
			SetName("Support"),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed generating statement: %w", err)
	}

	// Write the content to the migration directory.
	return w.FlushChange(
		"seed_initial_rooms",
		"Add the initial rooms to the database.",
	)
}
