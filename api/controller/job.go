package controller

import (
	"context"
	"errors"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
)

type JobController struct {
	cc container.Container
}

func NewJobController(cc container.Container) web.Controller {
	ctl := &JobController{cc: cc}
	return ctl
}

func (ctl JobController) Register(router *web.Router) {
	router.Group("/jobs", func(router *web.Router) {
		router.Get("/", ctl.Jobs)

		router.Post("/", ctl.CreateJob)
		router.Get("/{job_id}/", ctl.JobWithPlans)
		router.Put("/{job_id}/", ctl.UpdateJob)
		router.Delete("/{job_id}", ctl.RemoveJob)

		router.Post("/{job_id}/plan/", ctl.CreateJobPlan)
		router.Get("/{job_id}/plan/{plan_id}/", ctl.GetJobPlan)
		router.Put("/{job_id}/plan/{plan_id}/", ctl.UpdateJobPlan)
		router.Delete("/{job_id}/plan/{plan_id}/", ctl.RemoveJobPlan)
	})
}

// Jobs 返回所有的 jobs
func (ctl JobController) Jobs(ctx context.Context, srv service.JobService) ([]repo.Job, error) {
	return srv.AllJobs(ctx)
}

// JobWithPlans 返回请求的 job 以及相关联的计划任务
func (ctl JobController) JobWithPlans(ctx context.Context, req web.Request, srv service.JobService) (*service.JobWithPlan, error) {
	jobID := req.PathVar("job_id")
	return srv.JobWithPlans(ctx, jobID)
}

// JobWithPlansCreateReq job 和 计划任务一起创建的请求
type JobWithPlansCreateReq struct {
	JobReq
	Plans []JobPlanReq `json:"plans"`
}

// JobPlanReq 计划任务更新请求
type JobPlanReq struct {
	Name          string            `json:"name"`
	Plan          string            `json:"plan"`
	NodeSelectors map[string]string `json:"node_selectors"`
	ExecuteMode   string            `json:"execute_mode"`
}

func (pl JobPlanReq) Validate(request web.Request) error {
	// TODO 表单校验
	if pl.Name == "" {
		return errors.New("invalid job name")
	}

	return nil
}

func (pl JobPlanReq) Transform() repo.JobPlan {
	return repo.JobPlan{
		Name:          pl.Name,
		Plan:          pl.Plan,
		NodeSelectors: pl.NodeSelectors,
		ExecuteMode:   repo.JobExecuteMode(pl.ExecuteMode),
	}
}

func (jobWithPlanReq JobWithPlansCreateReq) Validate(req web.Request) error {
	if err := jobWithPlanReq.JobReq.Validate(req); err != nil {
		return err
	}

	for _, p := range jobWithPlanReq.Plans {
		if err := p.Validate(req); err != nil {
			return err
		}
	}

	return nil
}

// Transform 将 job 请求转换为服务更新对象
func (jobWithPlanReq JobWithPlansCreateReq) Transform() service.JobWithPlan {
	plans := make([]repo.JobPlan, 0)
	for _, pl := range jobWithPlanReq.Plans {
		plans = append(plans, pl.Transform())
	}

	return service.JobWithPlan{
		Job:   jobWithPlanReq.JobReq.Transform(),
		Plans: plans,
	}
}

// CreateJob 创建一个 job
func (ctl JobController) CreateJob(ctx context.Context, req web.Request, srv service.JobService) (IDResponse, error) {
	var jobWithPlanReq JobWithPlansCreateReq
	if err := req.Unmarshal(&jobWithPlanReq); err != nil {
		return IDResponse{}, err
	}

	req.Validate(jobWithPlanReq, true)

	jobID, err := srv.CreateJobWithPlan(ctx, jobWithPlanReq.Transform())
	return IDResponse{ID: jobID}, err
}

// JobReq job 更新请求
type JobReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
}

func (jq JobReq) Validate(request web.Request) error {
	// TODO 表单校验
	if jq.Name == "" {
		return errors.New("invalid job name")
	}

	if jq.Command == "" {
		return errors.New("invalid command")
	}

	return nil
}

func (jq JobReq) Transform() repo.Job {
	return repo.Job{
		Name:        jq.Name,
		Description: jq.Description,
		Command:     jq.Command,
	}
}

// UpdateJob 更新 job
func (ctl JobController) UpdateJob(ctx context.Context, req web.Request, srv service.JobService) error {
	var jobReq JobReq
	if err := req.Unmarshal(&jobReq); err != nil {
		return err
	}

	jobID := req.PathVar("job_id")
	return srv.UpdateJob(ctx, jobID, jobReq.Transform())
}

// RemoveJob 移除一个 job
func (ctl JobController) RemoveJob(ctx context.Context, req web.Request, srv service.JobService) error {
	jobID := req.PathVar("job_id")
	return srv.RemoveJob(ctx, jobID)
}

// CreateJobPlan 创建一个计划任务
func (ctl JobController) CreateJobPlan(ctx context.Context, req web.Request, srv service.JobService) (IDResponse, error) {
	var planReq JobPlanReq
	if err := req.Unmarshal(&planReq); err != nil {
		return IDResponse{}, err
	}

	req.Validate(planReq, true)

	jobID := req.PathVar("job_id")
	planID, err := srv.CreateJobPlan(ctx, jobID, planReq.Transform())
	return IDResponse{ID: planID}, err
}

// UpdateJobPlan 更新计划任务
func (ctl JobController) UpdateJobPlan(ctx context.Context, req web.Request, srv service.JobService) error {
	var planReq JobPlanReq
	if err := req.Unmarshal(&planReq); err != nil {
		return err
	}

	req.Validate(planReq, true)

	jobID := req.PathVar("job_id")
	planID := req.PathVar("plan_id")

	return srv.UpdateJobPlan(ctx, jobID, planID, planReq.Transform())
}

// RemoveJobPlan 删除计划任务
func (ctl JobController) RemoveJobPlan(ctx context.Context, req web.Request, srv service.JobService) error {
	jobID := req.PathVar("job_id")
	planID := req.PathVar("plan_id")

	return srv.RemoveJobPlan(ctx, jobID, planID)
}

// GetJobPlan 查询某个计划任务信息
func (ctl JobController) GetJobPlan(ctx context.Context, req web.Request, srv service.JobService) (*repo.JobPlan, error) {
	jobID := req.PathVar("job_id")
	planID := req.PathVar("plan_id")

	return srv.JobPlan(ctx, jobID, planID)
}
