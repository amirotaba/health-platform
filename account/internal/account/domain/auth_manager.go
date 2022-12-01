package domain

import (
	"context"
	"fmt"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type AuthError struct {
	Status  int
	Message string
}

func (e AuthError) Error() string {
	return fmt.Sprintf("Authorizing error: %d, %s", e.Status, e.Message)
}

type AuthManager interface {
	AuthenticationByToken(ctx context.Context, token string) (account entity.Identity, err error)
	AuthorizingByToken(ctx context.Context, account entity.Identity) (err error)
}
