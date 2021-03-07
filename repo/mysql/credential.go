package mysql

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type credentialRepoImpl struct {
}

func NewCredentialRepo() repo.CredentialRepo {
	return &credentialRepoImpl{}
}

func (rp *credentialRepoImpl) All(ctx context.Context) ([]repo.Credential, error) {
	panic("implement me")
}

func (rp *credentialRepoImpl) Add(ctx context.Context, cred repo.Credential) (credentialID string, err error) {
	panic("implement me")
}

func (rp *credentialRepoImpl) GetByID(ctx context.Context, credentialID string) (credential *repo.Credential, err error) {
	panic("implement me")
}

func (rp *credentialRepoImpl) GetByIDs(ctx context.Context, credentialIDs ...string) (credentials []repo.Credential, err error) {
	panic("implement me")
}

func (rp *credentialRepoImpl) Update(ctx context.Context, credentialID string, cred repo.Credential) error {
	panic("implement me")
}

func (rp *credentialRepoImpl) Remove(ctx context.Context, credentialID string) error {
	panic("implement me")
}
