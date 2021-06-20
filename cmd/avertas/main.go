package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"

	conf "github.com/AxonC/avertas/pkg/configuration"
)

func GetKeysFromMap(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func RegisterDirectoryHandler(c *cli.Context) error {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot determine current directory")
	}
	configuration, err := conf.ReadConfiguration()
	if err != nil {
		conf.CreateConfiguration()
		fmt.Println("Created new configuration.")
	}
	var folderName string
	path := strings.Split(currentPath, "/")
	// deduce the folder name from the last part of the split string.
	folderName = path[len(path)-1]

	newConfig, err := configuration.RegisterFolder(folderName, currentPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// save the configuration with the new directory
	newConfig.PersistConfiguration()
	fmt.Println("Successfully added folder.")

	return nil
}

func ListProjects(c *cli.Context) error {
	configuration, err := conf.ReadConfiguration()

	if err != nil {
		return nil
	}

	var registeredTypes = make(map[string][]string)
	for _, s := range configuration.Folders {
		files, err := ioutil.ReadDir(s.Path)

		if err != nil {
			log.Fatal(err)
		}

		var folders []string
		for _, f := range files {
			// filter out system directories e.g. .DS_Store etc and only push directories
			if f.Name()[0] != '.' && f.IsDir() {
				folders = append(folders, f.Name())
			}
		}
		registeredTypes[s.Name] = folders
	}

	prompt := promptui.Select{
		Label: "Select Project Type",
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

// Show the location of the configuration file to the user.
func ShowConfigPathHandler(c *cli.Context) error {
	fmt.Println(conf.DefaultConfigPath())
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
			{
				Name:   "register",
				Usage:  "Register current directory containing projects",
				Action: RegisterDirectoryHandler,
			},
			{
				Name:   "config",
				Usage:  "Show path of config",
				Action: ShowConfigPathHandler,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
