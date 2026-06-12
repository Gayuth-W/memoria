package tests

import (
	"testing"

	"memoria/internal/cache"
)

func TestRedis(t *testing.T) {

	r := cache.NewRedisCache()

	err := r.Set(
		"hello",
		"world",
	)

	if err != nil {
		t.Fatal(err)
	}

	value, err := r.Get(
		"hello",
	)

	if err != nil {
		t.Fatal(err)
	}

	if value != "world" {
		t.Fatalf(
			"expected world got %s",
			value,
		)
	}
}
