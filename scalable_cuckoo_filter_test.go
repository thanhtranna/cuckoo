package cuckoo

import (
	"crypto/rand"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalUse(t *testing.T) {
	filter := NewScalableCuckooFilter()
	for i := 0; i < 100000; i++ {
		filter.Insert([]byte("NewScalableCuckooFilter_" + strconv.Itoa(i)))
	}
	testStr := []byte("NewScalableCuckooFilter")
	b := filter.Insert(testStr)
	assert.True(t, b)
	b = filter.Lookup(testStr)
	assert.True(t, b)
	b = filter.Delete(testStr)
	assert.True(t, b)
	b = filter.Lookup(testStr)
	assert.False(t, b)
	b = filter.Lookup([]byte("NewScalableCuckooFilter_233"))
	assert.True(t, b)
	b = filter.InsertUnique([]byte("NewScalableCuckooFilter_599"))
	assert.False(t, b)
}

func TestScalableCuckooFilter_DecodeEncode(t *testing.T) {
	filter := NewScalableCuckooFilter(func(filter *ScalableCuckooFilter) {
		filter.loadFactor = 0.8
	})
	for i := 0; i < 100000; i++ {
		filter.Insert([]byte("NewScalableCuckooFilter_" + strconv.Itoa(i)))
	}
	bytes := filter.Encode()
	decodeFilter, err := decodeScalableFilter(bytes)
	assert.Nil(t, err)
	assert.Equal(t, decodeFilter.loadFactor, float32(0.8))
	b := decodeFilter.Lookup([]byte("NewScalableCuckooFilter_233"))
	assert.True(t, b)
	for i, f := range decodeFilter.filters {
		assert.Equal(t, f.count, filter.filters[i].count)
	}
}

func BenchmarkScalableCuckooFilter_Insert(b *testing.B) {
	filter := NewScalableCuckooFilter(func(filter *ScalableCuckooFilter) {
		filter.loadFactor = 0.8
	})

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}
}

func BenchmarkScalableCuckooFilter_Reset(b *testing.B) {
	filter := NewScalableCuckooFilter(func(filter *ScalableCuckooFilter) {
		filter.loadFactor = 0.8
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Reset()
	}
}

func BenchmarkScalableCuckooFilter_Lookup(b *testing.B) {
	filter := NewScalableCuckooFilter(func(filter *ScalableCuckooFilter) {
		filter.loadFactor = 0.8
	})

	var hash [32]byte
	for i := 0; i < 10000; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Lookup(hash[:])
	}
}
