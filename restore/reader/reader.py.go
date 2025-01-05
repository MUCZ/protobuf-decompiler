package reader

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func PyDescriptorReader(filepath string) ([]byte, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filepath)
	}
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var re = regexp.MustCompile(`AddSerializedFile\(b\'([^)]*)\'\)`)
	match := re.FindStringSubmatch(string(bytes))
	if len(match) <= 1 {
		return nil, fmt.Errorf("no match found")
	}
	return pyStringToHex(match[1])
}

func pyStringToHex(s string) ([]byte, error) {
	var result []byte
	i := 0
	for i < len(s) {
		if value, _, tail, err := strconv.UnquoteChar(s[i:], '"'); err == nil {
			result = append(result, byte(value))
			i += len(string(value))
			s = s[:i] + tail
		} else {
			return nil, err
		}
	}
	return result, nil
}
