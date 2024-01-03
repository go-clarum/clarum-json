package comparator

import (
	"testing"
)

func TestIgnoreString(t *testing.T) {
	expectedValue := []byte("{" +
		"\"address\": \"@ignore@\"," +
		"\"name\": \"Bruce Wayne\"" +
		"}")
	actualValue := []byte("{" +
		"\"address\": \"some ignore value\"," +
		" \"name\": \"Bruce Wayne\"" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"address\":  <-- ignoring field\n" +
		"  \"name\": Bruce Wayne,\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}

func TestIgnoreBoolean(t *testing.T) {
	expectedValue := []byte("{" +
		"\"active\": \"@ignore@\"" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"active\":  <-- ignoring field\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}

func TestIgnoreNumberType(t *testing.T) {
	expectedValue := []byte("{" +
		" \"age\": \"@ignore@\"" +
		"}")
	actualValue := []byte("{" +
		" \"age\": 38" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"age\":  <-- ignoring field\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}

func TestIgnoreObjectType(t *testing.T) {
	expectedValue := []byte("{" +
		" \"location\": \"@ignore@\"" +
		"}")
	actualValue := []byte("{" +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"location\":  <-- ignoring field\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}

func TestIgnoreArrayType(t *testing.T) {
	expectedValue := []byte("{" +
		" \"aliases\": \"@ignore@\"" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"aliases\":  <-- ignoring field\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}

func TestIgnoreValueInArray(t *testing.T) {
	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"@ignore@\"," +
		"\"Batman\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"123," +
		"\"Batman\"" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"aliases\": [\n" +
		"     <-- ignoring field\n" +
		"    Batman,\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}

func TestIgnoreValueInRootArray(t *testing.T) {
	expectedValue := []byte("[" +
		"\"@ignore@\"," +
		"\"Batcave\"" +
		"]")
	actualValue := []byte("[" +
		"123," +
		"\"Batcave\"" +
		"]")

	expectedRecorderLog := "[\n" +
		"   <-- ignoring field\n" +
		"  Batcave,\n" +
		"]\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}
