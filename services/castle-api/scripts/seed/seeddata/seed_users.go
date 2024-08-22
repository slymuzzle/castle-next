package seeddata

import (
	"context"

	"journeyhub/internal/platform/db"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(dbService db.Service) error {
	entClient := dbService.Client()

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("12345678"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return entClient.User.CreateBulk(
		entClient.User.
			Create().
			SetFirstName("Admin").
			SetLastName("Admin").
			SetEmail("admin@admin.com").
			SetNickname("admin").
			SetPassword(string(hashedPassword)),
		entClient.User.
			Create().
			SetFirstName("Test").
			SetLastName("Test").
			SetEmail("test@test.com").
			SetNickname("test").
			SetPassword(string(hashedPassword)),
	).OnConflict().
		DoNothing().
		Exec(context.Background())
}
