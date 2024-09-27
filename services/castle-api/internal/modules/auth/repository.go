package auth

import (
	"context"

	"journeyhub/ent"
	"journeyhub/ent/device"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
)

type Repository interface {
	FindUserByID(
		ctx context.Context,
		userID pulid.ID,
	) (*ent.User, error)
	FindUserByNickname(
		ctx context.Context,
		nickname string,
	) (*ent.User, error)
	CreateUser(
		ctx context.Context,
		firstName string,
		lastName string,
		nickname string,
		password string,
	) (*ent.User, error)
	CreateOrUpdateUserDevice(
		ctx context.Context,
		userID pulid.ID,
		deviceID string,
		fcmToken string,
	) (pulid.ID, error)
}

type repository struct {
	entClient *ent.Client
}

func NewRepository(entClient *ent.Client) Repository {
	return &repository{
		entClient: entClient,
	}
}

func (r *repository) FindUserByID(ctx context.Context, userID pulid.ID) (*ent.User, error) {
	return r.getClient(ctx).User.
		Get(ctx, userID)
}

func (r *repository) FindUserByNickname(
	ctx context.Context,
	nickname string,
) (*ent.User, error) {
	return r.getClient(ctx).User.
		Query().
		Where(
			user.Nickname(nickname),
		).
		Only(ctx)
}

func (r *repository) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
	nickname string,
	password string,
) (*ent.User, error) {
	return r.getClient(ctx).User.
		Create().
		SetFirstName(firstName).
		SetLastName(lastName).
		SetNickname(nickname).
		SetPassword(password).
		Save(ctx)
}

func (r *repository) CreateOrUpdateUserDevice(
	ctx context.Context,
	userID pulid.ID,
	deviceID string,
	fcmToken string,
) (pulid.ID, error) {
	return r.getClient(ctx).Device.
		Create().
		SetUserID(userID).
		SetDeviceID(deviceID).
		SetFcmToken(fcmToken).
		OnConflictColumns(device.FieldDeviceID).
		UpdateNewValues().
		ID(ctx)
}

func (r *repository) getClient(ctx context.Context) *ent.Client {
	var client *ent.Client
	if clientFromCtx := ent.FromContext(ctx); clientFromCtx != nil {
		client = clientFromCtx
	} else {
		client = r.entClient
	}
	return client
}
