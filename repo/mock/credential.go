package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/go-utils/str"
)

type credentialRepoImpl struct {
	credentials []repo.Credential
	idSeq       int
}

func NewCredentialRepo() repo.CredentialRepo {
	return &credentialRepoImpl{
		credentials: make([]repo.Credential, 0),
		idSeq:       0,
	}
}

func (rp *credentialRepoImpl) All(ctx context.Context) ([]repo.Credential, error) {
	return rp.credentials, nil
}

func (rp *credentialRepoImpl) Add(ctx context.Context, cred repo.Credential) (credentialID string, err error) {
	rp.idSeq++
	cred.ID = strconv.Itoa(rp.idSeq)
	cred.CreatedAt = time.Now()
	cred.UpdatedAt = time.Now()

	rp.credentials = append(rp.credentials, cred)
	return cred.ID, nil
}

func (rp *credentialRepoImpl) GetByID(ctx context.Context, credentialID string) (credential *repo.Credential, err error) {
	for _, cred := range rp.credentials {
		if cred.ID == credentialID {
			return &cred, nil
		}
	}

	return nil, repo.ErrNotFound
}

func (rp *credentialRepoImpl) GetByIDs(ctx context.Context, credentialIDs ...string) (credentials []repo.Credential, err error) {
	err = coll.Filter(rp.credentials, &credentials, func(cred repo.Credential) bool { return str.In(cred.ID, credentialIDs) })
	return
}

func (rp *credentialRepoImpl) Update(ctx context.Context, credentialID string, cred repo.Credential) error {
	for i, c := range rp.credentials {
		if c.ID == credentialID {
			cred.ID = credentialID
			cred.CreatedAt = c.CreatedAt
			cred.UpdatedAt = time.Now()
			rp.credentials[i] = cred
			return nil
		}
	}

	return repo.ErrNotFound
}

func (rp *credentialRepoImpl) Remove(ctx context.Context, credentialID string) error {
	return coll.Filter(rp.credentials, &rp.credentials, func(cred repo.Credential) bool { return cred.ID != credentialID })
}
