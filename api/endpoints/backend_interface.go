package endpoints

import (
	"github.com/nortonlifelock/domain"
)

type endpoint interface {
	update(domain.User, domain.Permission, string) (*GeneralResp, int, error)
	create(domain.User, domain.Permission) (*GeneralResp, int, error)
	delete(domain.User, domain.Permission) (*GeneralResp, int, error)
	verify() bool
}
