package cli

import (
	"fmt"
	"strings"

	"path"

	"io/ioutil"

	"sort"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func GenerateVerbCommands(c *Config, envVars []*string) []cli.Command {
	var verbCommands []cli.Command
	for name, verb := range c.Verbs {
		verb.Name = name
		verbCommands = append(verbCommands, GenerateVerbCommand(verb, c, envVars))
	}

	sort.Sort(cli.CommandsByName(verbCommands))

	return verbCommands
}

func GenerateVerbCommand(verb *Verb, c *Config, envVars []*string) cli.Command {
	return cli.Command{
		Name:         verb.Name,
		Description:  verb.Description,
		Usage:        verb.Usage,
		Category:     verb.Category,
		HideHelp:     true,
		BashComplete: verbCompletions(c),
		Action: func(ctx *cli.Context) error {
			if verb.Args.Min != nil {
				if ctx.NArg() < *verb.Args.Min {
					cli.ShowCommandHelpAndExit(ctx, ctx.Command.Name, 0)
				}
			}

			if verb.Args.Max != nil {
				if ctx.NArg() > *verb.Args.Max {
					cli.ShowCommandHelpAndExit(ctx, ctx.Command.Name, 0)
				}
			}

			context := strings.Split(ctx.Command.FullName(), " ")
			for key, value := range c.Contexts[context[0]].Env {
				envVar := fmt.Sprintf("%s=%s", key, *value)
				envVars = append(envVars, &envVar)
			}

			if c.Contexts[context[0]].EnvFlags != nil {
				for flag, envMap := range c.Contexts[context[0]].EnvFlags {
					if ctx.GlobalBool(flag) {
						for key, value := range *envMap {
							envVar := fmt.Sprintf("%s=%s", key, *value)
							envVars = append(envVars, &envVar)
						}
					}
				}
			}

			for _, command := range verb.Commands {
				if err := ExecuteComposeCommand(c.Home, envVars, command, ctx.Args()); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func verbCompletions(c *Config) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		var completions []string
		var err error

		context := strings.Split(ctx.Command.FullName(), " ")
		if context[0] == "global" {
			var files []string
			if context[0] == "global" {
				for _, context := range c.Contexts {
					if context.Name != "global" && context.Name != "group" {
						files = append(files, path.Join(c.Home, fmt.Sprintf("%s.yaml", context.Name)))
					}
				}
			}

			for _, file := range files {
				completions, err = appendCompletions(file, completions)
				if err != nil {
					panic(err)
				}
			}
		} else {
			file := path.Join(c.Home, fmt.Sprintf("%s.yaml", context[0]))
			completions, err = appendCompletions(file, completions)
			if err != nil {
				panic(err)
			}
		}

		fmt.Fprintf(ctx.App.Writer, strings.Join(completions, " "))
	}
}

func appendCompletions(file string, completions []string) ([]string, error) {
	composeFile := map[interface{}]interface{}{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bytes, composeFile); err != nil {
		return nil, err
	}

	for service := range composeFile["services"].(map[interface{}]interface{}) {
		completions = append(completions, service.(string))
	}

	return completions, nil
}
