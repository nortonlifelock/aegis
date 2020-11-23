package job

import "context"

type ctxWrapper struct {
	ctx    context.Context
	cancel context.CancelFunc
}
