package main

import (
	"testing"

	"github.com/Enriquefft/SAMS/internal/helloworld"
)

func TestMainGreet(t *testing.T) {
	t.Parallel()
	got := helloworld.Greet("")
	if got != "Hello, World!" {
		t.Fatalf("got %q, want %q", got, "Hello, World!")
	}
}
