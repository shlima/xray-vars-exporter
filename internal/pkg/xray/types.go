package xray

import "context"

type IClient interface {
	GetDebugVars(ctx context.Context) (*XrayDebugVars, error)
}
