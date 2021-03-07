package mysql

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
)

type nodeCredentialRepoImpl struct {
}

func NewNodeCredentialRepo() repo.NodeCredentialRepo {
	return &nodeCredentialRepoImpl{}
}

func (rp *nodeCredentialRepoImpl) GetCredentialIDsForNode(ctx context.Context, nodeID string) (credentialIDs []string, err error) {
	panic("implement me")
}

func (rp *nodeCredentialRepoImpl) BindCredentialWithNode(ctx context.Context, nodeID string, credentialID string) (id string, err error) {
	panic("implement me")
}

func (rp *nodeCredentialRepoImpl) UnbindCredentialWithNode(ctx context.Context, nodeID string, credentialID string) error {
	panic("implement me")
}

func (rp *nodeCredentialRepoImpl) UnbindAllForNode(ctx context.Context, nodeID string) error {
	panic("implement me")
}

func (rp *nodeCredentialRepoImpl) UnbindAllForCredential(ctx context.Context, credentialID string) error {
	panic("implement me")
}
