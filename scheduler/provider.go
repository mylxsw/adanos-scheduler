package scheduler

import (
	"github.com/mylxsw/adanos-scheduler/scheduler/task"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/cron"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (p Provider) Register(app container.Container) {
	app.MustSingletonOverride(task.NewJobPlanLoaderTask)
}

func (p Provider) Boot(app infra.Glacier) {
	app.Cron(func(cr cron.Manager, cc container.Container) error {
		return cc.ResolveWithError(func(jobPlanLoaderTask *task.JobPlanLoaderTask) error {
			// 创建计划任务加载任务
			cc.Must(jobPlanLoaderTask.Handle(cc))
			cc.Must(cr.Add("job-plan-loader", "@every 5s", jobPlanLoaderTask))

			return nil
		})
	})
}
