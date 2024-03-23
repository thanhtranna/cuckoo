package cuckoo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket_Reset(t *testing.T) {
	var bkt bucket
	for i := byte(0); i < bucketSize; i++ {
		bkt[i] = fingerprint(i)
	}
	bkt.reset()
	for _, val := range bkt {
		assert.EqualValues(t, 0, val)
	}
}
