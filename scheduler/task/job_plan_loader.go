package task

import (
	"context"
	"fmt"

	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
)

// JobPlanLoaderTask 用于加载计划任务并分配个相应的执行器，更新执行计划
type JobPlanLoaderTask struct {
	jobSrv  service.JobService
	nodeSrv service.NodeService
	credSrv service.CredentialService
}

// NewJobPlanLoaderTask 创建 JobPlanLoaderTask
func NewJobPlanLoaderTask(jobSrv service.JobService, nodeSrv service.NodeService, credSrv service.CredentialService) *JobPlanLoaderTask {
	return &JobPlanLoaderTask{
		jobSrv:  jobSrv,
		nodeSrv: nodeSrv,
		credSrv: credSrv,
	}
}

func (job *JobPlanLoaderTask) Handle(cc infra.Resolver) error {
	return cc.ResolveWithError(func(ctx context.Context) error {
		plans, err := job.jobSrv.AllPlans(ctx)
		if err != nil {
			return fmt.Errorf("query job plans failed: %w", err)
		}

		for _, p := range plans {
			selectedNodes, err := job.nodeSrv.SelectNodes(ctx, p.NodeSelectors)
			if err != nil {
				log.With(p).Errorf("select nodes for job plan failed: %w", err)
				continue
			}

			log.WithFields(log.Fields{
				"plan":  p,
				"nodes": selectedNodes,
			}).Warningf("plan %s loaded", p.Name)
		}

		return nil
	})
}
