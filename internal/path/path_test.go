package path

import (
	"testing"
)

func TestGetObjectChildPath(t *testing.T) {
	result := GetObjectChildPath("object", "field")

	if result != "object.field" {
		t.Error("wrong path: " + result)
	}
}

func TestGetArrayIndexPath(t *testing.T) {
	result := GetArrayIndexPath("myarray", 1)

	if result != "myarray[1]" {
		t.Error("wrong path: " + result)
	}
}

func TestIsRoot(t *testing.T) {
	trueResult := IsRoot(RootPath)

	if !trueResult {
		t.Error("should be root")
	}

	falseResult := IsRoot("some.thing")

	if falseResult {
		t.Error("should not be root")
	}
}

func TestIsChildOfArray(t *testing.T) {
	trueResult := IsChildOfArray("myarray[3]")

	if !trueResult {
		t.Error("should be true")
	}

	falseResult := IsChildOfArray("myarray.")

	if falseResult {
		t.Error("should not false")
	}
}
