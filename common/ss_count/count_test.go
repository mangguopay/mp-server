package ss_count

import (
	"testing"
)

func TestSub(t *testing.T) {
	t.Logf(Sub("100", "50").String())
}
