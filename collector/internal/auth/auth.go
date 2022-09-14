package auth

import "context"

type AuthorizationCode int

var Authorized AuthorizationCode = 0

func Authorize(ctx context.Context, authToken string) (authCode AuthorizationCode) {
	return Authorized
}
