package envmap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	type dbCfg1 struct {
		Timeout time.Duration `env:"TIMEOUT"`
	}

	type testCfg1 struct {
		Mode          string `env:"MODE"`
		DB            dbCfg1 `envPrefix:"DB_"`
		EmptyField    string `env:"EMPTY_FIELD"`
		RequiredField int    `env:"REQUIRED_FIELD,required"`
	}

	cases := []struct {
		Title    string
		Config   interface{}
		Expected map[string]string
		Opts     []Opt
	}{
		{
			Title:    "Empty config",
			Config:   struct{}{},
			Expected: map[string]string{},
		},
		{
			Title:    "Non-filled application config",
			Config:   testCfg1{},
			Expected: map[string]string{},
			Opts: []Opt{
				WithPrefix("APP_"),
			},
		},
		{
			Title: "Filled application config",
			Config: testCfg1{
				Mode: "prod",
				DB: dbCfg1{
					Timeout: time.Second,
				},
				RequiredField: 3,
			},
			Expected: map[string]string{
				"APP_MODE":           "prod",
				"APP_DB_TIMEOUT":     "1s",
				"APP_REQUIRED_FIELD": "3",
			},
			Opts: []Opt{
				WithPrefix("APP_"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			got, err := Convert(c.Config, c.Opts...)
			require.NoError(t, err)
			assert.Equal(t, c.Expected, got)
		})
	}
}

func TestValueToString(t *testing.T) {
	cases := []struct {
		Title    string
		Input    interface{}
		Expected string
	}{
		{
			Title:    "time.Duration: filled",
			Input:    time.Second,
			Expected: "1s",
		},
		{
			Title:    "time.Duration: zero",
			Input:    0 * time.Second,
			Expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			assert.Equal(t, c.Expected, valueToString(c.Input))
		})
	}
}
