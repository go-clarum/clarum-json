package path

import (
	"fmt"
	"strings"
)

const RootPath = "$"

func GetObjectChildPath(pathParent string, key string) string {
	return fmt.Sprintf("%s.%s", pathParent, key)
}

func GetArrayIndexPath(pathParent string, index int) string {
	return fmt.Sprintf("%s[%d]", pathParent, index)
}

func IsRoot(path string) bool {
	if path == RootPath {
		return true
	} else {
		return false
	}
}

func IsChildOfArray(path string) bool {
	if strings.LastIndex(path, "]") == len(path)-1 {
		return true
	} else {
		return false
	}
}
