package producer

import (
	"context"
	"testing"
)

func TestWriteMessage(t *testing.T) {

	parent := context.Background()
	defer parent.Done()
	err := Push(parent, "someData.Key", []byte("Asdf"))
	if err != nil {
		t.Error("error while pushing message into kafka:", err)
	}

}
