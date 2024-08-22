package file

import (
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

	"github.com/99designs/gqlgen/graphql"
)

type (
	FileWhereInput = ent.FileWhereInput
	FileOrderInput = ent.FileOrder
)

type (
	UploadFile = graphql.Upload
)

type UploadInfo struct {
	ID          pulid.ID
	Filename    string
	ContentType string
	Size        int64
	Bucket      string
	Location    string
	Path        string
}

type CreateFileInput struct {
	UserID    pulid.ID
	ContactID pulid.ID
	RoomID    *pulid.ID
}

func (c *CreateFileInput) Mutate(m *ent.FileMutation) {
}

type UpdateFileInput struct {
	UserID    *pulid.ID
	ContactID *pulid.ID
	RoomID    *pulid.ID
}

func (c *UpdateFileInput) Mutate(m *ent.FileMutation) {
}
