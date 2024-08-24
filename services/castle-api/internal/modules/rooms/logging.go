package rooms

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"

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

func (s *loggingService) IncrementRoomVersion(
	ctx context.Context,
	ID pulid.ID,
	lastMessage *ent.Message,
) (room *ent.Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "FindRoomByMessage",
			"ID", ID,
			"lastMessage", lastMessage.ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.IncrementRoomVersion(ctx, ID, lastMessage)
}

func (s *loggingService) SubscribeToRoomsUpdatedEvent(
	ctx context.Context,
) (ch <-chan *model.RoomUpdatedEvent, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomsUpdatedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToRoomsUpdatedEvent(ctx)
}

func (s *loggingService) DeleteRoomMember(
	ctx context.Context,
	roomMemberID pulid.ID,
) (roomMember *ent.RoomMember, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteRoomMember",
			"roomMemberID", roomMemberID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteRoomMember(ctx, roomMemberID)
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
