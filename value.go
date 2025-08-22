package envmap

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structtag"
)

type value struct {
	value interface{}
	tags  *structtag.Tags
}

func (v *value) sliceSeparator() string {
	sep, err := v.tags.Get("envSeparator")
	if err != nil {
		return ","
	}
	return sep.Name
}

func valueToString(val *value) (string, error) {
	switch v := val.value.(type) {
	case encoding.TextMarshaler:
		text, err := v.MarshalText()
		if err != nil {
			return "", err
		}
		return string(text), nil
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		if v == 0 {
			return "", nil
		}
		return fmt.Sprintf("%v", v), nil
	case time.Duration:
		if v == 0 {
			return "", nil
		}
		return v.String(), nil
	case []string:
		return strings.Join(v, val.sliceSeparator()), nil
	case []interface{}:
		vs := make([]string, len(v))
		for i := range v {
			var err error
			vs[i], err = valueToString(&value{
				value: v[i],
			})
			if err != nil {
				return "", err
			}
		}
		return strings.Join(vs, val.sliceSeparator()), nil
	default:
		if reflect.ValueOf(val.value).Kind() == reflect.Map {
			mv, err := resolveMapValue(val)
			if err != nil {
				return "", fmt.Errorf("resolve map value: %w", err)
			}
			return mv, nil
		}
		return fmt.Sprintf("%v", val.value), nil
	}
}

func resolveMapValue(val *value) (string, error) {
	rv := reflect.ValueOf(val.value)
	iter := rv.MapRange()

	vs := make([]string, 0)
	for iter.Next() {
		mk, err := valueToString(&value{
			value: iter.Key(),
		})
		if err != nil {
			return "", err
		}
		mv, err := valueToString(&value{
			value: iter.Value().Interface(),
		})
		if err != nil {
			return "", err
		}

		vs = append(vs, fmt.Sprintf("%s:%s", mk, mv))
	}

	return strings.Join(vs, ","), nil
}
