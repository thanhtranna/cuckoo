package cuckoo

import (
	"math/bits"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexAndFP(t *testing.T) {
	data := []byte("seif")
	bucketPow := uint(bits.TrailingZeros(1024))
	i1, fp := getIndexAndFingerprint(data, bucketPow)
	i2 := getAltIndex(fp, i1, bucketPow)
	i11 := getAltIndex(fp, i2, bucketPow)
	i22 := getAltIndex(fp, i11, bucketPow)
	assert.EqualValues(t, i11, i1)
	assert.EqualValues(t, i22, i2)
}

func TestCap(t *testing.T) {
	const capacity = 10000
	res := getNextPow2(uint64(capacity)) / bucketSize
	assert.EqualValues(t, res, 4096)
}
