package ceramic

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/source"
	"github.com/samber/lo"
)

var (
	ErrorGraphQLSyntaxError = fmt.Errorf("GraphQLError [Object]: Syntax Error")
)

func ModelNameSyntaxCheck(name string) error {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !re.MatchString(name) {
		return fmt.Errorf("syntax error with input name")
	}
	return nil
}

func ExtractModelName(schema []byte) (name string, err error) {
	var obj *ast.ObjectDefinition
	if obj, err = ExtarctSingleObject(schema); err != nil {
		return
	}
	name = obj.Name.Value
	return
}

func ExtarctSingleObject(schema []byte) (obj *ast.ObjectDefinition, err error) {
	doc, err := parser.Parse(parser.ParseParams{
		Source:  source.NewSource(&source.Source{Body: schema}),
		Options: parser.ParseOptions{},
	})
	if err != nil {
		return
	}
	if len(doc.Definitions) != 1 {
		err = fmt.Errorf("schema should have only one object")
		return
	}
	var ok bool
	if obj, ok = doc.Definitions[0].(*ast.ObjectDefinition); !ok {
		err = fmt.Errorf("cannot parse schema as a graphql obj")
		return
	}
	return
}

func OriginModifyFn(src string) string {
	return src
}

// for single model only
func FormatSchema(src []byte) (string, error) {
	res, err := SchemaModifyFn(src, OriginModifyFn, OriginModifyFn)
	return string(res), err
}

// for single model only
func SchemaModifyFn(body []byte, nameMod, despMod func(old string) string) (result string, err error) {
	var obj *ast.ObjectDefinition
	if obj, err = ExtarctSingleObject(body); err != nil {
		return
	}
	obj.Name.Value = nameMod(obj.Name.Value)
	for _, v := range obj.Directives[0].Arguments {
		if v.Name.Value == "description" {
			if str, ok := v.Value.(*ast.StringValue); ok {
				str.Value = despMod(str.Value)
			}
		}
	}
	result = printer.Print(obj).(string)
	return
}

func AddCustomField(body []byte, field *ast.FieldDefinition) (result string, err error) {
	var obj *ast.ObjectDefinition
	if obj, err = ExtarctSingleObject(body); err != nil {
		return
	}
	for _, v := range obj.Fields {
		if v.Name.Value == field.Name.Value {
			err = fmt.Errorf("input schema %s cannot with field named %s",
				lo.Must(json.Marshal(string(body))), field.Name.Value)
		}
	}
	obj.Fields = append(obj.Fields, field)
	result = printer.Print(obj).(string)
	return
}
