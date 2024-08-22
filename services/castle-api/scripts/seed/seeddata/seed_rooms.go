package seeddata

import (
	"context"

	"journeyhub/ent/room"
	"journeyhub/internal/platform/db"
)

func SeedRooms(dbService db.Service) error {
	entClient := dbService.Client()

	return entClient.Room.CreateBulk(
		entClient.Room.
			Create().
			SetName("Support").
			SetType(room.TypePersonal),
	).OnConflict().
		DoNothing().
		Exec(context.Background())
}
