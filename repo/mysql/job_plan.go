package mysql

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type jobPlanRepoImpl struct {
}

func NewJobPlanRepo() repo.JobPlanRepo {
	return &jobPlanRepoImpl{}
}

func (rp *jobPlanRepoImpl) Add(ctx context.Context, jobID string, jobPlan repo.JobPlan) (jobPlanID string, err error) {
	panic("implement me")
}

func (rp *jobPlanRepoImpl) All(ctx context.Context) (jobPlans []repo.JobPlan, err error) {
	panic("implement me")
}

func (rp *jobPlanRepoImpl) AllForJob(ctx context.Context, jobID string) (jobPlans []repo.JobPlan, err error) {
	panic("implement me")
}

func (rp *jobPlanRepoImpl) Get(ctx context.Context, planID string) (*repo.JobPlan, error) {
	panic("implement me")
}

func (rp *jobPlanRepoImpl) Update(ctx context.Context, jobPlanID string, jobPlan repo.JobPlan) error {
	panic("implement me")
}

func (rp *jobPlanRepoImpl) Remove(ctx context.Context, jobPlanID string) error {
	panic("implement me")
}

func (rp *jobPlanRepoImpl) RemoveAllForJob(ctx context.Context, jobID string) error {
	panic("implement me")
}
