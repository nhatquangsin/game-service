package database

import (
	"context"
	"fmt"

	"github.com/nhatquangsin/game-service/infra/repo/entc"
	"github.com/nhatquangsin/game-service/infra/utils/endpoint"
)

// EndpointTx returns a Middleware that wraps the `next` Endpoint in an
// Ent transaction.
func EndpointTx(client Client) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			tx, err := client.Master(ctx).Tx(ctx)
			if err != nil {
				return nil, err
			}

			defer func() {
				if v := recover(); v != nil {
					_ = tx.Rollback()
					panic(v)
				}
			}()

			ctx = entc.NewTxContext(ctx, tx)
			response, err := next(ctx, request)
			if err != nil {
				if rerr := tx.Rollback(); rerr != nil {
					err = fmt.Errorf("rolling back transaction: %v", rerr)
				}
				return nil, err
			}

			if err := tx.Commit(); err != nil {
				return nil, fmt.Errorf("rolling back transaction: %v", err)
			}

			return response, nil
		}
	}
}
