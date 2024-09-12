package rooms

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type serviceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &serviceLogging{logger, s}
}

// func (s *loggingService) FindOrCreatePersonalRoom(
// 	ctx context.Context,
// 	targetUserID pulid.ID,
// ) (room *ent.Room, err error) {
// 	defer func(begin time.Time) {
// 		level.Debug(s.logger).Log(
// 			"method", "FindOrCreatePersonalRoom",
// 			"targetUserID", targetUserID,
// 			"took", time.Since(begin),
// 			"err", err,
// 		)
// 	}(time.Now())
// 	return s.Service.FindOrCreatePersonalRoom(ctx, targetUserID)
// }

func (s *serviceLogging) FindRoomByMessage(
	ctx context.Context,
	messageID pulid.ID,
) (room *ent.Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "FindRoomByMessage",
			"messageID", messageID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.FindRoomByMessage(ctx, messageID)
}

func (s *serviceLogging) IncrementRoomVersion(
	ctx context.Context,
	ID pulid.ID,
	lastMessage *ent.Message,
) (room *ent.Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "FindRoomByMessage",
			"ID", ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.IncrementRoomVersion(ctx, ID, lastMessage)
}

func (s *serviceLogging) DeleteRoom(
	ctx context.Context,
	ID pulid.ID,
) (room *ent.Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteRoom",
			"ID", ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteRoom(ctx, ID)
}
