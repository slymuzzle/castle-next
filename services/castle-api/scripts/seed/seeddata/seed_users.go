package seeddata

import (
	"context"
	"fmt"

	"journeyhub/ent/room"
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

	ctx := context.Background()

	adminUser, err := entClient.User.
		Create().
		SetFirstName("Admin").
		SetLastName("Admin").
		SetEmail("admin@admin.com").
		SetNickname("admin").
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return err
	}

	testUser1, err := entClient.User.
		Create().
		SetFirstName("Эпона").
		SetLastName("Ортен").
		SetEmail("test@test.com").
		SetNickname("test").
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return err
	}

	testUser2, err := entClient.User.
		Create().
		SetFirstName("Боеслав").
		SetLastName("Рей").
		SetEmail("test1@test.com").
		SetNickname("test1").
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return err
	}

	testUser3, err := entClient.User.
		Create().
		SetFirstName("Реджинолд").
		SetLastName("Вуд").
		SetEmail("test3@test.com").
		SetNickname("test3").
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return err
	}

	testUser4, err := entClient.User.
		Create().
		SetFirstName("Готчолк").
		SetLastName("Галилей").
		SetEmail("test4@test.com").
		SetNickname("test4").
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return err
	}

	testRoom1, err := entClient.Room.
		Create().
		SetType(room.TypePersonal).
		Save(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", testUser1.FirstName, testUser1.LastName)).
		SetUser(adminUser).
		SetRoom(testRoom1).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", adminUser.FirstName, adminUser.LastName)).
		SetUser(testUser1).
		SetRoom(testRoom1).
		Exec(ctx)
	if err != nil {
		return err
	}

	testRoom2, err := entClient.Room.
		Create().
		SetType(room.TypePersonal).
		Save(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", testUser2.FirstName, testUser2.LastName)).
		SetUser(adminUser).
		SetRoom(testRoom2).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", adminUser.FirstName, adminUser.LastName)).
		SetUser(testUser2).
		SetRoom(testRoom2).
		Exec(ctx)
	if err != nil {
		return err
	}

	testRoom3, err := entClient.Room.
		Create().
		SetType(room.TypePersonal).
		Save(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", testUser3.FirstName, testUser3.LastName)).
		SetUser(adminUser).
		SetRoom(testRoom3).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", adminUser.FirstName, adminUser.LastName)).
		SetUser(testUser3).
		SetRoom(testRoom3).
		Exec(ctx)
	if err != nil {
		return err
	}

	testRoom4, err := entClient.Room.
		Create().
		SetType(room.TypePersonal).
		Save(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", testUser4.FirstName, testUser4.LastName)).
		SetUser(adminUser).
		SetRoom(testRoom4).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = entClient.RoomMember.
		Create().
		SetName(fmt.Sprintf("%s %s", adminUser.FirstName, adminUser.LastName)).
		SetUser(testUser4).
		SetRoom(testRoom4).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = entClient.UserContact.CreateBulk(
		entClient.UserContact.
			Create().
			SetRoom(testRoom1).
			SetUser(adminUser).
			SetContact(testUser1),
		entClient.UserContact.
			Create().
			SetRoom(testRoom1).
			SetUser(testUser1).
			SetContact(adminUser),
		entClient.UserContact.
			Create().
			SetRoom(testRoom2).
			SetUser(adminUser).
			SetContact(testUser2),
		entClient.UserContact.
			Create().
			SetRoom(testRoom2).
			SetUser(testUser2).
			SetContact(adminUser),
		entClient.UserContact.
			Create().
			SetRoom(testRoom3).
			SetUser(adminUser).
			SetContact(testUser3),
		entClient.UserContact.
			Create().
			SetRoom(testRoom3).
			SetUser(testUser3).
			SetContact(adminUser),
		entClient.UserContact.
			Create().
			SetRoom(testRoom4).
			SetUser(adminUser).
			SetContact(testUser4),
		entClient.UserContact.
			Create().
			SetRoom(testRoom4).
			SetUser(testUser4).
			SetContact(adminUser),
	).OnConflict().
		DoNothing().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
