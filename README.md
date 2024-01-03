# Clarum JSON

Library used by the clarum framework to do JSON validation.

## Features

The `Comparator` validates if two JSON objects match and:

- returns detailed errors on where and how they do not match
- errors are accompanied by json paths (when one can be provided)
- allows ignoring values of fields

## How to use

The recommended way is to create a `Comparator` using the builder and then use it to compare two JSON objects:

```go
expectedValue := []byte("{" +
"\"active\": true," +
" \"name\": \"Bruce\"," +
" \"age\": 37," +
" \"height\": 1.879," +
"\"location\": {" +
"\"street\": \"Mountain Drive\"," +
"\"number\": 1007," +
"\"hidden\": false," +
"\"timestamp\": \"@ignore@\"" +
"}" +
"}")

actualValue := []byte("{" +
"\"active\": true," +
" \"name\": \"Bruce Wayne\"," +
" \"age\": 38," +
" \"height\": 1.879," +
"\"location\": {" +
"\"address\": \"Mountain Drive\"," +
"\"number\": 1008," +
"\"hidden\": true," +
"\"timestamp\": \"2024-01-03 23:42:00\"" +
"}" +
"}")

jc := comparator.NewComparator().Build()

_, err := jc.Compare(expectedValue, actualValue)

// err will contain the following validation errors:
//
// "[$.name] - value mismatch - expected [Bruce] but received [Bruce Wayne]"
// "[$.age] - value mismatch - expected [37] but received [38]"
// "[$.location.street] - field is missing"
// "[$.location.number] - value mismatch - expected [1007] but received [1008]"
// "[$.location.hidden] - value mismatch - expected [false] but received [true]"
// "[$.location.address] - unexpected field"
```

## Recorder

The `Recorder` is an optional feature that returns a user-friendly output which makes it easier to see where the
validation has failed.
By default, a NoopRecorder is configured.

In the example above we can configure a different recorder like this:

```go
jc := NewComparator().
Recorder(recorder.NewDefaultRecorder()).
Build()

recorderLog, err := jc.Compare(expectedValue, actualValue)
```

The `recorderLog` will be a string that will look like this:

```
{
  "location": {
    "timestamp":  <-- ignoring field
     X-- missing field [street]
    "number": 1008, <-- value mismatch - expected [1007]
    "hidden": true, <-- value mismatch - expected [false]
    "address":  <-- unexpected field
  },
  "active": true,
  "name": Bruce Wayne, <-- value mismatch - expected [Bruce]
  "age": 38, <-- value mismatch - expected [37]
  "height": 1.879,
}
```

## Configuration

| Key               | Default          | Description                                                                                                                                                                                                                         |
|-------------------|------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| StrictObjectCheck | `true`           | Determines if the Comparator will do a strict check on object fields<br/><br/>If set to `true`, the following checks will be done:<br/>  - actual JSON has the same number of fields<br/> - actual JSON has extra unexpected fields |
| Logger            | `slog.Default()` | Logger used internally by the Comparator                                                                                                                                                                                            |
| Recorder          | `NoopRecorder`   | Recorder implementation to be used                                                                                                                                                                                                  |

## Ignoring field values

The comparator allows you to validate if a field is present **but** ignore its value. 

For example some entities may have a `modifiedAt` field which is a timestamp. We want to validate that the field exists but cannot really predict its value.
In such a case we can use the special `@ignore@` marker as the value of the field in our `expected` JSON object.
For an example check the code from [How to use](#how-to-use).
