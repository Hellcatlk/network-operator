package switchs

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New(context.Background(), "test", "test://1234")
	if err != nil {
		t.Fatal(err)
	}
}
