package main

import (
	"testing"
)

func TestServer(t *testing.T) {
	if 2+2 != 4 {
		t.Error("Test failed")
	}

}
