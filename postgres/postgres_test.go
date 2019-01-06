package postgres

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigDatabaseURL(t *testing.T) {
	// this is what we want URL to be set to
	testVal := "this is a test haha"

	// first make sure it doesn't work
	v := viper.New()
	cfg, err := NewConfig(v)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	res := v.GetString("postgres.url")
	if res == testVal {
		t.Errorf("URL is %s", cfg.URL)
	}

	// now see if it works
	os.Setenv("DATABASE_URL", testVal)
	v = viper.New()
	cfg, err = NewConfig(v)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	res = v.GetString("postgres.url")
	if res != testVal {
		t.Errorf("URL is %s; expected %s", cfg.URL, testVal)
	}
}
