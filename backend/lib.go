package main

import (
	"bytes"
	"fmt"
	templateHTML "html/template"
	"os"
	"reflect"
	"text/template"
)

func strucToMap(data any) (map[string]any, error) {
	result := make(map[string]any)

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %T", data)
	}

	for ii := 0; ii < v.NumField(); ii++ {
		field := v.Type().Field(ii)
		value := v.Field(ii)

		// skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		result[field.Tag.Get("json")] = value.Interface()
	}

	return result, nil
}

func loadTemplate(pth string) (*template.Template, error) {
	if buf, err := os.ReadFile(pth); err != nil {
		return nil, err
	} else {
		return template.New(pth).Parse(string(buf))
	}
}

func parseTemplate(pth string, vals any) (string, error) {
	if tpl, err := loadTemplate(pth); err != nil {
		return "", err
	} else {
		var buf bytes.Buffer

		err = tpl.Execute(&buf, vals)

		return buf.String(), err
	}
}

func loadHTMLTemplate(pth string) (*templateHTML.Template, error) {
	if buf, err := os.ReadFile(pth); err != nil {
		return nil, err
	} else {
		return templateHTML.New(pth).Parse(string(buf))
	}
}

func parseHTMLTemplate(pth string, vals any) (string, error) {
	if tpl, err := loadTemplate(pth); err != nil {
		return "", err
	} else {
		var buf bytes.Buffer

		err = tpl.Execute(&buf, vals)

		return buf.String(), err
	}
}
