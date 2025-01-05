package restore

const fileStencil = `
// {{.Name}}
// File restored by protodec.
// DO NOT EDIT!
🏷️
{{$FileDesc := . -}}

{{if gt (len .Syntax) 0}}syntax = "{{.Syntax}}";{{end}}
{{if gt (len .Package) 0}}package {{.Package}};{{end}}
🏷️
{{range $index, $import := .Dependency -}}
import "{{$import}}";
{{end}}
{{if gt (len .Dependency) 0}}🏷️{{end}}
{{with .Options}}
{{with .GoPackage}}
{{if gt (len .) 0}}option go_package = "{{ . }}";{{end}}
{{end}}
{{end}}
🏷️

{{extend .}}
{{ if gt (len .Service) 0 -}}
// Services
🏷️
{{range $index, $service := .Service -}}
service {{$service.Name}} {
{{- range $index, $method := $service.Method }}
	{{rpc $method $FileDesc }}
{{- end }}
}
{{end}}
{{- end -}}

{{- if gt (len .MessageType) 0 }}

// Messages
🏷️
{{range $index, $message := .MessageType -}}
{{message $message 0}}
{{ end -}}
{{- end -}}

{{- if gt (len .EnumType) 0 }}

// Enums
🏷️
{{range $index, $enum := .EnumType -}}
{{enum $enum 0 }}
{{ end -}}
{{- end -}}
`

var (
	TypeStringMap = map[string]string{
		"TYPE_DOUBLE":   "double",
		"TYPE_FLOAT":    "float",
		"TYPE_INT64":    "int64",
		"TYPE_UINT64":   "uint64",
		"TYPE_INT32":    "int32",
		"TYPE_FIXED64":  "fixed64",
		"TYPE_FIXED32":  "fixed32",
		"TYPE_BOOL":     "bool",
		"TYPE_STRING":   "string",
		"TYPE_BYTES":    "bytes",
		"TYPE_UINT32":   "uint32",
		"TYPE_SFIXED32": "sfixed32",
		"TYPE_SFIXED64": "sfixed64",
		"TYPE_SINT32":   "sint32",
		"TYPE_SINT64":   "sint64",
	}
)
