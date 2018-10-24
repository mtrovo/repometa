package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type RepoMeta struct {
	Name          string
	IssueTracking string `yaml:"issue-tracking"`
	Criticality   string
	ArchDomain    string `yaml:"arch-domain"`
	SlackChannel  string `yaml:"slack-channel"`
	Email         string
	BasePath      string `yaml:"basepath"`
}

var maintainersTemplate = `Slack: %v
%v
`
var contributingTemplate string

func init() {
	rootCmd.AddCommand(generateCmd)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	templateFile := path.Join(path.Dir(filename), "contributing.template")
	contentBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}

	contributingTemplate = string(contentBytes)
}
func generateMaintainers(meta *RepoMeta) error {

	content := fmt.Sprintf(maintainersTemplate, meta.SlackChannel, meta.Email)
	return ioutil.WriteFile("MAINTAINERS", []byte(content), 0644)
}

func generateContributing(meta *RepoMeta) error {
	template, err := template.New("CONTRIBUTING").Parse(contributingTemplate)
	if err != nil {
		return err
	}
	out, err := os.Create("CONTRIBUTING.md")
	defer out.Close()
	if err != nil {
		return err
	}

	template.Execute(out, meta)
	return nil
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate all files related to repometa",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return fmt.Errorf("Could not read file: %s", err)
		}

		meta := RepoMeta{BasePath: "/"}
		err = yaml.Unmarshal(data, &meta)
		if err != nil {
			return fmt.Errorf("Could not parse file: %s", err)
		}

		err = generateMaintainers(&meta)
		if err != nil {
			return err
		}
		err = generateContributing(&meta)
		if err != nil {
			return err
		}
		return nil
	},
}
