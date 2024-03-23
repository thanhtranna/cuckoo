package cuckoo

import (
	"bufio"
	"crypto/rand"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertion(t *testing.T) {
	cf := NewFilter(1000000)
	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fd)

	var values [][]byte
	var lineCount uint
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if cf.InsertUnique(s) {
			lineCount++
		}
		values = append(values, s)
	}

	count := cf.Count()
	if count != lineCount {
		t.Errorf("Expected count = %d, instead count = %d", lineCount, count)
	}

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.Count()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}

func TestEncodeDecode(t *testing.T) {
	cf := NewFilter(8)
	cf.buckets = []bucket{
		[4]fingerprint{1, 2, 3, 4},
		[4]fingerprint{5, 6, 7, 8},
	}
	cf.count = 8
	bytes := cf.Encode()
	ncf, err := Decode(bytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(cf, ncf) {
		t.Errorf("Expected %v, got %v", cf, ncf)
	}
}

func TestDecode(t *testing.T) {
	ncf, err := Decode([]byte(""))
	if err == nil {
		t.Errorf("Expected err, got nil")
	}
	if ncf != nil {
		t.Errorf("Expected nil, got %v", ncf)
	}
}

func TestInsert(t *testing.T) {
	const cap = 10000
	filter := NewFilter(cap)

	var hash [32]byte
	io.ReadFull(rand.Reader, hash[:])

	for i := 0; i < 100; i++ {
		filter.Insert(hash[:])
	}

	assert.EqualValues(t, filter.Count(), 8)
}

func TestFilter_Lookup(t *testing.T) {
	const cap = 10000

	var (
		m      = make(map[[32]byte]struct{})
		filter = NewFilter(cap)
		hash   [32]byte
	)

	for i := 0; i < cap; i++ {
		io.ReadFull(rand.Reader, hash[:])
		m[hash] = struct{}{}
		filter.Insert(hash[:])
	}

	assert.EqualValues(t, len(m), 10000)

	var lookFail int
	for k := range m {
		if !filter.Lookup(k[:]) {
			lookFail++
		}
	}

	assert.EqualValues(t, lookFail, 0)
}

func TestFilter_Reset(t *testing.T) {
	const cap = 10000

	var (
		filter        = NewFilter(cap)
		hash          [32]byte
		insertSuccess int
		insertFails   int
	)

	for i := 0; i < 10*cap; i++ {
		io.ReadFull(rand.Reader, hash[:])

		if filter.Insert(hash[:]) {
			insertSuccess++
		} else {
			insertFails++
			filter.Reset()
		}
	}

	assert.EqualValues(t, insertSuccess, 99994)
	assert.EqualValues(t, insertFails, 6)
}

func BenchmarkFilter_Reset(b *testing.B) {
	const cap = 10000
	filter := NewFilter(cap)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Reset()
	}
}

func BenchmarkFilter_Insert(b *testing.B) {
	const cap = 10000
	filter := NewFilter(cap)

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}
}

func BenchmarkFilter_Lookup(b *testing.B) {
	const cap = 10000
	filter := NewFilter(cap)

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
