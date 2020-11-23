package domain

import (
	"context"
	"github.com/nortonlifelock/aegis/pkg/log"
)

// Job specifies the interface required by job implementations in order to execute properly through the dispatcher
type Job interface {
	Process(ctx context.Context, id string, appconfig Config, db DatabaseConnection, lstream log.Logger, payload string, jobConfig JobConfig, inSource []SourceConfig, outSource []SourceConfig) (err error)
}
