package cache

import "time"

const (
	defaultScrollTTL        = time.Minute
	defaultGraphTTL         = time.Minute
	defaultFileChunkSize    = 1024 // 1 KB
	defaultGraphIDKeyPrefix = "graph-"
	expirationEvent         = "__keyevent@0__:expired"
)
