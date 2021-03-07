package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type jobRepoImpl struct {
	jobs  []repo.Job
	idSeq int
}

func NewJobRepo() repo.JobRepo {
	return &jobRepoImpl{
		jobs:  make([]repo.Job, 0),
		idSeq: 0,
	}
}

func (rp *jobRepoImpl) All(ctx context.Context) (jobs []repo.Job, err error) {
	return rp.jobs, nil
}

func (rp *jobRepoImpl) Add(ctx context.Context, job repo.Job) (jobID string, err error) {
	rp.idSeq++
	job.ID = strconv.Itoa(rp.idSeq)
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	rp.jobs = append(rp.jobs, job)
	return job.ID, err
}

func (rp *jobRepoImpl) Update(ctx context.Context, jobID string, job repo.Job) error {
	for i, j := range rp.jobs {
		if j.ID == jobID {
			job.ID = jobID
			job.CreatedAt = rp.jobs[i].CreatedAt
			job.UpdatedAt = time.Now()
			rp.jobs[i] = job
			return nil
		}
	}

	return repo.ErrNotFound
}

func (rp *jobRepoImpl) Remove(ctx context.Context, jobID string) error {
	deleteIndex := -1
	for i, j := range rp.jobs {
		if j.ID == jobID {
			deleteIndex = i
			break
		}
	}

	if deleteIndex >= 0 {
		rp.jobs = append(rp.jobs[:deleteIndex], rp.jobs[deleteIndex+1:]...)
	}
	return nil
}

func (rp *jobRepoImpl) GetByID(ctx context.Context, jobID string) (job *repo.Job, err error) {
	for _, j := range rp.jobs {
		if j.ID == jobID {
			return &j, nil
		}
	}

	return nil, repo.ErrNotFound
}
