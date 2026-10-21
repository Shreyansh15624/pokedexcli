package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://www.example.com/",
			val: []byte("testdata"),
		},
		{
			key: "https://www.example.com/",
			val: []byte("moretestdata"),
		},
		{
			key: "",
			val: []byte(nil),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test Case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("the cache value doesn't exists")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected value doesn't match the input data")
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
