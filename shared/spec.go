package shared

import (
	"log"
	"strings"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)


const (
	SWAGGER = iota
	UNKNOWN = -1
)


func DetectFileSpec(path string) int {
	_, filename := filepath.Split(path)
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		yamlFile, err := ioutil.ReadFile(path)
		if err != nil {
	        log.Println(err)
	        return UNKNOWN
	    }
	    c := map[string]interface{}{}
	    err = yaml.Unmarshal(yamlFile, &c)
	    if err != nil {
	    	log.Println(err)
	    	return UNKNOWN
	    }
	    if _, ok := c["swagger"]; ok {
	    	return SWAGGER
	    }
	    return UNKNOWN
	} else {
		return UNKNOWN
	}
}
