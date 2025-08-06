package envmap

import (
	"fmt"
	"reflect"
	"strings"
)

func MustConvert(v interface{}, opt ...Opt) map[string]string {
	res, err := Convert(v, opt...)
	if err != nil {
		panic(err)
	}

	return res
}

func Convert(v interface{}, opt ...Opt) (map[string]string, error) {
	emap := &envMap{
		values: make(map[string]string),
	}

	cfg := &config{}
	cfg.apply(opt...)

	err := convert(v, cfg.Prefix, emap)
	if err != nil {
		return nil, err
	}

	return emap.values, nil
}

func convert(v interface{}, prefix string, emap *envMap) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		envKey := field.Tag.Get("env")

		if envPrefix := field.Tag.Get("envPrefix"); envPrefix != "" {
			prefix += envPrefix
		}

		if fv := rv.Field(i); fv.Kind() == reflect.Struct {
			err := convert(fv.Interface(), prefix, emap)
			if err != nil {
				return fmt.Errorf("converting field %s: %v", field.Name, err)
			}
		} else {
			envKey = prefix + envKey

			if envKey != "" {
				emap.add(envKey, fmt.Sprintf("%v", valueToString(rv.Field(i).Interface())))
			}
		}
	}

	return nil
}

func valueToString(value interface{}) string {
	switch v := value.(type) {
	case []string:
		return strings.Join(v, ",")
	case []interface{}:
		vs := make([]string, len(v))
		for i := range v {
			vs[i] = valueToString(v[i])
		}

		return strings.Join(vs, ",")
	default:
		return fmt.Sprintf("%v", value)
	}
}

type envMap struct {
	values map[string]string
}

func (e envMap) add(key string, val string) {
	e.values[key] = val
}
