package dealer

import (
	"testing"
)

func TestNewDealer(t *testing.T) {
	d := NewDealer()
	if d == nil {
		t.Fatal("NewDealer returned nil")
	}
	if d.Hand != nil {
		t.Error("NewDealer should have nil Hand")
	}
}
