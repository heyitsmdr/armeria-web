package armeria

import (
	"io/ioutil"
	"log"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type config struct {
	HTTPPort   int    `yaml:"httpPort"`
	PublicPath string `yaml:"publicPath"`
	Production bool   `yaml:"production"`
	DataPath   string `yaml:"dataPath"`
}

func parseConfigFile(filePath string) config {
	data := readConfigFile(filePath)
	c := unmarshalConfig(data)
	return c

}

func readConfigFile(filePath string) []byte {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	return data
}

func unmarshalConfig(data []byte) config {
	c := config{}
	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		Armeria.log.Debug("Unmarshaling error",
			zap.Error(err),
		)
	}
	return c
}
