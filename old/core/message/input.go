package message

import (
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/internal/platform/media"
)

type (
	MessageOrderInput = ent.MessageOrder
	MessageWhereInput = ent.MessageWhereInput
)

type (
	UploadAttachmentsFn func(*ent.Message) ([]*media.UploadInfo, error)
	UploadVoiceFn       func(*ent.Message) (*media.UploadInfo, error)
)

type CreateMessageInput struct {
	RoomID              pulid.ID
	UserID              pulid.ID
	ReplyToID           *pulid.ID
	Content             string
	UploadAttachmentsFn UploadAttachmentsFn
	UploadVoiceFn       UploadVoiceFn
}

func (i *CreateMessageInput) Mutate(m *ent.MessageMutation) {
	m.SetRoomID(i.RoomID)
	m.SetUserID(i.UserID)
	if v := i.ReplyToID; v != nil {
		m.SetReplyToID(*v)
	}
	m.SetContent(i.Content)
}

type UpdateMessageInput struct {
	RoomID              *pulid.ID
	UserID              *pulid.ID
	ReplyToID           *pulid.ID
	Content             *string
	UploadAttachmentsFn UploadAttachmentsFn
	UploadVoiceFn       UploadVoiceFn
}

func (i *UpdateMessageInput) Mutate(m *ent.MessageMutation) {
	if v := i.RoomID; v != nil {
		m.SetRoomID(*v)
	}
	if v := i.UserID; v != nil {
		m.SetUserID(*v)
	}
	if v := i.ReplyToID; v != nil {
		m.SetReplyToID(*v)
	}
	if v := i.Content; v != nil {
		m.SetContent(*v)
	}
}
