package golang

import (
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"path/filepath"

	swagger_generator "github.com/go-swagger/go-swagger/generator"
	codon_shared "github.com/grofers/go-codon/shared"
)

type generator struct {
	CurrentSpecFile string
	CurrentSpecFilePath string
	CurrentAPIName string
}

var upstream_generator_list = map[int]func(*generator)bool {
	codon_shared.SWAGGER: GenerateUpstreamSwagger,
	codon_shared.UNKNOWN: GenerateUnknown,
}

func GenerateUpstreamSwagger(gen *generator) bool {
	gen.CurrentAPIName = strings.TrimSuffix(gen.CurrentSpecFile, ".yml")
	gen.CurrentAPIName = strings.TrimSuffix(gen.CurrentAPIName, ".yaml")

	// -t clients/$api_name/ -T client-templates/
	opts := &swagger_generator.GenOpts{
		APIPackage:        "operations",
		ModelPackage:      "models",
		ServerPackage:     "restapi",
		ClientPackage:     "client",
		Principal:         "",
		DefaultScheme:     "http",
		DefaultProduces:   "application/json",
		IncludeModel:      true,
		IncludeValidator:  true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		ValidateSpec:      true,
		Tags:              []string{},
		IncludeSupport:    true,
		DumpData:          false,
		Spec: gen.CurrentSpecFilePath,
		Target:            fmt.Sprintf("clients/%s/", gen.CurrentAPIName),
		TemplateDir:       "spec/templates/",
	}
	if err := opts.EnsureDefaults(true); err != nil {
		log.Println(err)
		return false
	}
	if err := swagger_generator.GenerateClient("", []string{}, []string{}, opts); err != nil {
		log.Println(err)
		return false
	}
	return true
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

func (gen *generator) Generate() bool {
	log.Println("Generating a codon project in golang ...")

	if !gen.GenerateUpstream() {
		return false
	}

	return true
}

var Generator = generator{}
