package service

import (
	"context"
	"fmt"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
)

// NodeWithCredentialID 节点信息，包含关联的凭证 id 列表
type NodeWithCredentialID struct {
	Node          repo.Node `json:"node"`
	CredentialIDs []string  `json:"credential_ids"`
}

// NodeService 操作节点的接口
type NodeService interface {
	// AllNodes 返回所有的 Node
	AllNodes(ctx context.Context) ([]repo.Node, error)
	// SelectNodes 选择匹配标签的节点列表
	SelectNodes(ctx context.Context, labelSelector repo.Labels) ([]repo.Node, error)
	// Node 查询单个节点信息
	Node(ctx context.Context, nodeID string) (*NodeWithCredentialID, error)
	// CreateNode 新增一个节点
	CreateNode(ctx context.Context, node repo.Node, credentialIDs ...string) (nodeID string, err error)
	// UpdateNode 更新节点
	UpdateNode(ctx context.Context, nodeID string, node repo.Node, credentialIDs ...string) error
	// RemoveNode 删除一个节点
	RemoveNode(ctx context.Context, nodeID string) error
}

// nodeService NodeService 接口实现
type nodeService struct {
	cc                 container.Container
	nodeRepo           repo.NodeRepo           `autowire:"@"`
	credentialRepo     repo.CredentialRepo     `autowire:"@"`
	nodeCredentialRepo repo.NodeCredentialRepo `autowire:"@"`
}

// NewNodeService 创建一个 NodeService 实例
func NewNodeService(cc container.Container) NodeService {
	srv := &nodeService{cc: cc}
	cc.MustAutoWire(srv)
	return srv
}

func (srv *nodeService) AllNodes(ctx context.Context) ([]repo.Node, error) {
	return srv.nodeRepo.All(ctx)
}

func (srv *nodeService) SelectNodes(ctx context.Context, labelSelector repo.Labels) ([]repo.Node, error) {
	return srv.nodeRepo.SelectByLabels(ctx, labelSelector)
}

func (srv *nodeService) Node(ctx context.Context, nodeID string) (*NodeWithCredentialID, error) {
	node, err := srv.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	credIds, err := srv.nodeCredentialRepo.GetCredentialIDsForNode(ctx, node.ID)
	if err != nil {
		return &NodeWithCredentialID{Node: *node}, err
	}

	return &NodeWithCredentialID{
		Node:          *node,
		CredentialIDs: credIds,
	}, nil
}

func (srv *nodeService) CreateNode(ctx context.Context, node repo.Node, credentialIDs ...string) (nodeID string, err error) {
	nodeID, err = srv.nodeRepo.Add(ctx, node)
	if err != nil {
		return "", err
	}

	for _, credID := range credentialIDs {
		if _, err := srv.nodeCredentialRepo.BindCredentialWithNode(ctx, nodeID, credID); err != nil {
			log.WithFields(log.Fields{
				"node_id":       nodeID,
				"credential_id": credID,
			}).Errorf("bind node with credential failed: %v", err)
		}
	}

	return nodeID, nil
}

func (srv *nodeService) UpdateNode(ctx context.Context, nodeID string, node repo.Node, credentialIDs ...string) error {
	if err := srv.nodeRepo.Update(ctx, nodeID, node); err != nil {
		return err
	}

	if err := srv.nodeCredentialRepo.UnbindAllForNode(ctx, nodeID); err != nil {
		return fmt.Errorf("unbind all credentials for node failed: %w", err)
	}

	for _, credID := range credentialIDs {
		if _, err := srv.nodeCredentialRepo.BindCredentialWithNode(ctx, nodeID, credID); err != nil {
			log.WithFields(log.Fields{
				"node_id":       nodeID,
				"credential_id": credID,
			}).Errorf("bind node with credential failed: %v", err)
		}
	}

	return nil
}

func (srv *nodeService) RemoveNode(ctx context.Context, nodeID string) error {
	if err := srv.nodeCredentialRepo.UnbindAllForNode(ctx, nodeID); err != nil {
		return fmt.Errorf("unbind all credentials for node failed: %w", err)
	}

	return srv.nodeRepo.Remove(ctx, nodeID)
}
