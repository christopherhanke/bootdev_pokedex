package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moredata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}

}

func TestReapLoop(t *testing.T) {
	const baseTime = time.Millisecond * 5
	const waitTime = time.Millisecond*5 + baseTime
	cache := NewCache(baseTime)

	testcase := struct {
		key string
		val []byte
	}{
		key: "https://example.com",
		val: []byte("testdata"),
	}

	cache.Add(testcase.key, testcase.val)

	_, ok := cache.Get(testcase.key)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(testcase.key)
	if ok {
		t.Errorf("did not expect to find key")
		return
	}

}
