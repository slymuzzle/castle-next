package call

import (
	"context"
	"time"

	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/modules/auth"
	"journeyhub/internal/platform/config"

	"github.com/k0kubun/pp/v3"
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
	userID, err := s.authService.Auth(ctx)
	if err != nil {
		return "", err
	}

	pp.Print(s.config)
	at := livekitauth.NewAccessToken(s.config.AccessKey, s.config.SecretKey)
	grant := &livekitauth.VideoGrant{
		RoomJoin: true,
		Room:     string(roomID),
	}
	at.AddGrant(grant).
		SetIdentity(string(userID)).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
