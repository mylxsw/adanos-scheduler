package repo

import (
	"context"
	"time"
)

// JobHistory 任务执行历史记录
type JobHistory struct {
	ID         string
	JobID      string
	JobPlanID  string
	NodeID     string
	Output     string
	ReturnCode int
	Error      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type JobHistoryRepo interface {
	// Add 新增任务执行历史记录
	Add(ctx context.Context, history JobHistory) (historyID string, err error)
	// FindByJobID 根据 Job 查询所有的执行记录
	FindByJobID(ctx context.Context, jobID string, limit int) ([]JobHistory, error)
	// FindByJobPlanID 根据执行计划 id 查询所有执行记录
	FindByJobPlanID(ctx context.Context, planID string, limit int) ([]JobHistory, error)
	// FindByNodeID 根据 Node 查询所有的执行记录
	FindByNodeID(ctx context.Context, nodeID string, limit int) ([]JobHistory, error)
	// FindAll 查询所有任务执行历史记录
	FindAll(ctx context.Context, limit int) ([]JobHistory, error)
	// RemoveByID 根据id删除任务执行历史记录
	RemoveByID(ctx context.Context, historyID string) error
	// RemoveAllForJob 根据 job 删除所有执行记录
	RemoveAllForJob(ctx context.Context, jobID string) error
	// RemoveAllForJobPlan 根据计划任务 id 删除所有执行记录
	RemoveAllForJobPlan(ctx context.Context, planID string) error
}
