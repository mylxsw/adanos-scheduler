package controller

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
)

type CredentialController struct {
	cc infra.Resolver
}

func NewCredentialController(cc infra.Resolver) web.Controller {
	return &CredentialController{cc: cc}
}

func (ctl CredentialController) Register(router web.Router) {
	router.Group("/credentials/", func(router web.Router) {
		router.Get("/", ctl.All)

		router.Post("/", ctl.Create)
		router.Get("/{cred_id}/", ctl.Credential)
		router.Put("/{cred_id}/", ctl.Update)
		router.Delete("/{cred_id}/", ctl.Remove)
	})
}

// CredentialReq 凭据更新请求
// Password/PrivateKey/PrivateKeyPassphrase/Token 这几个字段，如果更新时设置为值 ********，则不更新该字段，保留原值
type CredentialReq struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Type                 string `json:"type"`
	User                 string `json:"user"`
	Password             string `json:"password"`
	PrivateKey           string `json:"private_key"`
	PrivateKeyPassphrase string `json:"private_key_passphrase"`
	Token                string `json:"token"`
}

func (credReq CredentialReq) Validate(request web.Request) error {
	// TODO 表单校验
	return nil
}

func (credReq CredentialReq) Transform() repo.Credential {
	return repo.Credential{
		Name:                 credReq.Name,
		Description:          credReq.Description,
		Type:                 repo.CredentialType(credReq.Type),
		User:                 credReq.User,
		Password:             credReq.Password,
		PrivateKey:           credReq.PrivateKey,
		PrivateKeyPassphrase: credReq.PrivateKeyPassphrase,
		Token:                credReq.Token,
	}
}

// Create 创建一个新的凭据
func (ctl CredentialController) Create(ctx context.Context, req web.Request, credSrv service.CredentialService) (IDResponse, error) {
	var credReq CredentialReq
	if err := req.Unmarshal(&credReq); err != nil {
		return IDResponse{}, err
	}

	req.Validate(credReq, true)

	credID, err := credSrv.CreateCredential(ctx, credReq.Transform())
	return IDResponse{ID: credID}, err
}

// All 返回所有的凭据列表（不包含敏感信息）
func (ctl CredentialController) All(ctx context.Context, credSrv service.CredentialService) ([]service.CredentialMasked, error) {
	return credSrv.AllMasked(ctx)
}

// Update 更新凭证信息
func (ctl CredentialController) Update(ctx context.Context, req web.Request, credSrv service.CredentialService) error {
	var credReq CredentialReq
	if err := req.Unmarshal(&credReq); err != nil {
		return err
	}

	req.Validate(credReq, true)

	credID := req.PathVar("cred_id")
	return credSrv.UpdateCredential(ctx, credID, credReq.Transform())
}

// Credential 查看单个凭据信息，密码秘钥等敏感字段返回 ******** 代替
func (ctl CredentialController) Credential(ctx context.Context, req web.Request, credSrv service.CredentialService) (*repo.Credential, error) {
	credID := req.PathVar("cred_id")
	cred, err := credSrv.Credential(ctx, credID)
	if err != nil {
		return nil, err
	}

	if cred.Password != "" {
		cred.Password = service.PasswordZeroValue
	}
	if cred.PrivateKeyPassphrase != "" {
		cred.PrivateKeyPassphrase = service.PasswordZeroValue
	}
	if cred.PrivateKey != "" {
		cred.PrivateKey = service.PasswordZeroValue
	}
	if cred.Token != "" {
		cred.Token = service.PasswordZeroValue
	}

	return cred, nil
}

// Remove 移除一个凭据
func (ctl CredentialController) Remove(ctx context.Context, req web.Request, credSrv service.CredentialService) error {
	credID := req.PathVar("cred_id")
	return credSrv.RemoveCredential(ctx, credID)
}
