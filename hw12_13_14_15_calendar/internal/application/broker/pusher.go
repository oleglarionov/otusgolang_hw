package broker

import "context"

type Pusher interface {
	Push(ctx context.Context, data interface{}) error
}
