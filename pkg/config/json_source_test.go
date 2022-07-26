package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jsonConfigProvider_Load_Simple(t *testing.T) {
	json := `
{
	// comments are allowed
	"string": "string value",
	"number": 1,
	"float": 1.23,
	"bool": true,
	"null": null,
	"array_of_strings": [
		// comments are allowed
		"element_0",
		"element_1",
		"element_2"
	],
	"array_of_objects": [
		{
			"string": "string value 1",
			"number": 1,
			"float": 1.23,
			"bool": true,
			"null": null
		},
		{
			"string": "string value 2",
			"number": 2,
			"float": 2.23,
			"bool": true,
			"null": null
		}
	],

	"nested1": {
		"string": "nested1 string value",
		"number": 2,
		"float": 9.23,
		"bool": true,
		"null": null
	}
}
`

	provider, err := NewJsonSource([]byte(json)).Build()
	assert.NoError(t, err)

	// elements at the root
	assert.Equal(t, "string value", provider.Get("string"))
	assert.Equal(t, "1", provider.Get("number"))
	assert.Equal(t, "1.23", provider.Get("float"))
	assert.Equal(t, "true", provider.Get("bool"))
	assert.Equal(t, "", provider.Get("null"))

	// elements at the root, using different case to grab values
	assert.Equal(t, "string value", provider.Get("STRING"))
	assert.Equal(t, "1", provider.Get("NUMBER"))
	assert.Equal(t, "1.23", provider.Get("FLOAT"))
	assert.Equal(t, "true", provider.Get("BOOL"))
	assert.Equal(t, "", provider.Get("NULL"))

	// array of strings
	assert.Equal(t, "", provider.Get("array_of_strings"))
	assert.Equal(t, "element_0", provider.Get("array_of_strings:0"))
	assert.Equal(t, "element_1", provider.Get("array_of_strings:1"))

	//// array of objects
	assert.Equal(t, "", provider.Get("array_of_objects"))
	assert.Equal(t, "string value 1", provider.Get("array_of_objects:0:string"))
	assert.Equal(t, "1", provider.Get("array_of_objects:0:number"))
	assert.Equal(t, "1.23", provider.Get("array_of_objects:0:float"))
	assert.Equal(t, "true", provider.Get("array_of_objects:0:bool"))
	assert.Equal(t, "", provider.Get("array_of_objects:0:null"))

	assert.Equal(t, "string value 2", provider.Get("array_of_objects:1:string"))
	assert.Equal(t, "2", provider.Get("array_of_objects:1:number"))
	assert.Equal(t, "2.23", provider.Get("array_of_objects:1:float"))
	assert.Equal(t, "true", provider.Get("array_of_objects:1:bool"))
	assert.Equal(t, "", provider.Get("array_of_objects:1:null"))

	// nested1
	assert.Equal(t, "", provider.Get("nested1"))
	assert.Equal(t, "nested1 string value", provider.Get("nested1:string"))
	assert.Equal(t, "2", provider.Get("nested1:number"))
	assert.Equal(t, "9.23", provider.Get("nested1:float"))
	assert.Equal(t, "true", provider.Get("nested1:bool"))
	assert.Equal(t, "", provider.Get("nested1:null"))
}

func Test_jsonConfigProvider_Load_ArraysNotSupported(t *testing.T) {
	// .NET does not support arrays as root objects
	json := `
[
	"element_0",
	"element_1",
	"element_2"
]
`

	_, err := NewJsonSource([]byte(json)).Build()
	if assert.Error(t, err) {
		hasCorrectMessage := strings.Contains(err.Error(), "arrays are not supported as root json object")
		assert.Truef(t, hasCorrectMessage, "array objects should return correct message, was: %v", err.Error())
	}
}
