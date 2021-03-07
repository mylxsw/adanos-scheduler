package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/coll"
)

type jobHistoryRepoImpl struct {
	histories []repo.JobHistory
	idSeq     int
}

func NewJobHistoryRepo() repo.JobHistoryRepo {
	return &jobHistoryRepoImpl{
		histories: make([]repo.JobHistory, 0),
		idSeq:     0,
	}
}

func (rp *jobHistoryRepoImpl) Add(ctx context.Context, history repo.JobHistory) (historyID string, err error) {
	rp.idSeq++
	history.ID = strconv.Itoa(rp.idSeq)
	history.CreatedAt = time.Now()
	rp.histories = append(rp.histories, history)
	return history.ID, err
}

func (rp *jobHistoryRepoImpl) FindByJobID(ctx context.Context, jobID string, limit int) (histories []repo.JobHistory, err error) {
	err = coll.MustNew(rp.histories).Filter(func(his repo.JobHistory) bool { return his.JobID == jobID }).All(&histories)
	return histories[:limit], err
}

func (rp *jobHistoryRepoImpl) FindByJobPlanID(ctx context.Context, planID string, limit int) (histories []repo.JobHistory, err error) {
	err = coll.MustNew(rp.histories).Filter(func(his repo.JobHistory) bool { return his.JobPlanID == planID }).All(&histories)
	return histories[:limit], err
}

func (rp *jobHistoryRepoImpl) FindByNodeID(ctx context.Context, nodeID string, limit int) (histories []repo.JobHistory, err error) {
	err = coll.MustNew(rp.histories).Filter(func(his repo.JobHistory) bool { return his.NodeID == nodeID }).All(&histories)
	return histories[:limit], err
}

func (rp *jobHistoryRepoImpl) FindAll(ctx context.Context, limit int) ([]repo.JobHistory, error) {
	return rp.histories[:limit], nil
}

func (rp *jobHistoryRepoImpl) RemoveByID(ctx context.Context, historyID string) error {
	return coll.Filter(rp.histories, &rp.histories, func(his repo.JobHistory) bool { return his.ID != historyID })
}

func (rp *jobHistoryRepoImpl) RemoveAllForJob(ctx context.Context, jobID string) error {
	return coll.Filter(rp.histories, &rp.histories, func(his repo.JobHistory) bool { return his.JobID == jobID })
}

func (rp *jobHistoryRepoImpl) RemoveAllForJobPlan(ctx context.Context, planID string) error {
	return coll.Filter(rp.histories, &rp.histories, func(his repo.JobHistory) bool { return his.JobPlanID == planID })
}
