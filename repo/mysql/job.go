package mysql

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type jobRepoImpl struct {
}

func NewJobRepo() repo.JobRepo {
	return &jobRepoImpl{}
}

func (rp *jobRepoImpl) All(ctx context.Context) (jobs []repo.Job, err error) {
	panic("implement me")
}

func (rp *jobRepoImpl) Add(ctx context.Context, job repo.Job) (jobID string, err error) {
	panic("implement me")
}

func (rp *jobRepoImpl) Update(ctx context.Context, jobID string, job repo.Job) error {
	panic("implement me")
}

func (rp *jobRepoImpl) Remove(ctx context.Context, jobID string) error {
	panic("implement me")
}

func (rp *jobRepoImpl) GetByID(ctx context.Context, jobID string) (job *repo.Job, err error) {
	panic("implement me")
}
