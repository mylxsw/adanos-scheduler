package repo

import (
	"context"
	"time"
)

// CredentialType 节点访问凭证类型
type CredentialType string

const (
	// CredentialTypePassword 用户名密码鉴权模式
	CredentialTypePassword CredentialType = "password"
	// CredentialTypeToken 基于 Token 的鉴权模式
	CredentialTypeToken CredentialType = "token"
	// CredentialTypePrivateKey 基于私钥的鉴权模式
	CredentialTypePrivateKey CredentialType = "private_key"
)

// Credential 节点访问凭证
type Credential struct {
	ID                   string
	Name                 string
	Description          string
	Type                 CredentialType
	User                 string
	Password             string
	PrivateKey           string
	PrivateKeyPassphrase string
	Token                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// CredentialRepo 凭证管理仓库
type CredentialRepo interface {
	// All 返回所有凭据
	All(ctx context.Context) ([]Credential, error)
	// Add 新增凭证
	Add(ctx context.Context, cred Credential) (credentialID string, err error)
	// GetByID 根据 ID 查询凭证
	GetByID(ctx context.Context, credentialID string) (credential *Credential, err error)
	// GetByIDs 根据 ID 批量查询凭证
	GetByIDs(ctx context.Context, credentialIDs ...string) (credentials []Credential, err error)
	// Update 更新凭证
	Update(ctx context.Context, credentialID string, cred Credential) error
	// Remove 删除凭证
	Remove(ctx context.Context, credentialID string) error
}
