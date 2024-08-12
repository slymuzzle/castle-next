package migratedata

import (
	"context"
	"fmt"
	"journeyhub/ent"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect/sql/schema"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers add the initial users to the database.
func SeedInitialUsers(dialect string, dir *migrate.LocalDir) error {
	w := &schema.DirWriter{Dir: dir}
	client := ent.NewClient(ent.Driver(schema.NewWriteDriver(dialect, w)))

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("12345678"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// The statement that generates the INSERT statement.
	err = client.User.CreateBulk(
		client.User.
			Create().
			SetFirstName("Admin").
			SetLastName("Admin").
			SetEmail("admin@admin.com").
			SetNickname("admin").
			SetPassword(string(hashedPassword)),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed generating statement: %w", err)
	}

	// Write the content to the migration directory.
	return w.FlushChange(
		"seed_initial_users",
		"Add the initial users to the database.",
	)
}
