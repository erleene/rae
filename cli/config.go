package cli

import (
	"io/ioutil"

	"path"

	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Home string

	Contexts   map[string]*Context `yaml:"contexts"`
	Env        map[string]*string  `yaml:"env"`
	EnvFiles   []string            `yaml:"env_files"`
	GroupVerbs []string            `yaml:"group_verbs"`
	Groups     map[string]*Group   `yaml:"groups"`
	Verbs      map[string]*Verb    `yaml:"verbs"`
	Recipes    map[string]*Recipe  `yaml:"recipes"`
}

type Context struct {
	Name string

	Aliases     []string                       `yaml:"aliases"`
	Category    string                         `yaml:"category"`
	Description string                         `yaml:"description"`
	Env         map[string]*string             `yaml:"env"`
	EnvFlags    map[string]*map[string]*string `yaml:"env_flags"`
	Usage       string                         `yaml:"usage"`
}

type Group struct {
	Name string

	Aliases     []string              `yaml:"aliases"`
	Category    string                `yaml:"category"`
	Description string                `yaml:"description"`
	Members     []map[string][]string `yaml:"members"`
	Usage       string                `yaml:"usage"`
}

type Verb struct {
	Name string

	Args struct {
		Min *int `yaml:"min"`
		Max *int `yaml:"max"`
	} `yaml:"args"`
	Category    string     `yaml:"category"`
	Commands    [][]string `yaml:"commands"`
	Description string     `yaml:"description"`
	Usage       string     `yaml:"usage"`
}

type Recipe struct {
	Name string

	Usage        string          `yaml:"usage"`
	Description  string          `yaml:"description"`
	Instructions []*Instructions `yaml:"instructions"`
}

type Instructions struct {
	Group       string   `yaml:"group"`
	Verb        string   `yaml:"verb"`
	Healthcheck []string `yaml:"healthcheck"`
}

func (c *Config) Load(home string) error {
	bytes, err := ioutil.ReadFile(path.Join(home, "rae.yaml"))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return err
	}

	c.Home = home

	return nil
}

func (c *Config) LoadOverride(home string) error {
	override := path.Join(home, "rae.override.yaml")

	if _, err := os.Stat(override); os.IsNotExist(err) {
		return nil
	}

	bytes, err := ioutil.ReadFile(override)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return err
	}

	c.Home = home

	return nil
}

func (c *Config) MergeOverride(o *Config) error {
	if o != nil {
		if err := mergo.Merge(c, o, mergo.WithOverride); err != nil {
			return err
		}
	}

	return nil
}
