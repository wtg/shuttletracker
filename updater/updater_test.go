package updater

import (
	"testing"
	"time"
)

func TestITrakTimeDate(t *testing.T) {
	parsed, err := itrakTimeDate("time:52957", "date:04162018")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	expected := time.Date(2018, time.April, 16, 5, 29, 57, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Errorf("got %+v, expected %+v", parsed, expected)
	}

	parsed, err = itrakTimeDate("time:200546", "date:04162018")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	expected = time.Date(2018, time.April, 16, 20, 5, 46, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Errorf("got %+v, expected %+v", parsed, expected)
	}

	parsed, err = itrakTimeDate("time:2310", "date:04222018")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	expected = time.Date(2018, time.April, 22, 0, 23, 10, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Errorf("got %+v, expected %+v", parsed, expected)
	}
}
