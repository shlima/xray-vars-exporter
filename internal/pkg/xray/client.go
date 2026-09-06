package xray

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

type Client struct {
	host string
	r    *resty.Client
}

func NewClient(host string) *Client {
	return &Client{host: host, r: resty.New()}
}

func (c *Client) GetDebugVars(ctx context.Context) (*XrayDebugVars, error) {
	out := new(XrayDebugVars)

	res, err := c.r.R().
		SetResult(out).
		SetContext(ctx).
		SetContentType("application/json").
		Get(c.host + "/debug/vars")

	switch {
	case err != nil:
		return nil, err
	case res.IsStatusFailure():
		return nil, fmt.Errorf("HTTP %s: %s", res.Status(), res.String())
	default:
		return out, nil
	}
}
