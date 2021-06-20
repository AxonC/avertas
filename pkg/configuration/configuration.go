package configuration // import "github.com/AxonC/avertas/configuration"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Struct to store configuration information
type Configuration struct {
	Folders []Folder `json:"folders"`
}

// Struct to store registered folder.
type Folder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func DefaultConfigPath() string {
	homeDirectory, _ := os.UserHomeDir()

	return homeDirectory + "/.avertas/config.json"
}

func ReadConfiguration() (Configuration, error) {
	file, err := os.Open(DefaultConfigPath())
	if err != nil {
		fmt.Println("Unable to find configuration")
	}
	defer file.Close()

	contents, _ := ioutil.ReadAll(file)
	var configuration Configuration
	json.Unmarshal([]byte(contents), &configuration)

	return configuration, err
}

func CreateConfiguration() Configuration {
	emptyConfig := Configuration{}

	configurationJson, err := emptyConfig.JsonString()
	if err != nil {
		fmt.Println("Unable to export configuration")
	}

	writeErr := ioutil.WriteFile(DefaultConfigPath(), configurationJson, 0644)

	if writeErr != nil {
		fmt.Println("Unable to write configuration at the default path.")
	}

	return emptyConfig
}

func (c Configuration) JsonString() ([]byte, error) {
	json, err := json.MarshalIndent(c, " ", "   ")

	return json, err
}

func (c Configuration) RegisterFolder(name string, path string) (Configuration, error) {
	for f := range c.Folders {
		if c.Folders[f].Path == path {
			return Configuration{}, errors.New("Folder already registered")
		}
	}
	folder := Folder{Name: name, Path: path}
	c.Folders = append(c.Folders, folder)

	return c, nil
}

func (c Configuration) PersistConfiguration() {
	configurationJson, err := json.MarshalIndent(c, " ", "    ")

	writeErr := ioutil.WriteFile(DefaultConfigPath(), configurationJson, 0644)
	if err != nil || writeErr != nil {
		fmt.Println("Unable to persist config")
	}
}
