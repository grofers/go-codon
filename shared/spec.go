package shared

import (
	"os"
	"log"
	"errors"
	"strings"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
	goruntime "runtime"
	"github.com/go-openapi/swag"
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


// Copyright 2015 go-swagger maintainers
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
func BaseImport(tgt string) (string, error) {
	// Modified error reporting structure to match go-codon's
	p, err := filepath.Abs(tgt)
	if err != nil {
		return "", err
	}

	var pth string
	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		pp := filepath.Join(filepath.Clean(gp), "src")
		var np, npp string
		if goruntime.GOOS == "windows" {
			np = strings.ToLower(p)
			npp = strings.ToLower(pp)
		}
		if strings.HasPrefix(np, npp) {
			pth, err = filepath.Rel(pp, p)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if pth == "" {
		return "", errors.New("target must reside inside a location in the $GOPATH/src")
	}
	return pth, nil
}

func Pascalize(arg string) string {
	if len(arg) == 0 || arg[0] > '9' {
		return swag.ToGoName(arg)
	}
	if arg[0] == '+' {
		return swag.ToGoName("Plus " + arg[1:])
	}
	if arg[0] == '-' {
		return swag.ToGoName("Minus " + arg[1:])
	}

	return swag.ToGoName("Nr " + arg)
}
