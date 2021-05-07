package configreader

import (
	"io/ioutil"

	"github.com/joeshaw/envdecode"
	"gopkg.in/yaml.v2"
)

func ReadYamlAndSetEnv(path string, v interface{}) (err error) {
	if err = ReadYaml(path, v); err != nil {
		return err
	}

	return envdecode.Decode(v)
}

func ReadYaml(path string, v interface{}) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, v)
}
