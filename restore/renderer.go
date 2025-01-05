package restore

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/protocolbuffers/protoscope"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

func renderProtoFile(fileDesc *descriptorpb.FileDescriptorProto) string {
	var messageRenderer func(message *descriptorpb.DescriptorProto, depth int) string
	var messageRendererW = func(message *descriptorpb.DescriptorProto, depth int) string { return messageRenderer(message, depth) }

	var fieldRendererBoundToMdFd func(field *descriptorpb.FieldDescriptorProto) string
	var fieldRendererBoundToMdFdW = func(field *descriptorpb.FieldDescriptorProto) string { return fieldRendererBoundToMdFd(field) }

	var oneOfRendererBoundToFd = func(message *descriptorpb.DescriptorProto, depth int) string {
		return oneOfRendereer(message, fileDesc, depth)
	}
	renderersMap := template.FuncMap{
		"extend":  extendRenderer,
		"field":   fieldRendererBoundToMdFdW,
		"message": messageRendererW,
		"rpc":     rpcRenderer,
		"enum":    enumRenderer,
		"oneof":   oneOfRendererBoundToFd,
	}

	var render = func(name, stencil string, data interface{}) string {
		t := template.Must(template.New(name).Funcs(renderersMap).Parse(stencil))
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, data); err != nil {
			panic(err)
		}
		return tpl.String()
	}

	messageRenderer = func(messageDesc *descriptorpb.DescriptorProto, depth int) string {
		if messageDesc.Options != nil && (messageDesc.Options.MapEntry != nil && *messageDesc.Options.MapEntry) {
			return ""
		}

		fieldRendererBoundToMdFd = func(field *descriptorpb.FieldDescriptorProto) string {
			if field.OneofIndex != nil { // render oneof field outside
				return ""
			}
			return fieldRenderer(field, messageDesc, fileDesc)
		}

		messageStencil := `
message {{$.Name}} {
{{- range $index, $field := $.Field }}
	{{field $field}}
{{- end }}
	{{- oneof . ` + strconv.Itoa(depth+1) + ` -}}
{{range $index, $nestedType := $.NestedType -}}
	{{- message $nestedType ` + strconv.Itoa(depth+1) + ` -}}
{{- end }}
{{range $index, $enum := .EnumType }}
	{{- enum $enum ` + strconv.Itoa(depth+1) + ` -}}
{{ end }}
}`
		messageStencil = prependTabs(messageStencil, depth)
		result := render("message-"+*messageDesc.Name, messageStencil, messageDesc)
		if messageDesc.Options == nil {
			return result
		}
		opt := parseRawOptionsInUnknownFields(messageDesc.Options.ProtoReflect().GetUnknown())
		return fmt.Sprintf("//? option: [%s]\n%s", opt, result)
	}

	return render("file-"+*fileDesc.Name, fileStencil, fileDesc)
}

func rpcRenderer(methodDesc *descriptorpb.MethodDescriptorProto, fileDesc *descriptorpb.FileDescriptorProto) string {
	name, input, output := *methodDesc.Name, *methodDesc.InputType, *methodDesc.OutputType
	input = normalizeMessageName(input, *fileDesc.Package)
	output = normalizeMessageName(output, *fileDesc.Package)
	return fmt.Sprintf("rpc %s (%s) returns (%s) {}", name, input, output)
}

func enumRenderer(enumDesc *descriptorpb.EnumDescriptorProto, depth int) string {
	ret := "\nenum " + *enumDesc.Name + " {"
	for _, item := range enumDesc.Value {
		ret += fmt.Sprintf("\n\t%s = %d;", *item.Name, *item.Number)
	}
	ret += "\n}"
	return prependTabs(ret, depth)
}

func extendRenderer(file *descriptorpb.FileDescriptorProto) string {
	if len(file.Extension) == 0 {
		return ""
	}
	ret := ""
	type extend struct {
		extendee string
		fields   []*descriptorpb.FieldDescriptorProto
	}
	extends := make([]*extend, 0) // use slice to assure the order
	for _, ext := range file.Extension {
		found := false
		for _, x := range extends {
			if x.extendee == *ext.Extendee {
				x.fields = append(x.fields, ext)
				found = true
				break
			}
		}
		if !found {
			extends = append(extends, &extend{
				extendee: *ext.Extendee,
				fields:   []*descriptorpb.FieldDescriptorProto{ext},
			})
		}
	}
	for _, x := range extends {
		extendee, exts := x.extendee, x.fields
		ret += fmt.Sprintf("extend %s {", normalizeMessageName(extendee, *file.Package))
		for _, ext := range exts {
			typeName, isComplexType := getTypeNameAndGenre(ext)
			if isComplexType {
				typeName = *ext.TypeName
			}
			ret += fmt.Sprintf("\n\t%s %s = %d;", normalizeMessageName(typeName, *file.Package), *ext.Name, *ext.Number)
		}
		ret += "\n}\n"
	}
	return ret
}

func oneOfRendereer(messageDesc *descriptorpb.DescriptorProto, fileDesc *descriptorpb.FileDescriptorProto, depth int) string {
	if len(messageDesc.OneofDecl) == 0 {
		return ""
	}
	ret := ""
	for i, oneOf := range messageDesc.OneofDecl {
		ret += fmt.Sprintf("\noneof %s {\n", *oneOf.Name)
		for _, field := range messageDesc.Field {
			if field.OneofIndex != nil && *field.OneofIndex == int32(i) {
				ret += fmt.Sprintf("\t%s\n", fieldRenderer(field, messageDesc, fileDesc))
			}
		}
		ret += "}"
	}
	return prependTabs(ret, depth)
}

func fieldRenderer(field *descriptorpb.FieldDescriptorProto, messageDesc *descriptorpb.DescriptorProto, fileDesc *descriptorpb.FileDescriptorProto) string {
	ret := ""
	getTypeName := func(f *descriptorpb.FieldDescriptorProto) string {
		typeName, _ := getTypeNameAndGenre(f)
		return normalizeMessageName(typeName, *fileDesc.Package)
	}
	if *field.Label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		ret += "repeated "
	}
	typeName, isComplexType := getTypeNameAndGenre(field)
	if isComplexType {
		typeName = *field.TypeName

		for _, nestedType := range messageDesc.NestedType {
			if nestedType.Options != nil && *nestedType.Options.MapEntry {
				nestedTypeName := *nestedType.Name
				parts := strings.Split(typeName, ".")
				if parts[len(parts)-1] == nestedTypeName {
					typeName = fmt.Sprintf("map<%s, %s>", getTypeName(nestedType.Field[0]), getTypeName(nestedType.Field[1]))
					ret = ""
				}
			}
		}
	}
	typeName = normalizeMessageName(typeName, *fileDesc.Package)
	ret += fmt.Sprintf("%s %s = %d;", typeName, *field.Name, *field.Number)
	if field.Options == nil {
		return ret
	}
	rawOptions := parseRawOptionsInUnknownFields(field.Options.ProtoReflect().GetUnknown())
	return ret + fmt.Sprintf(" //? option: [%s]", rawOptions)
}

func prependTabs(s string, n int) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.Repeat("\t", n) + lines[i]
	}
	return strings.Join(lines, "\n")
}

func parseRawOptionsInUnknownFields(unknownFields []byte) string {
	if unknownFields == nil {
		return ""
	}
	r := protoscope.Write(unknownFields, protoscope.WriterOptions{})
	r = strings.ReplaceAll(r, "\n", " ")
	r = strings.ReplaceAll(r, "\t", " ")
	r = strings.TrimSpace(r)
	return r
}
