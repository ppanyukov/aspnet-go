package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// jsonConfigFileParser implements .NET equivalent of JsonConfigurationFileParser.
//
// See: https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.Json/src/JsonConfigurationFileParser.cs
type jsonConfigFileParser struct {
	data  map[string]string
	paths stringStack
}

func newJsonConfigFileParser() *jsonConfigFileParser {
	return &jsonConfigFileParser{
		data: make(map[string]string),
	}
}

func (j *jsonConfigFileParser) parseJson(r io.Reader) (map[string]interface{}, error) {
	// strip out comments as these are common
	b := bytes.Buffer{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Give a friendly message if json input is an array.
	var rootElement interface{}
	err := json.Unmarshal(b.Bytes(), &rootElement)
	if err != nil {
		return nil, err
	}

	switch val := rootElement.(type) {
	case []interface{}:
		return nil, errors.New("arrays are not supported as root json object")
	case map[string]interface{}:
		err := json.Unmarshal(b.Bytes(), &rootElement)
		if err != nil {
			return val, err
		}
		return val, nil
	default:
		return nil, errors.New("unexpected, input json is neither array nor an object")
	}
}

func (j *jsonConfigFileParser) Parse(r io.Reader) (map[string]string, error) {
	rootElement, err := j.parseJson(r)
	if err != nil {
		return j.data, err
	}

	//data, err := ioutil.ReadAll(r)
	//if err != nil {
	//	return j.data, err
	//}
	//
	//var rootElement map[string]interface{}
	//err = json.Unmarshal(data, &rootElement)
	//if err != nil {
	//	return j.data, err
	//}

	err = j.visitElement(rootElement)
	return j.data, err
}

func (j *jsonConfigFileParser) visitElement(element map[string]interface{}) error {
	isEmpty := true

	for k, v := range element {
		isEmpty = false
		j.enterContext(k)
		err := j.visitValue(v)
		if err != nil {
			return err
		}
		j.exitContext()
	}

	if isEmpty && j.paths.Count() > 0 {
		j.data[j.paths.Peek()] = ""
	}

	return nil
}

func (j *jsonConfigFileParser) enterContext(k string) {
	path := k
	if j.paths.Count() > 0 {
		path = fmt.Sprintf("%s%s%s", j.paths.Peek(), keyDelimiter, path)
	}

	path = normalizeKey(path)
	j.paths.Push(path)
}

func (j *jsonConfigFileParser) exitContext() {
	j.paths.Pop()
}

func (j *jsonConfigFileParser) visitValue(v interface{}) error {
	switch val := v.(type) {
	case map[string]interface{}:
		err := j.visitElement(val)
		if err != nil {
			return err
		}
	case []interface{}:
		for index, arrayElement := range val {
			j.enterContext(fmt.Sprintf("%d", index))
			err := j.visitValue(arrayElement)
			if err != nil {
				return err
			}
			j.exitContext()
		}
	default:
		key := j.paths.Peek()
		if _, found := j.data[key]; found {
			return errors.Errorf("duplicate key '%s'", key)
		}

		value := fmt.Sprintf("%v", val)
		if val == nil {
			value = ""
		}
		j.data[key] = value
	}

	return nil
}
