package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"errors"
	"strings"
)

type Config struct {
	Endpoints	map[string]map[string]string	`yaml:"endpoints"`
	Constants	map[string]interface{}			`yaml:"constants"`
}

func (cfg *Config) GetEndpoint(endpoint string) *map[string]string {
	if endpoint_cfg, ok := cfg.Endpoints[endpoint]; ok {
		return &endpoint_cfg
	} else {
		return nil
	}
}

func (cfg *Config) GetConstant(name string) (interface{}, error) {
	if val, ok := cfg.Constants[name]; ok {
		return val, nil
	} else {
		return nil, errors.New("Constant not found")
	}
}

func (cfg *Config) GetConstantPath(path string) (interface{}, error) {
	path_s := strings.Split(path, ".")
	obj, err := cfg.GetConstant(path_s[0])
	if err != nil {
		return nil, err
	}
	for i:=1; i<len(path_s); i++ {
		obj_t, ok := obj.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("Constant not found")
		}
		obj = obj_t[path_s[i]]
	}
	return obj, nil
}

func (cfg *Config) MustGetConstant(name string) interface{} {
	val, err := cfg.GetConstant(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (cfg *Config) MustGetConstantPath(path string) interface{} {
	val, err := cfg.GetConstantPath(path)
	if err != nil {
		panic(err)
	}
	return val
}

var YmlConfig = ReadYmlConfig()

func ReadYmlConfig() *Config {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil
	}

	var t Config

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return nil
	}

	t.Constants = cleanupStringInterfaceMap(t.Constants)

	return &t
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	case map[string]interface{}:
		return cleanupStringInterfaceMap(v)
	default:
		return v
	}
}

func cleanupStringInterfaceMap(in map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[k] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func GetEndpoint(endpoint string) *map[string]string {
	return YmlConfig.GetEndpoint(endpoint)
}

func GetConstant(name string) (interface{}, error) {
	return YmlConfig.GetConstant(name)
}

func MustGetConstant(name string) interface{} {
	return YmlConfig.MustGetConstant(name)
}

func MustGetConstantPath(path string) interface{} {
	return YmlConfig.MustGetConstantPath(path)
}
