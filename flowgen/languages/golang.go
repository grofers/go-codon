package languages

import (
	flowgen_shared "github.com/grofers/go-codon/flowgen/shared"
	shared "github.com/grofers/go-codon/shared"
	"text/template"
	"os"
	"fmt"
	"path/filepath"
	// "strings"
	html_template "html/template"
	// "fmt"
)

type GoGenerator struct {
	Data                 *flowgen_shared.PostSpec
	Dest                 string
	Templates            string
	BaseImport           string
}

type OutputObj struct {
	Type                 string
	Children             map[string]OutputObj
	ExpressionSrno       int
	FlowName             string
}

func (g GoGenerator) getOutputObj(output interface{}) OutputObj {
	outputObj := OutputObj{}
	switch output_val := output.(type) {
	case map[interface{}]interface{}:
		outputObj.Type = "map"
		outputObj.Children = map[string]OutputObj{}
		for key, val := range output_val {
			outputObj.Children[key.(string)] = g.getOutputObj(val)
		}
	default:
		outputObj.Type = "value"
		output_val_str := fmt.Sprintf("%v", output_val)
		outputObj.ExpressionSrno = g.Data.ExpressionMap[output_val_str].Srno
	}
	outputObj.FlowName = g.Data.OrigSpec.Name
	return outputObj
}

func (g *GoGenerator) postProcess() error {
	g.Data.LanguageSpec["OutputObj"] = g.getOutputObj(g.Data.OrigSpec.Output)

	if _, ok := g.Data.OrigSpec.References["go"]; !ok {
		g.Data.OrigSpec.References["go"] = make(map[string]string)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	import_path, err := shared.BaseImport(pwd)
	if err != nil {
		return err
	}
	g.BaseImport = import_path

	// Add imports required for expressions
	for _, expr_obj := range g.Data.ExpressionMap {
		switch expr_obj.Type {
		case "json":
			g.Data.OrigSpec.References["go"]["json"] = "encoding/json"
		case "jmes":
			g.Data.OrigSpec.References["go"]["jmespath"] = "github.com/jmespath/go-jmespath"
		}
	}

	// Add imports required for actions
	for _, action_obj := range g.Data.ActionMap {
		action_obj_type := action_obj.Type
		if action_obj_type == "" {
			continue
		}
		// If import location is not custom defined we assume
		// it is present in base import path
		if _, ok := g.Data.OrigSpec.References["go"][action_obj_type]; !ok {
			g.Data.OrigSpec.References["go"][action_obj_type] = g.BaseImport + "/" + action_obj_type
		}
	}

	return nil
}

func (g *GoGenerator) Generate() error {
	err := g.postProcess()
	if err != nil {
		return err
	}

	template_location := filepath.Join(g.Templates, "golang.gotmpl")

	tmpl, err := template.New("golang.gotmpl").Funcs(template.FuncMap{
		"escapestring": html_template.JSEscapeString,
		"pascalize": flowgen_shared.Pascalize,
	}).ParseFiles(template_location)
	if err != nil {
		return err
	}

	dest_base, _ := filepath.Split(g.Dest)

	err = os.MkdirAll(dest_base, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(g.Dest)
	defer f.Close()
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, g.Data)
	if err != nil {
		return err
	}

	return nil
}
