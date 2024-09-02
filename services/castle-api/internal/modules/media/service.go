package media

import (
	"context"
	"fmt"
	"path/filepath"

	"journeyhub/ent"
	"journeyhub/ent/file"
	"journeyhub/ent/messageattachment"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"journeyhub/internal/platform/config"

	"github.com/gabriel-vasile/mimetype"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/sync/errgroup"
)

type UploadInfo struct {
	ID          pulid.ID
	Type        messageattachment.Type
	Name        string
	Filename    string
	ContentType string
	Size        int64
	Bucket      string
	Location    string
	Path        string
}

type Service interface {
	UploadMessageFiles(
		ctx context.Context,
		prefix string,
		files []*model.UploadMessageFile,
	) ([]*UploadInfo, error)
	UploadFile(
		ctx context.Context,
		prefix string,
		file *graphql.Upload,
	) (*UploadInfo, error)
	Config() config.S3Config
}

type service struct {
	config         config.S3Config
	minioClient    *minio.Client
	uploadIDPrefix string
}

func NewService(config config.S3Config) (Service, error) {
	minioClient, err := minio.New(
		config.Host,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.AccessKey,
				config.SecretKey,
				"",
			),
			Secure: config.Ssl,
		})
	if err != nil {
		return nil, err
	}

	uploadIDPrefix, err := ent.TableToPrefix(file.Table)
	if err != nil {
		return nil, err
	}

	return &service{
		config:         config,
		minioClient:    minioClient,
		uploadIDPrefix: uploadIDPrefix,
	}, nil
}

func (s *service) UploadMessageFiles(
	ctx context.Context,
	prefix string,
	files []*model.UploadMessageFile,
) ([]*UploadInfo, error) {
	uploadInfoCh := make(chan *UploadInfo, len(files))

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(10)

	for _, file := range files {
		eg.Go(func() error {
			uploadInfo, err := s.UploadFile(egCtx, prefix, &file.File)
			uploadInfo.Type = file.Type
			if err != nil {
				return err
			}
			uploadInfoCh <- uploadInfo
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	close(uploadInfoCh)

	uploadInfo := make([]*UploadInfo, 0, len(files))
	for chValue := range uploadInfoCh {
		uploadInfo = append(uploadInfo, chValue)
	}

	return uploadInfo, nil
}

func (s *service) UploadFile(
	ctx context.Context,
	prefix string,
	file *graphql.Upload,
) (*UploadInfo, error) {
	uploadID := pulid.MustNew(s.uploadIDPrefix)

	objectName := fmt.Sprintf(
		"%s/%s%s",
		prefix,
		uploadID,
		filepath.Ext(file.Filename),
	)

	mtype, err := mimetype.DetectReader(file.File)
	if err != nil {
		return nil, err
	}

	uploadInfo, err := s.minioClient.PutObject(
		ctx,
		s.config.Bucket,
		objectName,
		file.File,
		file.Size,
		minio.PutObjectOptions{
			ContentType: mtype.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &UploadInfo{
		ID:          uploadID,
		Filename:    file.Filename,
		ContentType: mtype.String(),
		Size:        uploadInfo.Size,
		Location:    uploadInfo.Location,
		Bucket:      uploadInfo.Bucket,
		Path:        uploadInfo.Key,
	}, nil
}

func (s *service) Config() config.S3Config {
	return s.config
}
