package templates

import (
	"github.com/flosch/pongo2"
	"errors"
	"io"
	"time"
	// "fmt"
)

var TemplateMap = map[string](*pongo2.Template) {}

var init_completed = false

func Init(assetNames func()[]string, get_asset func(name string)([]byte, error)) error {
	if init_completed {
		return nil
	}
	for _, asset_name := range assetNames() {
		asset_b, err := get_asset(asset_name)
		if err == nil {
			tmpl, err2 := pongo2.FromString(string(asset_b))
			if err2 == nil {
				TemplateMap[asset_name] = tmpl
			}
		}
	}
	if err3 := initVars(); err3 != nil {
		return err3
	}
	init_completed = true
	return nil
}

func initVars() error {
	pongo2.Globals["_go"] = map[string]interface{} {
		"time": map[string] interface{} {
			"Unix": time.Unix,
		},
	}
	pongo2.Globals["_pongo"] = map[string]interface{} {
		"AsValue": pongo2.AsValue,
	}
	return nil
}

func Execute(templatePath string, context pongo2.Context, writer io.Writer) error {
	tmpl, ok := TemplateMap[templatePath]
	if ok == false {
		return errors.New("Template not found")
	}
	return tmpl.ExecuteWriter(context, writer)
}
