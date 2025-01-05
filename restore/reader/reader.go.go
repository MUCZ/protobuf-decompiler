package reader

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GoDescriptorReader(filepath string) ([]byte, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filepath)
	}
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var re = regexp.MustCompile(`var .*_rawDesc = \[\]byte\{([^}]*)\}`)
	match := re.FindStringSubmatch(string(bytes))
	if len(match) <= 1 {
		return nil, fmt.Errorf("no match found")
	}
	return goStrToHex(match[1]), nil
}

// txt reader is used for tests.
func TxtRawDescReader(filepath string) ([]byte, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filepath)
	}
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return goStrToHex(string(bytes)), nil
}

func goStrToHex(str string) []byte {
	parts := strings.Split(str, ",")
	strBytes := make([]string, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		s := parts[i]
		s = strings.TrimSpace(s)
		if s != "" {
			strBytes = append(strBytes, s)
		}
	}
	value := make([]byte, 0, len(strBytes))
	for _, part := range strBytes {
		part = strings.TrimSpace(part)
		val, _ := strconv.ParseUint(part, 0, 8)
		value = append(value, byte(val))
	}
	return value
}
