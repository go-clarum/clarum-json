package comparator

import (
	"strings"
	"testing"
)

func TestEmptyInExpectedJson(t *testing.T) {
	expectedErrors := []string{
		"[$.aliases] - array size mismatch - expected [0] but received [1]",
	}

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"aliases\": [ <-- size mismatch - expected [0]\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestEmptyInActualJson(t *testing.T) {
	expectedErrors := []string{
		"[$.aliases] - array size mismatch - expected [1] but received [0]",
	}

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"]" +
		"}")

	expectedRecorderLog := "{\n  " +
		"\"aliases\": [ <-- size mismatch - expected [1]\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestTypeMismatchJson(t *testing.T) {
	expectedErrors := []string{
		"[$.aliases[1]] - value type mismatch - expected [string] but found [number]",
		"[$.aliases[2]] - value type mismatch - expected [string] but found [object]",
		"[$.aliases[3]] - value type mismatch - expected [string] but found [array]",
	}

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"The Dark Knight\"," +
		"\"Batsy\"," +
		"\"The Gotham Guardian\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"123," +
		"{" +
		"\"someStringKey\": \"someValue\"," +
		"\"someNumberKey\": 123" +
		"}," +
		"[1,2,3]" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"aliases\": [\n" +
		"    Batman,\n" +
		"    123, <-- value type mismatch - expected [string] but found [number]\n" +
		"    object, <-- value type mismatch - expected [string] but found [object]\n" +
		"    array, <-- value type mismatch - expected [string] but found [array]\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestStringValidation(t *testing.T) {
	expectedErrors := []string{
		"[$.aliases[1]] - value mismatch - expected [The Dark Knight] but received [Robin]",
	}

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"The Dark Knight\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"Robin\"" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"aliases\": [\n" +
		"    Batman,\n" +
		"    Robin, <-- value mismatch - expected [The Dark Knight]\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestNumberValidation(t *testing.T) {
	expectedErrors := []string{
		"[$.measures[1]] - value mismatch - expected [82] but received [83]",
		"[$.measures[3]] - value mismatch - expected [64.1] but received [64.2]",
	}

	expectedValue := []byte("{" +
		"\"measures\": [" +
		"45," +
		"82," +
		"32.2," +
		"64.1" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"measures\": [" +
		"45," +
		"83," +
		"32.2," +
		"64.2" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"measures\": [\n" +
		"    45,\n" +
		"    83, <-- value mismatch - expected [82]\n" +
		"    32.2,\n" +
		"    64.2, <-- value mismatch - expected [64.1]\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestBoolValidation(t *testing.T) {
	expectedErrors := []string{
		"[$.someBooleanArray[1]] - value mismatch - expected [true] but received [false]",
	}

	expectedValue := []byte("{" +
		"\"someBooleanArray\": [" +
		"true," +
		"true" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"someBooleanArray\": [" +
		"true," +
		"false" +
		"]" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"someBooleanArray\": [\n" +
		"    true,\n" +
		"    false, <-- value mismatch - expected [true]\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestDeepArrayValidation(t *testing.T) {
	expectedErrors := []string{
		"[$.parent[1][1]] - value type mismatch - expected [string] but found [number]",
	}

	expectedValue := []byte("{" +
		"\"parent\": [" +
		"[" +
		"\"child11\"," +
		"\"child12\"" +
		"]," +
		"[" +
		"\"child21\"," +
		"\"child22\"" +
		"]" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"parent\": [" +
		"[" +
		"\"child11\"," +
		"\"child12\"" +
		"]," +
		"[" +
		"\"child21\"," +
		"123" +
		"]" +
		"]" +
		"}")
	expectedRecorderLog := "{\n" +
		"  \"parent\": [\n" +
		"    [\n" +
		"      child11,\n" +
		"      child12,\n" +
		"    ],\n" +
		"    [\n" +
		"      child21,\n" +
		"      123, <-- value type mismatch - expected [string] but found [number]\n" +
		"    ],\n" +
		"  ],\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedErrors, expectedRecorderLog)
}

func TestObjectValidation(t *testing.T) {
	expectedErrors := []string{
		"[$.addresses[0].number] - value mismatch - expected [1007] but received [1035]",
		"[$.addresses[1].hidden] - value mismatch - expected [true] but received [false]",
	}

	expectedValue := []byte("{" +
		"\"addresses\": [" +
		"{" +
		"\"name\": \"Home\"," +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}," +
		"{" +
		"\"name\": \"Batcave\"," +
		"\"street\": \"unknown\"," +
		"\"number\": 0," +
		"\"hidden\": true" +
		"}" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"addresses\": [" +
		"{" +
		"\"name\": \"Home\"," +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1035," +
		"\"hidden\": false" +
		"}," +
		"{" +
		"\"name\": \"Batcave\"," +
		"\"street\": \"unknown\"," +
		"\"number\": 0," +
		"\"hidden\": false" +
		"}" +
		"]" +
		"}")

	// we ignore the recorder log because the order of the elements in the object is always different
	recorderResult := testComparator(t, expectedValue, actualValue, expectedErrors, "")

	if !strings.Contains(recorderResult, "\n    },\n    {\n      ") {
		t.Error("indentation between objects is wrong")
	}
	if !strings.Contains(recorderResult, "\n    },\n  ],\n}\n") {
		t.Error("indentation at the end is wrong")
	}
}

func TestRootArrayValidation(t *testing.T) {
	expectedValue := []byte("[" +
		"\"Batcave\"" +
		"]")
	actualValue := []byte("[" +
		"\"Batcave\"" +
		"]")

	expectedRecorderLog := "[\n" +
		"  Batcave,\n" +
		"]\n"

	testComparator(t, expectedValue, actualValue, []string{}, expectedRecorderLog)
}
