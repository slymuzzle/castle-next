package seeddata

import (
	"context"
	"journeyhub/internal/db"
)

func SeedRooms(dbService db.Service) error {
	entClient := dbService.Client()

	return entClient.Room.CreateBulk(
		entClient.Room.
			Create().
			SetName("Support"),
	).OnConflict().
		DoNothing().
		Exec(context.Background())
}
