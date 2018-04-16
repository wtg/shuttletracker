package updater

import (
	"testing"
)

func TestITrakTimeDate(t *testing.T) {
	time := "time:52957"
	date := "date:04162018"
	parsed, err := itrakTimeDate(time, date)
	if err != nil {
		t.Fail(err)
	}
	t.Log()
}
