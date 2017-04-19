package golang

import (
	"os"
	"log"
	"strings"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/go-openapi/swag"

	codon_shared "github.com/grofers/go-codon/shared"
	config "github.com/grofers/go-codon/runtime/config"
)

// Important: Dont use lists, use maps
// While go randomizes iteration results, template does not.
// So we let template ensure idiomatic (and sorted) results.
type generator struct {
	CurrentSpecFile string
	CurrentSpecFilePath string
	CurrentAPIName string
	CurrentDirTarget string
	CurrentDirPath string
	CurrentDirName string
	ProjectName string
	ClientImports map[string]string
	ClientEndpoints map[string]map[string]string
	ClientsUsed map[string]bool
	ClientEndpointGoNames map[string]string
}

var upstream_generator_list = map[int]func(*generator)bool {
	codon_shared.SWAGGER: GenerateUpstreamSwagger,
	codon_shared.UNKNOWN: GenerateUnknown,
}

func (gen *generator) Init() {
	gen.ClientImports = make(map[string]string)
	gen.ClientEndpoints = make(map[string]map[string]string)
	gen.ClientsUsed = make(map[string]bool)
	gen.ClientEndpointGoNames = make(map[string]string)
}

func (gen *generator) UpdateCurrentDirPath() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	gen.CurrentDirPath = pwd
	log.Println("Updated working directory:", pwd)
	_, gen.CurrentDirName = filepath.Split(pwd)
	gen.ProjectName = gen.CurrentDirName
	log.Println("Working with project name:", gen.ProjectName)
	return nil
}

func (gen *generator) process_templates() error {
	for _, asset := range AssetNames() {
		t := template.New(asset)

		// If the content being templated is a template itself
		if strings.HasSuffix(asset, ".gotmpl") {
			t = t.Delims("{|{", "}|}")
		}

		t, err := t.Parse(string(MustAsset(asset)))
		if err != nil {
			log.Println("Failed to get asset:", err)
			return err
		}

		var new_asset_path string
		if strings.HasSuffix(asset, ".gofile") {
			new_asset_path = filepath.Join(gen.CurrentDirPath, strings.TrimSuffix(asset, ".gofile") + ".go")
		} else {
			new_asset_path = filepath.Join(gen.CurrentDirPath, asset)
		}
		base_path, _ := filepath.Split(new_asset_path)

		err = os.MkdirAll(base_path, 0755)
		if err != nil {
			log.Println("Failed to create file:", err)
			return err
		}

		f, err := os.Create(new_asset_path)
		defer f.Close()
		if err != nil {
			log.Println("Failed to create file:", err)
			return err
		}

		err = t.Execute(f, gen)
		if err != nil {
			log.Println("Failed to execute template:", err)
			return err
		}
	}
	return nil
}

func (gen *generator) process_config() error {
	gen.ClientEndpoints = config.YmlConfig.Endpoints

	for endpoint_name, endpoint := range gen.ClientEndpoints {
		client, ok := endpoint["client"]
		if !ok {
			continue
		}
		gen.ClientsUsed[client] = true
		gen.ClientEndpointGoNames[endpoint_name] = swag.ToGoName(endpoint_name)
	}

	return nil
}

func GenerateUnknown(gen *generator) bool {
	log.Println("Ignoring unknown file format for:", gen.CurrentSpecFile)
	return true
}

func (gen *generator) GenerateUpstream() bool {
	// Get list of all the files in spec/clients
	files, err := ioutil.ReadDir("spec/clients")
	if err != nil {
		log.Println(err)
		return false
	}

	for _, file := range files {
		gen.CurrentSpecFile = file.Name()
		gen.CurrentSpecFilePath = filepath.Join("spec/clients", file.Name())
		log.Println("Processing upstream spec: ", file.Name())
		if file.IsDir() {
			log.Println(file.Name(), "is a directory. Ignoring.")
			continue
		}
		spec_type := codon_shared.DetectFileSpec(gen.CurrentSpecFilePath)
		gen_func := upstream_generator_list[spec_type]
		if ok := gen_func(gen); !ok {
			log.Println("Failed to generate code for spec", file.Name())
			return false
		}
	}
	return true
}

func (gen *generator) GenerateContent() bool {
	if err := gen.process_config(); err != nil {
		log.Println(err)
		return false
	}
	if err := gen.process_templates(); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (gen *generator) Generate() bool {
	gen.Init()
	log.Println("Generating a codon project in golang ...")

	if err := gen.UpdateCurrentDirPath(); err != nil {
		return false
	}

	if !gen.GenerateUpstream() {
		return false
	}

	if !gen.GenerateContent() {
		return false
	}

	return true
}

var Generator = generator{}
