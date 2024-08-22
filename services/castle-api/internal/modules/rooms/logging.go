package rooms

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

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

func (s *loggingService) FindOrCreatePersonalRoom(
	ctx context.Context,
	targetUserID pulid.ID,
) (room *ent.Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "FindOrCreatePersonalRoom",
			"targetUserID", targetUserID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.FindOrCreatePersonalRoom(ctx, targetUserID)
}

func (s *loggingService) FindRoomByMessage(
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

func (s *loggingService) CreateRoom(
	ctx context.Context,
	input CreateRoomInput,
) (room *ent.Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "CreateRoom",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateRoom(ctx, input)
}

func (s *loggingService) UpdateRoom(
	ctx context.Context,
	ID pulid.ID,
	input UpdateRoomInput,
) (room *ent.RoomEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "UpdateRoom",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateRoom(ctx, ID, input)
}

func (s *loggingService) SubscribeToRoomsUpdatedEvent(
	ctx context.Context,
) (ch <-chan *ent.RoomEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomsUpdatedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToRoomsUpdatedEvent(ctx)
}

func (s *loggingService) DeleteRoom(
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
