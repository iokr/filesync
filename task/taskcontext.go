package task

import (
	"sync"
	"context"
)

var (
	taskContext 	map[string]context.CancelFunc
	ctxMutex		sync.Mutex
)

func init() {
	taskContext = make(map[string]context.CancelFunc)
}

func AddTaskContext(taskID string, ctxCancel context.CancelFunc) {
	ctxMutex.Lock()
	defer ctxMutex.Unlock()
	taskContext[taskID] = ctxCancel
}

func DeleteTaskContext(taskID string) {
	ctxMutex.Lock()
	defer ctxMutex.Unlock()
	delete(taskContext, taskID)
}

func FindTaskContext(taskID string) (context.CancelFunc, bool) {
	ctxMutex.Lock()
	defer ctxMutex.Unlock()

	taskCtx, ok := taskContext[taskID]
	return taskCtx, ok
}

func IsExistsTaskContext(taskID string) bool {
	ctxMutex.Lock()
	defer ctxMutex.Unlock()

	_, ok := taskContext[taskID]
	return ok
}