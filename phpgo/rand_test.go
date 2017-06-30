package phpgo

import (
	"testing"
)

func TestRandFloat64(t *testing.T) {
	t.Log("RandFloat64", RandFloat64())
}

func TestRandFloat(t *testing.T) {
	t.Log("RandFloat", RandFloat())
}

func TestRandInt(t *testing.T) {
	t.Logf("RandInt: %d", RandInt(99999))
}
