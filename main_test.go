package main

import (
	"testing"
)

func TestGetHello(t *testing.T) {
	if GetHello() != "Hello, world!" {
		t.Error("GetHello() did not return \"Hello, world!\"")
	}
}