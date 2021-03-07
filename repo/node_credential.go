package repo

import (
	"context"
	"time"
)

// NodeCredential 节点和凭证的关联关系
type NodeCredential struct {
	ID           string
	NodeID       string
	CredentialID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type NodeCredentialRepo interface {
	// GetCredentialIDsForNode 获取节点的凭证
	GetCredentialIDsForNode(ctx context.Context, nodeID string) (credentialIDs []string, err error)
	// BindCredentialWithNode 绑定节点和凭证的关系
	BindCredentialWithNode(ctx context.Context, nodeID string, credentialID string) (id string, err error)
	// UnbindCredentialWithNode 移除节点和凭证的关系
	UnbindCredentialWithNode(ctx context.Context, nodeID string, credentialID string) error
	// UnbindAllForNode 解除节点绑定的所有凭证
	UnbindAllForNode(ctx context.Context, nodeID string) error
	// UnbindAllForCredential 解除凭证绑定的所有节点
	UnbindAllForCredential(ctx context.Context, credentialID string) error
}
