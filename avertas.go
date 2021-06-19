package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type Configuration struct {
	Folders []Folder `json:"folders"`
}

type Folder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func GetKeysFromMap(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func ReadConfiguration(filePath string) (Configuration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to find configuration")
	}
	defer file.Close()

	contents, _ := ioutil.ReadAll(file)

	var configuration Configuration
	json.Unmarshal([]byte(contents), &configuration)

	return configuration, err
}

func ListProjects(c *cli.Context) error {
	configuration, err := ReadConfiguration("config.json")

	if err != nil {
		return nil
	}

	var registeredTypes = make(map[string][]string)
	for _, s := range configuration.Folders {
		homeDirectory, _ := os.UserHomeDir()
		files, err := ioutil.ReadDir(homeDirectory + s.Path)

		if err != nil {
			log.Fatal(err)
		}

		var folders []string
		for _, f := range files {
			// filter out system directories
			if f.Name()[0] != '.' {
				folders = append(folders, f.Name())
			}
		}
		fmt.Println(folders)
		registeredTypes[s.Name] = folders
	}
	fmt.Println(registeredTypes)

	prompt := promptui.Select{
		Label: "Select Project Testing",
		Items: GetKeysFromMap(registeredTypes),
	}

	_, selected, err := prompt.Run()

	if err != nil {
		fmt.Println("Unknown Option")
	}

	projectPrompt := promptui.Select{
		Label: "Select Folder",
		Items: registeredTypes[selected],
	}

	_, selectedProject, err := projectPrompt.Run()

	if err != nil {
		fmt.Println("Unknown option selected.")
		return err
	}
	fmt.Println(selectedProject)
	return nil
}

func main() {
	app := &cli.App{
		Name:  "avertas",
		Usage: "Project switcher",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "List active projects",
				Action: ListProjects,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
