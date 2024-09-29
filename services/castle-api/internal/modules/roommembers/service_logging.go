package roommembers

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

func (s *serviceLogging) IncrementUnreadMessagesCount(
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

func (s *serviceLogging) DeleteRoomMember(
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

func (s *serviceLogging) MarkRoomMemberAsSeen(
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

func (s *serviceLogging) RestoreRoomMembersByRoom(
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
