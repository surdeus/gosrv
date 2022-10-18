package hashes

import (
	"time"
	"hash/fnv"
)

func TimestampHash(t time.Time) {
	h := fnv.New64a()
	h.Write
}

