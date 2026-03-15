package gluetun

import (
	"strings"

	z "github.com/Oudwins/zog"
)

type gluetunEnv struct {
	_name        string
	validateFunc func(value string) z.ZogIssueList
}

func newGluetunEnv(name string, schema *z.StringSchema[string]) gluetunEnv {
	return newStringGluetunEnv(name, schema)
}

func newStringGluetunEnv[T ~string](name string, schema *z.StringSchema[T]) gluetunEnv {
	return newValidatedGluetunEnv(name, func(value string) z.ZogIssueList {
		var dest T
		return schema.Parse(normalizeGluetunEnvValue(value), &dest)
	})
}

func newBoolGluetunEnv[T ~bool](name string, schema *z.BoolSchema[T]) gluetunEnv {
	return newValidatedGluetunEnv(name, func(value string) z.ZogIssueList {
		var dest T
		return schema.Parse(normalizeGluetunEnvValue(value), &dest)
	})
}

func newNumberGluetunEnv[T z.Numeric](name string, schema *z.NumberSchema[T]) gluetunEnv {
	return newValidatedGluetunEnv(name, func(value string) z.ZogIssueList {
		var dest T
		return schema.Parse(normalizeGluetunEnvValue(value), &dest)
	})
}

func newCSVGluetunEnv[T any](name string, schema *z.SliceSchema) gluetunEnv {
	return newValidatedGluetunEnv(name, func(value string) z.ZogIssueList {
		dest := new([]T)
		return schema.Parse(splitCSVGluetunEnvValue(value), dest)
	})
}

func newValidatedGluetunEnv(name string, validateFunc func(value string) z.ZogIssueList) gluetunEnv {
	return gluetunEnv{
		_name:        name,
		validateFunc: validateFunc,
	}
}

func (e gluetunEnv) name() string {
	return e._name
}

func (e gluetunEnv) validate(value string) z.ZogIssueList {
	return e.validateFunc(value)
}

func normalizeGluetunEnvValue(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return trimmed
}

func splitCSVGluetunEnvValue(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	parts := strings.Split(trimmed, ",")
	items := make([]any, len(parts))
	for i, part := range parts {
		items[i] = strings.TrimSpace(part)
	}

	return items
}
