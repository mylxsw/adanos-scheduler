package mysql

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type jobHistoryRepoImpl struct {
}

func NewJobHistoryRepo() repo.JobHistoryRepo {
	return &jobHistoryRepoImpl{}
}

func (rp *jobHistoryRepoImpl) Add(ctx context.Context, history repo.JobHistory) (historyID string, err error) {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) FindByJobID(ctx context.Context, jobID string, limit int) ([]repo.JobHistory, error) {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) FindByJobPlanID(ctx context.Context, planID string, limit int) ([]repo.JobHistory, error) {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) FindByNodeID(ctx context.Context, nodeID string, limit int) ([]repo.JobHistory, error) {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) FindAll(ctx context.Context, limit int) ([]repo.JobHistory, error) {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) RemoveByID(ctx context.Context, historyID string) error {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) RemoveAllForJob(ctx context.Context, jobID string) error {
	panic("implement me")
}

func (rp *jobHistoryRepoImpl) RemoveAllForJobPlan(ctx context.Context, planID string) error {
	panic("implement me")
}
