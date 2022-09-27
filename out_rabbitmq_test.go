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

func TestLogInfo(t *testing.T) {
	var testmessage = "hoge"
	LogInfo(testmessage)
}

func TestLogError(t *testing.T) {
	var testmessage = "hoge"
	var err error
	LogError(testmessage, err)
}
