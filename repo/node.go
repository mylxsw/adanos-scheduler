package repo

import (
	"context"
	"time"
)

// NodeType 节点类型
type NodeType string

const (
	// NodeTypeSSH ssh 连接型节点，目标服务器不需要部署 Agent，使用 ssh 远程执行
	NodeTypeSSH NodeType = "ssh"
	// NodeTypeAgent agent 型节点，需要在目标服务器上部署 Agent 节点
	NodeTypeAgent NodeType = "agent"
)

// NodeStatus 节点状态
type NodeStatus string

const (
	// NodeStatusEnabled 节点状态为启用
	NodeStatusEnabled NodeStatus = "enabled"
	// NodeStatusDisabled 节点状态为禁用
	NodeStatusDisabled NodeStatus = "disabled"
)

// Node 服务器节点
type Node struct {
	ID        string
	Name      string
	Labels    Labels
	IP        string
	Port      int
	Type      NodeType
	Status    NodeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NodeRepo 服务器节点操作接口
type NodeRepo interface {
	// All 返回所有节点
	All(ctx context.Context) ([]Node, error)
	// Add 新增节点
	Add(ctx context.Context, node Node) (nodeID string, err error)
	// GetByID 通过节点 ID 查询节点
	GetByID(ctx context.Context, id string) (node *Node, err error)
	// SelectByLabels 根据 labels 选择匹配的节点
	SelectByLabels(ctx context.Context, labels Labels) (nodes []Node, err error)
	// Update 更新节点
	Update(ctx context.Context, id string, node Node) error
	// Remove 删除节点
	Remove(ctx context.Context, id string) error
}
