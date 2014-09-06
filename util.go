// Package tplutil provides more convient way to use text/template inside
// the code.
//
// Consider (using text/template):
//
//	var myTpl = template.Must(template.New("name").Parse(
//		"Some list:\n" +
//			"{{range .}}" +
//			"# {{.}}\n" +
//			"{{end}}"))
//
// `gofmt` will ruin any attempt to format code above.
//
// And with tplutil:
//
//	var myTpl = template.Must(template.New("name").Parse(tplutil.Strip(`
//		Some list:{{"\n"}}
//
//		{{range .}}
//			# {{.}}{{"\n"}}
//		{{end}}
//	`)))
//
// Output will be exactly the same.
//
// Any indenting whitespaces and newlines will be ignored. If must, they
// should be specified by using syntax
//
//	`{{" "}}` or `{{"\n"}}`.
//
// It also provide `{{last}}` function to check on last element of pipeline:
//
//	var myTpl = template.Must(template.New("asd").Funcs(tplutil.Last).Parse(
//		tplutil.Strip(`
//			Some list:{{"\n"}}
//
//			{{range $i, $_ := .}}
//				{{.}}
//				{{if last $i $ | not}}
//					{{"\n"}}{{/* do not append newline to the last element */}}
//				{{end}}
//			{{end}}
//		`)
//
package tplutil

import (
	"bytes"
	"reflect"
	"regexp"
	"text/template"
)

var reInsignificantWhitespace = regexp.MustCompile(`(?m)\n?^\s*`)

var Last = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
}

func Strip(text string) string {
	return reInsignificantWhitespace.ReplaceAllString(text, ``)
}

// ExecuteToString applies a parsed template to specified data object and
// returns it output as return value. It can return partial result if
// execution can't be proceed because of error.
func ExecuteToString(tpl *template.Template, v interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, v)

	return buf.String(), err
}
