package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/coll"
)

type nodeCredentialRepoImpl struct {
	nodeCreds []repo.NodeCredential
	idSeq     int
}

func NewNodeCredentialRepo() repo.NodeCredentialRepo {
	return &nodeCredentialRepoImpl{
		nodeCreds: make([]repo.NodeCredential, 0),
		idSeq:     0,
	}
}

func (rp *nodeCredentialRepoImpl) GetCredentialIDsForNode(ctx context.Context, nodeID string) (credentialIDs []string, err error) {
	err = coll.MustNew(rp.nodeCreds).Filter(func(nc repo.NodeCredential) bool {
		return nc.NodeID == nodeID
	}).Map(func(nc repo.NodeCredential) string {
		return nc.CredentialID
	}).All(&credentialIDs)
	return
}

func (rp *nodeCredentialRepoImpl) BindCredentialWithNode(ctx context.Context, nodeID string, credentialID string) (id string, err error) {
	for _, nc := range rp.nodeCreds {
		if nc.NodeID == nodeID && nc.CredentialID == credentialID {
			return nc.ID, nil
		}
	}

	rp.idSeq++
	nc := repo.NodeCredential{
		ID:           strconv.Itoa(rp.idSeq),
		NodeID:       nodeID,
		CredentialID: credentialID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rp.nodeCreds = append(rp.nodeCreds, nc)
	return nc.ID, nil
}

func (rp *nodeCredentialRepoImpl) UnbindCredentialWithNode(ctx context.Context, nodeID string, credentialID string) error {
	return coll.Filter(rp.nodeCreds, &rp.nodeCreds, func(nc repo.NodeCredential) bool {
		return nc.NodeID != nodeID || nc.CredentialID != credentialID
	})
}

func (rp *nodeCredentialRepoImpl) UnbindAllForNode(ctx context.Context, nodeID string) error {
	return coll.Filter(rp.nodeCreds, &rp.nodeCreds, func(nc repo.NodeCredential) bool {
		return nc.NodeID != nodeID
	})
}

func (rp *nodeCredentialRepoImpl) UnbindAllForCredential(ctx context.Context, credentialID string) error {
	return coll.Filter(rp.nodeCreds, &rp.nodeCreds, func(nc repo.NodeCredential) bool {
		return nc.CredentialID == credentialID
	})
}
