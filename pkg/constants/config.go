package constants

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type SystemConfig struct {
	Agent struct {
		Config struct {
			WAN_LINK string `yaml:"wan_link"`
			Mode     string `yaml:"mode"`
		} `yaml:"config"`
	} `yaml:"mrz_mpmgmt_agent"`
}

var SystemConfigEnv SystemConfig

func ParseSystemConfig() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	err = yaml.Unmarshal(yamlFile, &SystemConfigEnv)
	if err != nil {
		panic(err)
	}

}
