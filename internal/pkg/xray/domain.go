package xray

type XrayObservatory struct {
	Alive           bool   `json:"alive"`
	Delay           int64  `json:"delay"`
	OutboundTag     string `json:"outbound_tag"`
	LastSeenTime    int64  `json:"last_seen_time"`
	LastTryTime     int64  `json:"last_try_time"`
	LastErrorReason string `json:"last_error_reason,omitempty"`
}

type XrayDebugVars struct {
	Observatories map[string]XrayObservatory `json:"observatory"`
}
