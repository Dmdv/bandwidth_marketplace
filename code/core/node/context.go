package node

import (
	ctx "context"

	"github.com/0chain/bandwidth_marketplace/code/core/context"
)

const SelfNode context.CtxKey = "SELF_NODE"

// GetNodeContext setups a common with the self node.
func GetNodeContext() ctx.Context {
	return ctx.WithValue(ctx.Background(), SelfNode, self)
}
