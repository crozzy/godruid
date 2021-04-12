package godruid

import (
	"testing"
)

func TestInstaciateClient(t *testing.T) {
	_ = NewClient([]string{}, "")
}
