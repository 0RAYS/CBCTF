package utils

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	wg  *errgroup.Group
	ctx context.Context
}

func (g *Group) Go(f func() error) {
	g.wg.Go(func() error {
		select {
		case <-g.ctx.Done():
			return g.ctx.Err()
		default:
			return f()
		}
	})
}

func (g *Group) Wait() error {
	return g.wg.Wait()
}

func NewGroup(ctx context.Context) *Group {
	wg, ctx := errgroup.WithContext(ctx)
	return &Group{wg: wg, ctx: ctx}
}
