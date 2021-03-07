package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/coll"
)

type nodeRepoImpl struct {
	nodes []repo.Node
	idSeq int
}

func NewNodeRepo() repo.NodeRepo {
	return &nodeRepoImpl{
		nodes: make([]repo.Node, 0),
		idSeq: 0,
	}
}

func (rp *nodeRepoImpl) All(ctx context.Context) ([]repo.Node, error) {
	return rp.nodes, nil
}

func (rp *nodeRepoImpl) Add(ctx context.Context, node repo.Node) (nodeID string, err error) {
	rp.idSeq++
	node.ID = strconv.Itoa(rp.idSeq)
	node.CreatedAt = time.Now()
	node.UpdatedAt = time.Now()

	rp.nodes = append(rp.nodes, node)
	return node.ID, nil
}

func (rp *nodeRepoImpl) GetByID(ctx context.Context, id string) (node *repo.Node, err error) {
	for _, node := range rp.nodes {
		if node.ID == id {
			return &node, err
		}
	}

	return nil, repo.ErrNotFound
}

func (rp *nodeRepoImpl) SelectByLabels(ctx context.Context, labels repo.Labels) (nodes []repo.Node, err error) {
	err = coll.Filter(rp.nodes, &nodes, func(node repo.Node) bool {
		for k, v := range labels {
			if nodeV, ok := node.Labels[k]; ok && nodeV == v {
				continue
			}

			return false
		}
		return true
	})
	return
}

func (rp *nodeRepoImpl) Update(ctx context.Context, id string, node repo.Node) error {
	for i, n := range rp.nodes {
		if n.ID == id {
			node.CreatedAt = n.CreatedAt
			node.UpdatedAt = time.Now()
			node.ID = id
			rp.nodes[i] = node
			return nil
		}
	}

	return repo.ErrNotFound
}

func (rp *nodeRepoImpl) Remove(ctx context.Context, id string) error {
	return coll.Filter(rp.nodes, &rp.nodes, func(node repo.Node) bool { return node.ID != id })
}
