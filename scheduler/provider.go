package scheduler

import (
	"github.com/mylxsw/adanos-scheduler/scheduler/task"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/scheduler"
)

type Provider struct{}

func (p Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		scheduler.Provider(p.jobs),
	}
}

func (p Provider) Register(app infra.Binder) {
	app.MustSingletonOverride(task.NewJobPlanLoaderTask)
}

func (p Provider) Boot(app infra.Resolver) {
}

func (p Provider) jobs(cc infra.Resolver, creator scheduler.JobCreator) {
	cc.MustResolve(func(jobPlanLoaderTask *task.JobPlanLoaderTask) error {
		// 创建计划任务加载任务
		cc.Must(jobPlanLoaderTask.Handle(cc))
		cc.Must(creator.Add("job-plan-loader", "@every 5s", jobPlanLoaderTask))

		return nil
	})
}
