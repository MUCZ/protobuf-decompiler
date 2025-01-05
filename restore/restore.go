package restore

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mucz/protobuf-decompiler/restore/reader"
	"google.golang.org/protobuf/proto"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

var (
	readerMap = map[string]func(string) ([]byte, error){
		".go":   reader.GoDescriptorReader,
		".py":   reader.PyDescriptorReader,
		".txt":  reader.TxtRawDescReader,
		".cc":   reader.CcDescriptorReader,
		".cpp":  reader.CcDescriptorReader,
		".c":    reader.CcDescriptorReader,
		".h":    reader.CcDescriptorReader,
		".java": reader.JavaDescriptorReader,
		".rs":   reader.RsDescriptorReader,
		".cs":   reader.CsDescriptorReader,
	}
)

func Do(file string) (string, error) {
	fileExtension := filepath.Ext(file)
	reader, ok := readerMap[fileExtension]
	if !ok {
		return "", fmt.Errorf("unsupported file extension: %s", fileExtension)
	}
	bytes, err := reader(file)
	if err != nil {
		return "", err
	}
	data, _, err := restoreSingleProtoFile(bytes)
	if err != nil {
		return "", err
	}
	return data, nil
}

func restoreSingleProtoFile(rawDesc []byte) (string, string, error) {
	if len(rawDesc) == 0 {
		return "", "", nil
	}
	fileDesc := &descriptorpb.FileDescriptorProto{}
	if err := proto.Unmarshal(rawDesc, fileDesc); err != nil {
		return "", "", err
	}
	str := renderProtoFile(fileDesc)
	lines := strings.Split(str, "\n")
	ret := ""
	for _, line := range lines {
		cleanLine := strings.TrimSpace(line)
		cleanLine = strings.TrimPrefix(cleanLine, "\t")
		cleanLine = strings.TrimSuffix(cleanLine, "\t")
		if cleanLine != "" {
			if line == "🏷️" { // empty line marker
				line = ""
			}
			ret += line + "\n"
			if line == "}" {
				ret += "\n"
			}
		}
	}
	return ret, *fileDesc.Name, nil
}

func normalizeMessageName(name string, parentPkg string) string {
	name = strings.TrimPrefix(name, "."+parentPkg+".")
	name = strings.TrimPrefix(name, ".")
	return name
}

func getTypeNameAndGenre(f *descriptorpb.FieldDescriptorProto) (string, bool) {
	fieldType := f.Type.String()
	if fieldType == "TYPE_GROUP" || fieldType == "TYPE_MESSAGE" {
		// nested type
		return *f.TypeName, true
	} else if fieldType == "TYPE_ENUM" {
		return *f.TypeName, false
	} else {
		// simple type
		return TypeStringMap[fieldType], false
	}
}
