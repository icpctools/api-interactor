package interactor

import (
	"archive/zip"
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

	_ json.Marshaler   = new(ApiTime)
	_ json.Unmarshaler = new(ApiTime)
	_ fmt.Stringer     = new(ApiTime)

	_ json.Unmarshaler = new(ApiRelTime)
	_ fmt.Stringer     = new(ApiRelTime)

	_ json.Marshaler = new(LocalFileReference)
)

func TestApiTime_UnmarshalJSON(t *testing.T) {
	var ti ApiTime

	// Supported formats
	formats := []string{"2006-01-02T15:04:05Z07", time.RFC3339}
	now := time.Now()

	jsonFormat := `"%v"`
	for _, f := range formats {
		t.Run(f, func(t *testing.T) {
			jsonString := fmt.Sprintf(jsonFormat, now.Format(f))
			assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
			assert.EqualValues(t, now.Truncate(time.Second).UnixNano(), ti.Time().UnixNano())

			// Also test when the value is null
			jsonString = `null`
			assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
			assert.EqualValues(t, time.Time{}.UnixNano(), ti.Time().UnixNano())
		})
	}
}

func TestApiTime_MarshalJSON(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		n := ApiTime(time.Now())
		bts, err := json.Marshal(n)
		assert.Nil(t, err)
		assert.NotNil(t, bts)

		assert.EqualValues(t, fmt.Sprintf(`"%v"`, n.Time().Format(time.RFC3339Nano)), bts)
	})

	t.Run("empty", func(t *testing.T) {
		// Create value with a reference, this should be marshallable to an empty object
		ti := struct {
			T *ApiTime `json:"T,omitempty"`
		}{}

		bts, err := json.Marshal(ti)
		assert.Nil(t, err)
		assert.EqualValues(t, "{}", bts)
	})
}

func TestApiRelTime_UnmarshalJSON(t *testing.T) {
	ti := struct {
		T ApiRelTime `json:"T,omitempty"`
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

	// Test marshalling of an empty type
	t.Run("marshal-empty", func(t *testing.T) {
		var rt ApiTime
		bts, err := json.Marshal(rt)
		assert.Nil(t, err)
		assert.EqualValues(t, "null", bts)
	})
}

func TestLocalFileReference_MarshalJSON(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fr := new(LocalFileReference)

		// Test empty file reference
		data, err := json.Marshal(fr)
		assert.Nil(t, err)
		assert.NotNil(t, data)

		var base64zip string
		assert.Nil(t, json.Unmarshal(data, &base64zip))

		// Base64 decode the data
		decoded, err := base64.StdEncoding.DecodeString(base64zip)
		assert.Nil(t, err)
		assert.NotEmpty(t, decoded)

		// Read the ZIP file
		zipped, err := zip.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
		assert.Nil(t, err)
		assert.NotNil(t, zipped)

		// Check if it is empty
		assert.EqualValues(t, 0, len(zipped.File))
	})

	t.Run("filled", func(t *testing.T) {
		fr := new(LocalFileReference)

		// Add some files, test it again
		err := fr.FromString("sample.txt", "This is a sample")
		assert.Nil(t, err)

		goModFile, err := os.Open("go.mod")
		assert.Nil(t, err)
		defer goModFile.Close()

		err = fr.FromFile(goModFile)
		assert.Nil(t, err)

		// Test marshsalling
		var base64zip string

		data, err := json.Marshal(fr)
		assert.Nil(t, err)
		assert.Nil(t, json.Unmarshal(data, &base64zip))

		// Base64 decode the data
		decoded, err := base64.StdEncoding.DecodeString(base64zip)
		assert.Nil(t, err)
		assert.NotEmpty(t, decoded)
		// Read the ZIP file
		zipped, err := zip.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
		assert.Nil(t, err)
		assert.NotNil(t, zipped)
		// Check that it contains 2 files
		assert.EqualValues(t, 2, len(zipped.File))

		// Check that the sample.txt file is correct
		zipContentFile, err := zipped.Open("sample.txt")
		assert.Nil(t, err)
		fileContent, err := ioutil.ReadAll(zipContentFile)
		assert.Nil(t, err)
		assert.EqualValues(t, "This is a sample", fileContent)

		// Seek back to the beginning of goModFile to remove need for reopen
		_, err = goModFile.Seek(0, 0)
		assert.Nil(t, err)

		// Also check the go.mod file is correct
		goModContents, err := ioutil.ReadAll(goModFile)
		assert.Nil(t, err)
		zipContentFile, err = zipped.Open("go.mod")
		assert.Nil(t, err)
		fileContent, err = ioutil.ReadAll(zipContentFile)
		assert.Nil(t, err)
		assert.EqualValues(t, goModContents, fileContent)
	})
}
