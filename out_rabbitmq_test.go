package main

import (
	"context"
	"testing"
	"unsafe"
)

func TestFLBPluginRegister(t *testing.T) {
	bk := context.Background()
	var ctx = unsafe.Pointer(&bk)

	result := FLBPluginRegister(ctx)
	if result != 0 {
		t.Error("Failed to register plugin, expected", 0, "but found", result)
	}
}
