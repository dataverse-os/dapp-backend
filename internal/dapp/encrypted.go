package dapp

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/source"
	"github.com/samber/lo"
)

var encryptableTypes = map[string]struct{}{
	"String!":    {},
	"[String!]":  {},
	"[String!]!": {},
	"String":     {},
	"[String]":   {},
	"[String]!":  {},
}

var encryptableProtectedFields = map[string]struct{}{
	"encrypted": {},
}

var EncryptedField = ast.FieldDefinition{
	Kind: kinds.FieldDefinition,
	Name: &ast.Name{
		Kind:  kinds.Name,
		Value: "encrypted",
	},
	Type: &ast.Named{
		Kind: kinds.Named,
		Name: &ast.Name{
			Kind:  kinds.Name,
			Value: "String",
		},
	},
	Directives: []*ast.Directive{
		{
			Kind: kinds.Directive,
			Name: &ast.Name{
				Kind:  kinds.Name,
				Value: "string",
			},
			Arguments: []*ast.Argument{
				{
					Kind: kinds.Argument,
					Name: &ast.Name{
						Kind:  kinds.Name,
						Value: "maxLength",
					},
					Value: &ast.IntValue{
						Kind:  kinds.IntValue,
						Value: "300000000",
					},
				},
			},
		},
	},
}

func CheckEncryptable(schema StreamModel) (err error) {
	if schema.Encryptable != nil && schema.IsPublicDomain {
		err = errors.New("input model is public but want encrypt fields")
		return
	}
	if schema.Encryptable == nil {
		return
	}
	// check protected fields, just encrypted now
	var field string
	if lo.ContainsBy(schema.Encryptable, func(item string) bool {
		_, exists := encryptableProtectedFields[item]
		if exists {
			field = item
		}
		return exists
	}) {
		err = fmt.Errorf("input model %s with protected encryptable field %s", string(schema.Schema), field)
		return
	}
	// check input encryptable with duplicate
	if len(schema.Encryptable) != len(lo.Uniq(schema.Encryptable)) {
		err = fmt.Errorf("input model %s with duplicate encryptable fields", string(schema.Schema))
		return
	}
	// assert model schema as graphql
	// donot meanning available to deploy on ceramic
	var doc *ast.Document
	if doc, err = parser.Parse(parser.ParseParams{
		Source:  source.NewSource(&source.Source{Body: []byte(schema.Schema)}),
		Options: parser.ParseOptions{},
	}); err != nil {
		return
	}
	obj, ok := doc.Definitions[0].(*ast.ObjectDefinition)
	if !ok {
		err = errors.New("assert schema definition error")
		return
	}
	// check input model fileds contain encryptable fields
	fields := lo.Map(obj.Fields, func(d *ast.FieldDefinition, _ int) string {
		return d.Name.Value
	})
	if !lo.Every(fields, schema.Encryptable) {
		err = fmt.Errorf("input model %s donot have such fields %s in encryptable",
			obj.Name.Value,
			lo.T2(lo.Difference(fields, schema.Encryptable)).B,
		)
	}
	// check encryptable fields type is encryptable
	// ## should be string in ceramic
	for _, field := range obj.Fields {
		if lo.Contains(schema.Encryptable, field.Name.Value) {
			fieldType := printer.Print(field.Type).(string)
			_, isEncryptableType := encryptableTypes[fieldType]
			if !isEncryptableType {
				err = fmt.Errorf("field %s of model %s with unencryptable type %s, but taged encryptable",
					field.Name.Value,
					obj.Name.Value,
					fieldType,
				)
				return
			}
		}
	}
	return
}
