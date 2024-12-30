package authenticator

import "context"

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (bool, error)
}
