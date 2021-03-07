package repo

import (
	"context"
	"time"
)

// Job 定义了一个任务
type Job struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Command     string    `json:"command"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JobRepo job 管理仓库
type JobRepo interface {
	// Add 创建一个 Job
	Add(ctx context.Context, job Job) (jobID string, err error)
	// All 返回所有的 Job
	All(ctx context.Context) (jobs []Job, err error)
	// Update 更新 job
	Update(ctx context.Context, jobID string, job Job) error
	// Remove 删除 Job
	Remove(ctx context.Context, jobID string) error
	// GetByID 更过 ID 查询 job
	GetByID(ctx context.Context, jobID string) (job *Job, err error)
}

// JobPlanStatus 任务计划状态
type JobPlanStatus string

const (
	// JobPlanStatusEnabled 启用该计划任务
	JobPlanStatusEnabled JobPlanStatus = "enabled"
	// JobPlanStatusEnabled 禁用该计划任务
	JobPlanStatusDisabled JobPlanStatus = "disabled"
)

// JobExecuteMode 任务执行模式
type JobExecuteMode string

const (
	// JobExecuteModeAll 在所有的匹配节点上同时执行
	JobExecuteModeAll JobExecuteMode = "all"
	// JobExecuteModeSingle 只在一个节点上执行
	JobExecuteModeSingle JobExecuteMode = "single"
)

// JobPlan 任务执行计划
type JobPlan struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Plan          string         `json:"plan"`
	JobID         string         `json:"job_id"`
	NodeSelectors Labels         `json:"node_selectors"`
	ExecuteMode   JobExecuteMode `json:"execute_mode"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// JobPlanRepo 任务执行计划仓库
type JobPlanRepo interface {
	// Add 新增一个计划任务
	Add(ctx context.Context, jobID string, jobPlan JobPlan) (jobPlanID string, err error)
	// All 返回所有的计划任务
	All(ctx context.Context) (jobPlans []JobPlan, err error)
	// Get 返回某个计划任务
	Get(ctx context.Context, planID string) (*JobPlan, error)
	// AllForJob 返回任务关联的所有计划任务
	AllForJob(ctx context.Context, jobID string) (jobPlans []JobPlan, err error)
	// Update 更新计划任务
	Update(ctx context.Context, jobPlanID string, jobPlan JobPlan) error
	// Remove 删除一个计划任务
	Remove(ctx context.Context, jobPlanID string) error
	// RemoveAllForJob 移除任务关联的所有计划任务
	RemoveAllForJob(ctx context.Context, jobID string) error
}
