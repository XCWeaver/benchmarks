package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eu_deployment/pkg/model"
	"eu_deployment/pkg/utils"

	"github.com/XCWeaver/xcweaver"
)

type UniqueIdService interface {
	UploadUniqueId(ctx context.Context, reqID int64, postType model.PostType, id string) (string, error)
}

type uniqueIdOptions struct {
	Region string `toml:"region"`
}

type uniqueIdService struct {
	xcweaver.Implements[UniqueIdService]
	xcweaver.WithConfig[uniqueIdOptions]
	composePostService xcweaver.Ref[ComposePostService]
	currentTimestamp   int64
	counter            int64
	machineID          string
	mu                 sync.Mutex
}

func (u *uniqueIdService) Init(ctx context.Context) error {
	logger := u.Logger(ctx)
	u.machineID = utils.GetMachineID()
	u.currentTimestamp = -1
	u.counter = 0
	logger.Info("unique id service running!", "machine_id", u.machineID, "region", u.Config().Region)
	return nil
}

func (u *uniqueIdService) getCounter(timestamp int64) (int64, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.currentTimestamp > timestamp {
		return 0, fmt.Errorf("timestamps are not incremental")
	}
	if u.currentTimestamp == timestamp {
		u.counter += 1
		return u.counter, nil
	} else {
		u.currentTimestamp = timestamp
		u.counter = 1
		return u.counter, nil
	}
}

func (u *uniqueIdService) UploadUniqueId(ctx context.Context, reqID int64, postType model.PostType, id string) (string, error) {
	logger := u.Logger(ctx)
	logger.Debug("entering UploadUniqueId", "req_id", reqID, "post_type", postType)

	if id != "" {
		return id, u.composePostService.Get().UploadUniqueId(ctx, reqID, id, postType, false)
	}

	timestamp := time.Now().UnixMilli() - utils.CUSTOM_EPOCH
	counter, err := u.getCounter(timestamp)
	if err != nil {
		logger.Error("error getting counter", "msg", err.Error())
		return "", err
	}
	id, err = utils.GenUniqueID(u.machineID, timestamp, counter)
	logger.Debug("uniqueId created!", "machineID", u.machineID, "timestamp", timestamp, "counter", counter)
	if err != nil {
		return "", err
	}
	return id, u.composePostService.Get().UploadUniqueId(ctx, reqID, id, postType, true)
}
