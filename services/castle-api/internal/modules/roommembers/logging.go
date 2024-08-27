package roommembers

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

func (s *loggingService) CreateRoomMembers(
	ctx context.Context,
	inputs []CreateRoomMemberInput,
) (roomMembers []*ent.RoomMember, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "CreateRoomMembers",
			"ID", len(roomMembers),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateRoomMembers(ctx, inputs)
}

func (s *loggingService) SubscribeToRoomMemberCreatedEvent(
	ctx context.Context,
) (ch <-chan *ent.RoomMemberEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomMemberCreatedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToRoomMemberCreatedEvent(ctx)
}

func (s *loggingService) IncrementUnreadMessagesCount(
	ctx context.Context,
	ID pulid.ID,
) (roomMembersToNotify []*ent.RoomMember, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "IncrementUnreadMessagesCount",
			"ID", ID,
			"roomMembersToNotifyCount", len(roomMembersToNotify),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.IncrementUnreadMessagesCount(ctx, ID)
}

func (s *loggingService) SubscribeToRoomMemberUpdatedEvent(
	ctx context.Context,
) (ch <-chan *ent.RoomMemberEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomMemberUpdatedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToRoomMemberUpdatedEvent(ctx)
}

func (s *loggingService) DeleteRoomMember(
	ctx context.Context,
	ID pulid.ID,
) (roomMember *ent.RoomMember, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteRoomMember",
			"ID", ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteRoomMember(ctx, ID)
}

func (s *loggingService) SubscribeToRoomMemberDeletedEvent(
	ctx context.Context,
) (ch <-chan pulid.ID, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomMemberDeletedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToRoomMemberDeletedEvent(ctx)
}

func (s *loggingService) MarkRoomMemberAsSeen(
	ctx context.Context,
	ID pulid.ID,
) (roomMember *ent.RoomMember, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "MarkRoomMemberAsSeen",
			"ID", ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.MarkRoomMemberAsSeen(ctx, ID)
}

func (s *loggingService) RestoreRoomMembersByRoom(
	ctx context.Context,
	roomID pulid.ID,
) (roomMember []*ent.RoomMember, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "RestoreRoomMemberByRoom",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.RestoreRoomMembersByRoom(ctx, roomID)
}
