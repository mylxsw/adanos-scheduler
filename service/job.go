package service

import (
	"context"
	"fmt"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
)

// JobWithPlan 任务和执行计划的组合体
type JobWithPlan struct {
	Job   repo.Job       `json:"job"`
	Plans []repo.JobPlan `json:"plans"`
}

// JobService 任务管理服务接口
type JobService interface {
	// AllJobs 查询所有的 jobs
	AllJobs(ctx context.Context) ([]repo.Job, error)
	// AllPlansForJob 查询任务关联的所有执行计划
	JobWithPlans(ctx context.Context, jobID string) (*JobWithPlan, error)
	// AllPlans 查询所有的计划任务
	AllPlans(ctx context.Context) ([]repo.JobPlan, error)
	// JobPlan 查询某个计划任务详情
	JobPlan(ctx context.Context, jobID string, planID string) (*repo.JobPlan, error)

	// CreateJobWithPlan 创建一个带有计划的 job
	CreateJobWithPlan(ctx context.Context, job JobWithPlan) (jobID string, err error)
	// UpdateJob 更新一个 job
	UpdateJob(ctx context.Context, jobID string, job repo.Job) error
	// CreateJobPlan 创建任务执行计划
	CreateJobPlan(ctx context.Context, jobID string, plan repo.JobPlan) (jobPlanID string, err error)
	// UpdateJobPlan 更新任务执行计划
	UpdateJobPlan(ctx context.Context, jobID string, planID string, plan repo.JobPlan) error
	// RemoveJob 移除一个 Job
	RemoveJob(ctx context.Context, jobID string) error
	// RemoveJobPlan 移除一个执行计划
	RemoveJobPlan(ctx context.Context, jobID string, jobPlanID string) error
}

// jobService 任务管理服务实现
type jobService struct {
	cc          container.Container
	jobRepo     repo.JobRepo     `autowire:"@"`
	jobPlanRepo repo.JobPlanRepo `autowire:"@"`
	nodeRepo    repo.NodeRepo    `autowire:"@"`
}

// NewJobService 创建任务管理服务
func NewJobService(cc container.Container) JobService {
	s := &jobService{cc: cc}
	_ = cc.AutoWire(s)
	return s
}

func (srv *jobService) AllJobs(ctx context.Context) ([]repo.Job, error) {
	return srv.jobRepo.All(ctx)
}

func (srv *jobService) JobWithPlans(ctx context.Context, jobID string) (*JobWithPlan, error) {
	job, err := srv.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get job by id failed: %w", err)
	}

	plans, err := srv.jobPlanRepo.AllForJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get job plans failed: %w", err)
	}

	return &JobWithPlan{
		Job:   *job,
		Plans: plans,
	}, nil
}

func (srv *jobService) AllPlans(ctx context.Context) ([]repo.JobPlan, error) {
	return srv.jobPlanRepo.All(ctx)
}

func (srv *jobService) JobPlan(ctx context.Context, jobID string, planID string) (*repo.JobPlan, error) {
	plan, err := srv.jobPlanRepo.Get(ctx, planID)
	if err != nil {
		return nil, err
	}

	if plan.JobID != jobID {
		return nil, repo.ErrNotFound
	}

	return plan, nil
}

func (srv *jobService) CreateJobWithPlan(ctx context.Context, job JobWithPlan) (jobID string, err error) {
	jobID, err = srv.jobRepo.Add(ctx, job.Job)
	if err != nil {
		return "", fmt.Errorf("create job failed: %w", err)
	}

	if len(job.Plans) == 0 {
		return
	}

	for _, plan := range job.Plans {
		plan.JobID = jobID
		if _, err = srv.jobPlanRepo.Add(ctx, jobID, plan); err != nil {
			log.With(plan).Errorf("create job succeed, but create plan failed: %w", err)
		}
	}

	return
}

func (srv *jobService) UpdateJob(ctx context.Context, jobID string, job repo.Job) error {
	return srv.jobRepo.Update(ctx, jobID, job)
}

func (srv *jobService) CreateJobPlan(ctx context.Context, jobID string, plan repo.JobPlan) (jobPlanID string, err error) {
	_, err = srv.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return "", err
	}

	return srv.jobPlanRepo.Add(ctx, jobID, plan)
}

func (srv *jobService) UpdateJobPlan(ctx context.Context, jobID string, planID string, plan repo.JobPlan) error {
	if _, err := srv.jobRepo.GetByID(ctx, jobID); err != nil {
		return err
	}

	return srv.jobPlanRepo.Update(ctx, planID, plan)
}

func (srv *jobService) RemoveJob(ctx context.Context, jobID string) error {
	if err := srv.jobPlanRepo.RemoveAllForJob(ctx, jobID); err != nil {
		return fmt.Errorf("remove associated plans for job failed: %w", err)
	}
	return srv.jobRepo.Remove(ctx, jobID)
}

func (srv *jobService) RemoveJobPlan(ctx context.Context, jobID string, jobPlanID string) error {
	if _, err := srv.jobRepo.GetByID(ctx, jobID); err != nil {
		return err
	}

	return srv.jobPlanRepo.Remove(ctx, jobPlanID)
}
