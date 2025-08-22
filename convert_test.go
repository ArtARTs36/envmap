package envmap

import (
	"encoding/base64"
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
		Mode           string                       `env:"MODE"`
		DB             dbCfg1                       `envPrefix:"DB_"`
		EmptyField     string                       `env:"EMPTY_FIELD"`
		RequiredField  int                          `env:"REQUIRED_FIELD,required"`
		UserMap        map[string]string            `env:"USER_MAP"`
		Marshalling    marshallingString            `env:"MARSHALING"`
		MarshallingMap map[string]marshallingString `env:"MARSHALING_MAP"`
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
				UserMap: map[string]string{
					"id-1": "ab",
				},
				Marshalling: "test",
				MarshallingMap: map[string]marshallingString{
					"k1": "v1",
				},
			},
			Expected: map[string]string{
				"APP_MODE":           "prod",
				"APP_DB_TIMEOUT":     "1s",
				"APP_REQUIRED_FIELD": "3",
				"APP_USER_MAP":       "id-1:ab",
				"APP_MARSHALING":     "dGVzdA==",
				"APP_MARSHALING_MAP": "k1:djE=",
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
		Input    *value
		Expected string
	}{
		{
			Title: "time.Duration: filled",
			Input: &value{
				value: time.Second,
			},
			Expected: "1s",
		},
		{
			Title: "time.Duration: zero",
			Input: &value{
				value: 0 * time.Second,
			},
			Expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			v, err := valueToString(c.Input)
			require.NoError(t, err)

			assert.Equal(t, c.Expected, v)
		})
	}
}

type marshallingString string

func (u marshallingString) MarshalText() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(u))
	return []byte(encoded), nil
}
