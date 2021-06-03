package node

import (
	ctx "context"

	"github.com/MurashovVen/bandwidth-marketplace/code/core/context"
)

const SelfNode context.CtxKey = "SELF_NODE"

// GetNodeContext setups a common with the self node.
func GetNodeContext() ctx.Context {
	return ctx.WithValue(ctx.Background(), SelfNode, self)
}
