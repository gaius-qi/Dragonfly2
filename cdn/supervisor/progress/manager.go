/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//go:generate mockgen -destination ../mocks/progress/mock_progress_manager.go -package progress d7y.io/dragonfly/v2/cdn/supervisor/progress Manager

package progress

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"

	"d7y.io/dragonfly/v2/cdn/constants"
	"d7y.io/dragonfly/v2/cdn/supervisor/task"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/synclock"
)

// Manager as an interface defines all operations about seed progress
type Manager interface {

	// WatchSeedProgress watch task seed progress
	WatchSeedProgress(ctx context.Context, clientAddr string, taskID string) (<-chan *task.PieceInfo, error)

	// PublishPiece publish piece seed
	PublishPiece(ctx context.Context, taskID string, piece *task.PieceInfo) error

	// PublishTask publish task seed
	PublishTask(ctx context.Context, taskID string, task *task.SeedTask) error
}

var _ Manager = (*manager)(nil)

type manager struct {
	mu               *synclock.LockerPool
	taskManager      task.Manager
	seedTaskSubjects sync.Map
}

func NewManager(taskManager task.Manager) (Manager, error) {
	return newManager(taskManager)
}

func newManager(taskManager task.Manager) (*manager, error) {
	return &manager{
		mu:               synclock.NewLockerPool(),
		taskManager:      taskManager,
		seedTaskSubjects: sync.Map{},
	}, nil
}

func (pm *manager) WatchSeedProgress(ctx context.Context, clientAddr string, taskID string) (<-chan *task.PieceInfo, error) {
	pm.mu.Lock(taskID, false)
	defer pm.mu.UnLock(taskID, false)
	span := trace.SpanFromContext(ctx)
	span.AddEvent(constants.EventWatchSeedProgress)
	seedTask, err := pm.taskManager.Get(taskID)
	if err != nil {
		return nil, err
	}
	pieces, err := pm.taskManager.GetProgress(taskID)
	if err != nil {
		return nil, err
	}
	if seedTask.IsDone() {
		pieceChan := make(chan *task.PieceInfo)
		go func(pieceChan chan *task.PieceInfo) {
			defer func() {
				logger.Debugf("subscriber %s starts watching task %s seed progress", clientAddr, taskID)
				close(pieceChan)
			}()
			for _, piece := range pieces {
				logger.Debugf("notifies subscriber %s about %d piece info of taskID %s", clientAddr, piece.PieceNum, taskID)
				pieceChan <- piece
			}
		}(pieceChan)
		return pieceChan, nil
	}
	var progressPublisher, _ = pm.seedTaskSubjects.LoadOrStore(taskID, newProgressPublisher(taskID))
	observer := newProgressSubscriber(ctx, clientAddr, seedTask.ID, pieces)
	progressPublisher.(*publisher).AddSubscriber(observer)
	return observer.Receiver(), nil
}

func (pm *manager) PublishPiece(ctx context.Context, taskID string, record *task.PieceInfo) (err error) {
	pm.mu.Lock(taskID, false)
	defer pm.mu.UnLock(taskID, false)
	span := trace.SpanFromContext(ctx)
	jsonRecord, err := json.Marshal(record)
	if err != nil {
		return errors.Wrapf(err, "json marshal piece record: %#v", record)
	}
	span.AddEvent(constants.EventPublishPiece, trace.WithAttributes(constants.AttributeSeedPiece.String(string(jsonRecord))))
	logger.Debugf("publish task %s seed piece record: %s", taskID, jsonRecord)
	var progressPublisher, ok = pm.seedTaskSubjects.Load(taskID)
	if ok {
		progressPublisher.(*publisher).NotifySubscribers(record)
	}
	return pm.taskManager.UpdateProgress(taskID, record)
}

func (pm *manager) PublishTask(ctx context.Context, taskID string, seedTask *task.SeedTask) error {
	jsonTask, err := json.Marshal(seedTask)
	if err != nil {
		return errors.Wrapf(err, "json marshal seedTask: %#v", seedTask)
	}
	logger.Debugf("publish task %s seed piece record: %s", taskID, jsonTask)
	pm.mu.Lock(taskID, false)
	defer pm.mu.UnLock(taskID, false)
	span := trace.SpanFromContext(ctx)
	recordBytes, _ := json.Marshal(seedTask)
	span.AddEvent(constants.EventPublishTask, trace.WithAttributes(constants.AttributeSeedTask.String(string(recordBytes))))
	// task status must be updated before removing the subscribers
	if err := pm.taskManager.Update(taskID, seedTask); err != nil {
		return err
	}
	if progressPublisher, ok := pm.seedTaskSubjects.Load(taskID); ok {
		progressPublisher.(*publisher).RemoveAllSubscribers()
		pm.seedTaskSubjects.Delete(taskID)
	}
	return nil
}
