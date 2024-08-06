package chat

import (
	"context"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) SendMessage(
	ctx context.Context,
	roomID pulid.ID,
	userID pulid.ID,
	content string,
) (msg *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SendMessage",
			"roomID", roomID,
			"userID", userID,
			"content", content,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendMessage(ctx, roomID, userID, content)
}

func (s *loggingService) Subscribe(roomID pulid.ID) (ch <-chan *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Subscribe",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Subscribe(roomID)
}
