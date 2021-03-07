package service

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
)

// PasswordZeroValue  密码类型字段默认值，使用该值则认为不修改
const PasswordZeroValue = "********"

// CredentialService 凭据管理服务
type CredentialService interface {
	// AllMasked 所有的凭据，敏感信息脱敏返回
	AllMasked(ctx context.Context) ([]CredentialMasked, error)
	// Credential 单个凭据信息
	Credential(ctx context.Context, credID string) (*repo.Credential, error)
	// Credentials 批量查询凭据
	Credentials(ctx context.Context, credIDs ...string) ([]repo.Credential, error)
	// Credentials 批量查询凭据，敏感信息脱敏返回
	CredentialsMasked(ctx context.Context, credIDs ...string) ([]CredentialMasked, error)
	// CreateCredential 创建一个新的凭据
	CreateCredential(ctx context.Context, cred repo.Credential) (credID string, err error)
	// UpdateCredential 更新一个凭据
	UpdateCredential(ctx context.Context, credID string, cred repo.Credential) error
	// RemoveCredential 删除一个凭据
	RemoveCredential(ctx context.Context, credID string) error
}

type credentialService struct {
	cc                 container.Container
	nodeCredentialRepo repo.NodeCredentialRepo `autowire:"@"`
	credentialRepo     repo.CredentialRepo     `autowire:"@"`
}

// NewCredentialService 创建一个凭据管理服务实现
func NewCredentialService(cc container.Container) CredentialService {
	srv := &credentialService{cc: cc}
	cc.MustAutoWire(srv)

	return srv
}

type CredentialMasked struct {
	ID          string
	Name        string
	Description string
	Type        repo.CredentialType
	User        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// maskCredential 将凭据信息转换为脱敏后的对象
func maskCredential(cred repo.Credential) CredentialMasked {
	return CredentialMasked{
		ID:          cred.ID,
		Name:        cred.Name,
		Description: cred.Description,
		Type:        cred.Type,
		User:        cred.User,
		CreatedAt:   cred.CreatedAt,
		UpdatedAt:   cred.UpdatedAt,
	}
}

func (srv credentialService) AllMasked(ctx context.Context) ([]CredentialMasked, error) {
	credentials, err := srv.credentialRepo.All(ctx)
	if err != nil {
		return nil, err
	}

	var credentialsMasked []CredentialMasked
	err = coll.Map(credentials, &credentialsMasked, maskCredential)
	return credentialsMasked, err
}

func (srv credentialService) Credential(ctx context.Context, credID string) (*repo.Credential, error) {
	return srv.credentialRepo.GetByID(ctx, credID)
}

func (srv credentialService) Credentials(ctx context.Context, credIDs ...string) ([]repo.Credential, error) {
	return srv.credentialRepo.GetByIDs(ctx, credIDs...)
}

func (srv credentialService) CredentialsMasked(ctx context.Context, credIDs ...string) ([]CredentialMasked, error) {
	creds, err := srv.credentialRepo.GetByIDs(ctx, credIDs...)
	if err != nil {
		return nil, err
	}

	var credsMasked []CredentialMasked
	err = coll.Map(creds, &credsMasked, maskCredential)

	return credsMasked, err
}

func (srv credentialService) CreateCredential(ctx context.Context, cred repo.Credential) (credID string, err error) {
	return srv.credentialRepo.Add(ctx, cred)
}

func (srv credentialService) UpdateCredential(ctx context.Context, credID string, cred repo.Credential) error {
	credOld, err := srv.credentialRepo.GetByID(ctx, credID)
	if err != nil {
		return err
	}

	if cred.Password == PasswordZeroValue {
		cred.Password = credOld.Password
	}
	if cred.PrivateKeyPassphrase == PasswordZeroValue {
		cred.PrivateKeyPassphrase = credOld.PrivateKeyPassphrase
	}
	if cred.PrivateKey == PasswordZeroValue {
		cred.PrivateKey = credOld.PrivateKey
	}
	if cred.Token == PasswordZeroValue {
		cred.Token = credOld.Token
	}

	return srv.credentialRepo.Update(ctx, credID, cred)
}

func (srv credentialService) RemoveCredential(ctx context.Context, credID string) error {
	return srv.credentialRepo.Remove(ctx, credID)
}
