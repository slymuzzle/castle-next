package call

import (
	"context"
	"fmt"
	"time"

	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/platform/config"

	livekitauth "github.com/livekit/protocol/auth"
)

type Service interface {
	GetCallJoinToken(
		ctx context.Context,
		roomID pulid.ID,
	) (string, error)
}

type service struct {
	config      config.LivekitConfig
	authService auth.Service
}

func NewService(
	config config.LivekitConfig,
	authService auth.Service,
) Service {
	return &service{
		config:      config,
		authService: authService,
	}
}

func (s *service) GetCallJoinToken(
	ctx context.Context,
	roomID pulid.ID,
) (string, error) {
	user, err := s.authService.AuthUser(ctx)
	if err != nil {
		return "", err
	}

	at := livekitauth.NewAccessToken(s.config.AccessKey, s.config.SecretKey)
	grant := &livekitauth.VideoGrant{
		RoomJoin: true,
		Room:     string(roomID),
	}
	at.AddGrant(grant).
		SetIdentity(string(user.ID)).
		SetName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
