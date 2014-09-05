// Package tplutil provides more convient way to use text/template inside
// the code.
//
// Consider (using text/template):
//
//	var myTpl = template.Must(template.New("name").Parse(
//		"Some list:\n" +
//			"{{range .}}" +
//			"# {{.}}\n" +
//			"{{end}}")
//
// `gofmt` will ruin any attempt to format code above.
//
// And with tplutil:
//
//	var myTpl = tplutil.SparseTemplate("name", `
//		Some list:{{"\n"}}
//
//		{{range .}}
//			# {{.}}{{"\n"}}
//		{{end}}
//	`)
//
// Output will be exactly the same.
//
// Any indenting whitespaces and newlines will be ignored. If must, they
// should be specified by using syntax `{{" "}}` or `{{"\n"}}`.
//
// It also provide `{{last}}` function to check on last element of pipeline:
//
//	var myTpl = tplutil.SparseTemplate("name", `
//		Some list:{{"\n"}}
//
//		{{range $i, $_ := .}}
//			{{.}}
//			{{if last $i $ | not}}
//				{{"\n"}}{{/* do not append newline to the last element */}}
//			{{end}}
//		{{end}}
//	`)
//
// Behaviour of `Execute` is changed too: it will return `string` as template
// execution result.
//
package tplutil

import (
	"bytes"
	"reflect"
	"regexp"
	"text/template"
)

type Template struct {
	*template.Template
}

var reInsignificantWhitespace = regexp.MustCompile(`(?m)\n?^\s*`)

// SparseTemplate constructs template from "sparse" variant, trimming all
// insignificant indent whitespaces and newlines.
func SparseTemplate(name, text string) *Template {
	stripped := reInsignificantWhitespace.ReplaceAllString(text, ``)

	funcs := template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
	}

	tpl := &Template{
		template.Must(template.New("comment").Funcs(funcs).Parse(stripped)),
	}

	return tpl
}

// Execute applies a parsed template to specified data object and returns it
// output as return value.
//
// Default behaviour can be anytime restored by using `t.Template.Execute()`
// call.
func (t *Template) Execute(v interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := t.Template.Execute(buf, v)

	if err != nil {
		panic(err)
	}

	return buf.String(), err
}
