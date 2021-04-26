package interactor

import (
	zip2 "archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// Ensure all types adhere to required interfaces
var (
	_ ApiType = Contest{}
	_ ApiType = Problem{}
	_ ApiType = Submission{}
	_ ApiType = Clarification{}
	_ ApiType = Language{}

	_ Submittable = Clarification{}

	_ LocalFileReference = &localFileReference{}

	_ json.Marshaler   = new(ApiTime)
	_ json.Unmarshaler = new(ApiTime)
	_ fmt.Stringer     = new(ApiTime)

	_ json.Unmarshaler = new(ApiRelTime)
	_ fmt.Stringer     = new(ApiRelTime)

	_ json.Marshaler = new(localFileReference)
)

func TestApiTime_UnmarshalJSON(t *testing.T) {
	ti := struct {
		T ApiTime
	}{}

	// Supported formats
	formats := []string{"2006-01-02T15:04:05Z07", time.RFC3339}
	now := time.Now()

	jsonFormat := `{"T": "%v"}`
	for _, f := range formats {
		t.Run("format-"+f, func(t *testing.T) {
			jsonString := fmt.Sprintf(jsonFormat, now.Format(f))
			assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
			assert.EqualValues(t, now.Truncate(time.Second).UnixNano(), ti.T.Time().UnixNano())

			// Also test when the value is null
			jsonString = `{"T": null}`
			assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
			assert.EqualValues(t, time.Time{}.UnixNano(), ti.T.Time().UnixNano())
		})
	}
}

func TestApiRelTime_UnmarshalJSON(t *testing.T) {
	ti := struct {
		T ApiRelTime
	}{}

	// Only one format is allowed, test a single value. TODO perhaps some form of fuzz testing here?
	jsonString := `{"T": "0:03:38.749"}`
	duration := time.Minute*3 + time.Second*38 + time.Millisecond*749
	assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
	assert.EqualValues(t, duration, ti.T.Duration())

	// Test null
	jsonString = `{"T": null}`
	assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
	assert.EqualValues(t, time.Duration(0), ti.T.Duration())
}

func TestLocalFileReference_MarshalJSON(t *testing.T) {
	ti := struct {
		T LocalFileReference
	}{}

	ti.T = NewLocalFileReference()

	// Test empty file reference
	data, err := json.Marshal(ti)
	assert.Nil(t, err)
	// Now decode it as JSON again
	result := struct {
		T string
	}{}
	assert.Nil(t, json.Unmarshal(data, &result))
	// Base64 decode the data
	decoded, err := base64.StdEncoding.DecodeString(result.T)
	assert.Nil(t, err)
	assert.NotEmpty(t, decoded)
	// Read the ZIP file
	zip, err := zip2.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
	assert.Nil(t, err)
	assert.NotNil(t, zip)
	// Check if it is empty
	assert.EqualValues(t, 0, len(zip.File))

	ti.T = NewLocalFileReference()

	// Add some files, test it again
	err = ti.T.AddFromString("sample.txt", "This is a sample")
	assert.Nil(t, err)
	goModFile, err := os.Open("go.mod")
	assert.Nil(t, err)
	err = ti.T.AddFromFile(goModFile)
	assert.Nil(t, err)

	data, err = json.Marshal(ti)
	assert.Nil(t, err)
	assert.Nil(t, json.Unmarshal(data, &result))
	// Base64 decode the data
	decoded, err = base64.StdEncoding.DecodeString(result.T)
	assert.Nil(t, err)
	assert.NotEmpty(t, decoded)
	// Read the ZIP file
	zip, err = zip2.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
	assert.Nil(t, err)
	assert.NotNil(t, zip)
	// Check that it contains 2 files
	assert.EqualValues(t, 2, len(zip.File))

	// Check that the sample.txt file is correct
	zipContentFile, err := zip.Open("sample.txt")
	assert.Nil(t, err)
	fileContent, err := ioutil.ReadAll(zipContentFile)
	assert.Nil(t, err)
	assert.EqualValues(t, "This is a sample", string(fileContent))

	// Also check the go.mod file is correct
	goModFile, _ = os.Open("go.mod")
	goModContents, err := ioutil.ReadAll(goModFile)
	assert.Nil(t, err)
	zipContentFile, err = zip.Open("go.mod")
	assert.Nil(t, err)
	fileContent, err = ioutil.ReadAll(zipContentFile)
	assert.EqualValues(t, string(goModContents), string(fileContent))
}
