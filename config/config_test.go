package config

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLoadEmptyFilename(t *testing.T) {
	t.Run("Load config file with empty filename", func(t *testing.T) {
		Load("")
		if Config != nil {
			t.Error("Config should be null on empty config file")
		}
	})
}

func TestLoadNoPermissions(t *testing.T) {
	t.Run("Load config file without permissions", func(t *testing.T) {
		file, fileErr := ioutil.TempFile("", "no_permissions_file")
		if fileErr != nil {
			t.Error("Could not create test file.")
		}

		file.Chmod(0000)
		Load(file.Name())
		if Config != nil {
			t.Error("Config should be null on config file without permissions.")
		}
	})
}

func TestLoadInvalidFile(t *testing.T) {
	t.Run("Load config file - Invalid file", func(t *testing.T) {
		file, fileErr := ioutil.TempFile("", "no_read")
		if fileErr != nil {
			t.Error("Could not create test file.")
		}

		_, err := os.Open(file.Name())
		if err != nil {
			t.Error("Could not open test file.")

		}

		Load(file.Name())
		if Config != nil {
			t.Error("Config should be null on read/unmarshal error.")
		}
	})
}

func TestValidConfig(t *testing.T) {
	t.Run("Load config file - valid", func(t *testing.T) {
		file, fileErr := ioutil.TempFile("", "valid")
		if fileErr != nil {
			t.Error("Could not create test file.")
		}

		_, wErr := io.Copy(file, strings.NewReader("abc: def"))
		if wErr != nil {
			t.Error("Could not write to test file.")
		}
		Load(file.Name())
		if Config == nil {
			t.Error("Config should not be null on valid content.")
		}
	})
}
