package connector

import (
	"context"

	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/nexpose"
	"github.com/nortonlifelock/ttl"
	"github.com/pkg/errors"
)

// Connect creates a Nexpose session that implements the connector interface and is used to interact with the Nexpose API
func Connect(ctx context.Context, logger log.Logger, auth string, payload string) (session *Connection, err error) {

	ctx = ctxtest(ctx)

	if logger != nil {

		// Build the host data and http client
		var host *domain.Host
		var c client
		if host, c, err = unmarshalAuthCreateClient(ctx, auth, logger); err == nil {

			// Build the settings object
			var settings *Payload
			if settings, err = unmarshalPayload(payload); err == nil {

				// Create the API connection
				var api *nexpose.Session
				if api, err = createSession(ctx, c, host, logger); err == nil {

					// Create the connection object that implements the interface
					session = &Connection{
						ctx,
						api,
						settings,
						host,
						logger,
						ttl.Cache{},
						ttl.Cache{},
					}
				}
			}
		}
	} else {
		err = errors.New("nil logger passed to nexpose connection")
	}

	return session, err
}

// ctxtest processes the passed context and returns a context. If the ctx
// is nil then it returns a background context
func ctxtest(ctx context.Context) context.Context {

	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}
