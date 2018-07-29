package main

import (
	"os"
	"testing"
)

func BenchmarkWalkWithNative(b *testing.B) {
	gp := os.Getenv("GOPATH")
	err := walkWithNative(gp)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkWalkWithLib(b *testing.B) {
	gp := os.Getenv("GOPATH")
	err := walkWithLib(gp)
	if err != nil {
		b.Error(err)
	}
}
