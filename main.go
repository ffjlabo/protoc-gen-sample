package main

import (
	"bytes"
	"flag"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

type Field struct {
	Name string
	Type string
}

type Struct struct {
	Name   string
	Fields []Field
}

func generateFile(gen *protogen.Plugin, protoFile *protogen.File) *protogen.GeneratedFile {
	filename := protoFile.GeneratedFilenamePrefix + ".go"
	g := gen.NewGeneratedFile(filename, protoFile.GoImportPath)

	// 構造体を元にコードを生成
	// fmt.Println(*f.Proto.Package)
	file := NewFile(*protoFile.Proto.Package)

	// Messageから構造体を作成
	var Structs []Struct

	for _, m := range protoFile.Messages {
		strct := Struct{
			Name:   m.GoIdent.GoName,
			Fields: []Field{},
		}

		for _, field := range m.Fields {
			ff := Field{
				Name: field.GoName,
				Type: field.Desc.Kind().String(),
			}
			strct.Fields = append(strct.Fields, ff)
		}

		Structs = append(Structs, strct)
	}

	for _, strct := range Structs {
		var codes []Code
		for _, field := range strct.Fields {
			code := Id(field.Name).Id(field.Type)
			codes = append(codes, code)
		}

		// modelのstruct作成
		file.Type().Id(strct.Name).Struct(codes...)
	}

	buf := &bytes.Buffer{}
	err := file.Render(buf)
	if err != nil {
		panic(err)
	}

	g.P(buf)

	return g
}

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			generateFile(gen, f)
		}

		return nil
	})
}
