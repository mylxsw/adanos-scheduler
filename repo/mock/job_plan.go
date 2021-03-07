package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type jobPlanRepoImpl struct {
	plans []repo.JobPlan
	idSeq int
}

func NewJobPlanRepo() repo.JobPlanRepo {
	return &jobPlanRepoImpl{
		plans: make([]repo.JobPlan, 0),
		idSeq: 0,
	}
}

func (rp *jobPlanRepoImpl) Add(ctx context.Context, jobID string, jobPlan repo.JobPlan) (jobPlanID string, err error) {
	rp.idSeq++
	jobPlan.ID = strconv.Itoa(rp.idSeq)
	jobPlan.CreatedAt = time.Now()
	jobPlan.UpdatedAt = time.Now()
	jobPlan.JobID = jobID

	rp.plans = append(rp.plans, jobPlan)
	return jobPlan.ID, nil
}

func (rp *jobPlanRepoImpl) All(ctx context.Context) (jobPlans []repo.JobPlan, err error) {
	return rp.plans, nil
}

func (rp *jobPlanRepoImpl) AllForJob(ctx context.Context, jobID string) (jobPlans []repo.JobPlan, err error) {
	jobPlans = make([]repo.JobPlan, 0)
	for _, p := range rp.plans {
		if p.JobID == jobID {
			jobPlans = append(jobPlans, p)
		}
	}

	return
}

func (rp *jobPlanRepoImpl) Get(ctx context.Context, planID string) (*repo.JobPlan, error) {
	for _, p := range rp.plans {
		if p.ID == planID {
			return &p, nil
		}
	}

	return nil, repo.ErrNotFound
}

func (rp *jobPlanRepoImpl) Update(ctx context.Context, jobPlanID string, jobPlan repo.JobPlan) error {
	for i, j := range rp.plans {
		if j.ID == jobPlanID {
			jobPlan.JobID = j.JobID
			jobPlan.CreatedAt = j.CreatedAt
			jobPlan.UpdatedAt = time.Now()
			rp.plans[i] = jobPlan
			return nil
		}
	}

	return repo.ErrNotFound
}

func (rp *jobPlanRepoImpl) Remove(ctx context.Context, jobPlanID string) error {
	deleteIndex := -1
	for i, p := range rp.plans {
		if p.ID == jobPlanID {
			deleteIndex = i
		}
	}

	if deleteIndex >= 0 {
		rp.plans = append(rp.plans[:deleteIndex], rp.plans[deleteIndex+1:]...)
	}

	return nil
}

func (rp *jobPlanRepoImpl) RemoveAllForJob(ctx context.Context, jobID string) error {
	restPlans := make([]repo.JobPlan, 0)
	for _, p := range rp.plans {
		if p.JobID == jobID {
			continue
		}

		restPlans = append(restPlans, p)
	}

	rp.plans = restPlans

	return nil
}
