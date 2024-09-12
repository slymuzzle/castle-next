package roommembers

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type subscriptionsLogging struct {
	logger log.Logger
	Subscriptions
}

func NewSubscriptionsLogging(logger log.Logger, s Subscriptions) Subscriptions {
	return &subscriptionsLogging{logger, s}
}

func (s *subscriptionsLogging) PublishRoomMemberCreatedEvent(
	ctx context.Context,
	userID pulid.ID,
	roomMemberID pulid.ID,
) (subject string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "PublishRoomMemberCreatedEvent",
			"subject", subject,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.PublishRoomMemberCreatedEvent(ctx, userID, roomMemberID)
}

func (s *subscriptionsLogging) SubscribeToRoomMemberCreatedEvent(
	ctx context.Context,
) (ch <-chan *ent.RoomMemberEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomMemberCreatedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.SubscribeToRoomMemberCreatedEvent(ctx)
}

func (s *subscriptionsLogging) PublishRoomMemberUpdatedEvent(
	ctx context.Context,
	userID pulid.ID,
	roomMemberID pulid.ID,
) (subject string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "PublishRoomMemberUpdatedEvent",
			"subject", subject,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.PublishRoomMemberUpdatedEvent(ctx, userID, roomMemberID)
}

func (s *subscriptionsLogging) SubscribeToRoomMemberUpdatedEvent(
	ctx context.Context,
) (ch <-chan *ent.RoomMemberEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomMemberUpdatedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.SubscribeToRoomMemberUpdatedEvent(ctx)
}

func (s *subscriptionsLogging) PublishRoomMemberDeletedEvent(
	ctx context.Context,
	roomMemberID pulid.ID,
) (subject string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "PublishRoomMemberDeletedEvent",
			"subject", subject,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.PublishRoomMemberDeletedEvent(ctx, roomMemberID)
}

func (s *subscriptionsLogging) SubscribeToRoomMemberDeletedEvent(
	ctx context.Context,
) (ch <-chan pulid.ID, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomMemberDeletedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.SubscribeToRoomMemberDeletedEvent(ctx)
}
