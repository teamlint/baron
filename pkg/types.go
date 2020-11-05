package pkg

var (
	GoBaseTypes = []string{
		"uint32",
		"uint64",
		"int32",
		"int64",
		"bool",
		"float32",
		"float64",
		"string",
	}
)

// IsGoBaseType 判断是否为 Go 基本数据类型
func IsGoBaseType(goType string) bool {
	for _, v := range GoBaseTypes {
		if v == goType {
			return true
		}
	}
	return false
}
