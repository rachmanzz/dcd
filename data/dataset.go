package data

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var varRe = regexp.MustCompile(`\{\{[^}]+\}\}`)

type DataSet struct {
	data map[string]any
}

func NewDataSet(src any) *DataSet {
	ds := &DataSet{data: make(map[string]any)}
	switch v := src.(type) {
	case map[string]any:
		ds.data = v
	case map[string]string:
		for k, v := range v {
			ds.data[k] = v
		}
	default:
		ds.data = toMap(reflect.ValueOf(src))
	}
	return ds
}

func (ds *DataSet) Set(key string, val any) {
	ds.data[key] = val
}

func (ds *DataSet) Get(path string) (any, bool) {
	parts := strings.Split(path, ".")
	current, ok := ds.data[parts[0]]
	if !ok {
		return nil, false
	}
	for i := 1; i < len(parts); i++ {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = m[parts[i]]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func (ds *DataSet) Resolve(template string) string {
	return varRe.ReplaceAllStringFunc(template, func(match string) string {
		path := match[2 : len(match)-2]
		path = strings.TrimSpace(path)
		if val, ok := ds.Get(path); ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}

func toMap(v reflect.Value) map[string]any {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	m := make(map[string]any)
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := range v.NumField() {
			m[toSnake(t.Field(i).Name)] = v.Field(i).Interface()
		}
	}
	return m
}

func toSnake(s string) string {
	var out []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, r)
	}
	return strings.ToLower(string(out))
}
