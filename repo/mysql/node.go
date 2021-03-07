package mysql

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type nodeRepoImpl struct {
}

func NewNodeRepo() repo.NodeRepo {
	return &nodeRepoImpl{}
}

func (rp *nodeRepoImpl) All(ctx context.Context) ([]repo.Node, error) {
	panic("implement me")
}

func (rp *nodeRepoImpl) Add(ctx context.Context, node repo.Node) (nodeID string, err error) {
	return "", nil
}

func (rp *nodeRepoImpl) GetByID(ctx context.Context, id string) (node *repo.Node, err error) {
	return nil, repo.ErrNotFound
}

func (rp *nodeRepoImpl) SelectByLabels(ctx context.Context, labels repo.Labels) (nodes []repo.Node, err error) {
	return []repo.Node{}, nil
}

func (rp *nodeRepoImpl) Update(ctx context.Context, id string, node repo.Node) error {
	panic("implement me")
}

func (rp *nodeRepoImpl) Remove(ctx context.Context, id string) error {
	panic("implement me")
}