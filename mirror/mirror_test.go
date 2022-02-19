package mirror

import (
	"testing"
)

func TestMatchIds(t *testing.T) {
	ordId := "41de063f-0e7d-47dc-602c-fef9f7ab5ce5"
	clOrdId := clOrdId(ordId)

	if !matchIDs(ordId, clOrdId) {
		t.Error("matchIDs1 failed")
	}
}
