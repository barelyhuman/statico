package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Site struct {
	Name        string `yaml:"name"`
	Link        string `yaml:"link"`
	Description string `yaml:"description"`
}

type Config struct {
	Site          Site `yaml:"site"`
	TemplateNames struct {
		PostTemplateName      string `yaml:"post_template"`
		PostIndexTemplateName string `yaml:"post_index_template"`
		PageTemplateName      string `yaml:"page_template"`
		RSSTemplateName       string `yaml:"rss_template"`
	} `yaml:"template_names"`

	ContentPath    string   `yaml:"content_path"`
	TemplatesPath  string   `yaml:"templates_path"`
	OutPath        string   `yaml:"out_path"`
	PublicFolder   string   `yaml:"public_folder"`
	PostIndexPath  string   `yaml:"post_index_path"`
	IndexedFolders []string `yaml:"indexed_folders"`
	GenerateRss    bool     `yaml:"generate_rss"`
	RssOutPath     string   `yaml:"rss_out_path"`
	Port           string   `yaml:"port"`
}

func (cfg *Config) readConfig(configFilePath string) error {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), &cfg)
	if err != nil {
		return err
	}

	return nil
}

func readConfigFlags(flags ...string) *Config {
	configFile := "./config.yml"

	for _, flag := range flags {
		if len(flag) > 0 {
			configFile = flag
		}
	}
	configRef := &Config{}
	err := configRef.readConfig(configFile)
	if err != nil {
		log.Fatal("Error reading config: ", err)
	}
	return configRef
}
