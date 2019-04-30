package utils

import (
	"github.com/jeremywohl/flatten"
)

func Flatten(docName string, tree map[string]interface{}) (map[string]interface{}, error) {
	root := make(map[string]interface{})
	if len(docName) > 0 {
		root[docName] = tree
	} else {
		root = tree
	}
	return flatten.Flatten(root, "", flatten.DotStyle)
}
