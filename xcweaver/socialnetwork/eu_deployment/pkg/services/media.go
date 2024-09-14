package services

import (
	"context"
	"eu_deployment/pkg/model"
	"fmt"

	"github.com/TiagoMalhadas/xcweaver"
)

type MediaService interface {
	UploadMedia(ctx context.Context, reqID int64, mediaTypes []string, medaIDs []int64, new bool) error
}

type mediaService struct {
	xcweaver.Implements[MediaService]
	xcweaver.WithConfig[mediaServiceOptions]
	composePostService xcweaver.Ref[ComposePostService]
}

type mediaServiceOptions struct {
	Region string `toml:"region"`
}

func (m *mediaService) Init(ctx context.Context) error {
	logger := m.Logger(ctx)
	logger.Info("media service running!", "region", m.Config().Region)
	return nil
}

func (m *mediaService) UploadMedia(ctx context.Context, reqID int64, mediaTypes []string, mediaIDs []int64, new bool) error {
	logger := m.Logger(ctx)
	logger.Debug("entering UploadMedia", "req_id", reqID, "media_types", mediaTypes, "mediaIDs", mediaIDs)
	if len(mediaTypes) != len(mediaIDs) {
		errMsg := "the lengths of media_id list and media_type list are not equal"
		logger.Error(errMsg, "num_media_types", len(mediaTypes), "num_media_ids", len(mediaIDs))
		return fmt.Errorf(errMsg)
	}
	var medias []model.Media
	for i := range mediaIDs {
		medias = append(medias, model.Media{
			MediaID:   mediaIDs[i],
			MediaType: mediaTypes[i],
		})
	}
	return m.composePostService.Get().UploadMedia(ctx, reqID, medias, new)
}
