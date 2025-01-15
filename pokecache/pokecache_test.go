package pokecache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cases :=
		[]struct {
			input    time.Duration
			expected reflect.Type
		}{
			{
				input:    time.Duration(10) * time.Second,
				expected: reflect.TypeOf(Cache{}),
			},
		}

	for _, c := range cases {
		actual := NewCache(c.input)

		if reflect.TypeOf(actual) != c.expected {
			t.Errorf("Expected type: %v, Get %v", c.expected, actual)
		}
	}
}

func TestAddGetReap(t *testing.T) {
	cases := []struct {
		key       string
		val       []byte
		delay     time.Duration
		interval  time.Duration
		available bool
	}{
		{
			key:       "https://example.com",
			val:       []byte("testdata"),
			delay:     time.Duration(0),
			interval:  time.Duration(500) * time.Millisecond,
			available: true,
		},
		{
			key:       "https://example.com/path",
			val:       []byte("testdatamore"),
			delay:     time.Duration(0),
			interval:  time.Duration(500) * time.Millisecond,
			available: true,
		},
		{
			key:       "https://example.com/sleep2",
			val:       []byte("testdata"),
			delay:     time.Duration(2) * time.Millisecond,
			interval:  time.Duration(500) * time.Millisecond,
			available: true,
		},
		{
			key:       "https://example.com/sleep7",
			val:       []byte("testdatamore"),
			delay:     time.Duration(700) * time.Millisecond,
			interval:  time.Duration(500) * time.Millisecond,
			available: false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(c.interval)
			cache.Add(c.key, c.val)
			time.Sleep(c.delay)
			val, ok := cache.Get(c.key)
			if ok != c.available {
				if c.available {
					t.Errorf("Expected to find key: %s", c.key)
				} else {
					t.Errorf("Expected to not find key: %s", c.key)
				}
				return
			}
			if c.available && string(val) != string(c.val) {
				t.Errorf("Expected to find val: %s", string(c.val))
				return
			}
		})
		t.Logf("Successfully AddGet %s, %s", c.key, string(c.val))
	}
}
