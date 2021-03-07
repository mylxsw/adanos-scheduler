package service_test

import (
	"context"
	"testing"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/adanos-scheduler/repo/mock"
	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/container"
	"github.com/stretchr/testify/assert"
)

func createMockJobService(ctx context.Context) service.JobService {
	cc := container.NewWithContext(ctx)
	mock.Provider{}.Register(cc)

	return service.NewJobService(cc)
}

func TestJobService(t *testing.T) {
	srv := createMockJobService(context.TODO())
	// 创建 jobs
	jobID1, err := srv.CreateJobWithPlan(context.TODO(), service.JobWithPlan{
		Job: repo.Job{
			Name:        "example",
			Command:     "ping -c 4 baidu.com",
			Description: "Ping baidu.com is a example",
		},
		Plans: nil,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, jobID1)

	// 创建带有执行计划的 jobs
	jobID2, err := srv.CreateJobWithPlan(context.TODO(), service.JobWithPlan{
		Job: repo.Job{Name: "whoami", Command: "whoami", Description: "who am i?"},
		Plans: []repo.JobPlan{
			{
				Name: "whoami-five-minutes",
				Plan: "@every 5m",
				NodeSelectors: repo.Labels{
					"type": "infra",
					"core": "6",
				},
				ExecuteMode: repo.JobExecuteModeAll,
			},
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, jobID2)

	// 查询所有 jobs
	jobs, err := srv.AllJobs(context.TODO())
	assert.NoError(t, err)
	assert.True(t, len(jobs) == 2)

	// 查询指定的 job
	job1WithPlans, err := srv.JobWithPlans(context.TODO(), jobID1)
	assert.NoError(t, err)
	assert.NotEmpty(t, job1WithPlans.Job.Command)
	assert.Empty(t, job1WithPlans.Plans)

	job2WithPlans, err := srv.JobWithPlans(context.TODO(), jobID2)
	assert.NoError(t, err)
	assert.NotEmpty(t, job2WithPlans.Job.Command)
	assert.NotEmpty(t, job2WithPlans.Plans)
}
