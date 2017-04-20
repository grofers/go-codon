package golang

import (
	"os"
	"path/filepath"
	"log"
	"text/template"
	"strings"
)

// var templates = template.Must(template.ParseGlob("*"))

type bootstrapper struct {
	CurrentDirPath string
	CurrentDirName string
	ProjectName string
}

func (bs *bootstrapper) UpdateCurrentDirPath() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	bs.CurrentDirPath = pwd
	log.Println("Updated working directory:", pwd)
	_, bs.CurrentDirName = filepath.Split(pwd)
	bs.ProjectName = bs.CurrentDirName
	log.Println("Working with project name:", bs.ProjectName)
	return nil
}

func (bs *bootstrapper) process_templates() error {
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
			new_asset_path = filepath.Join(bs.CurrentDirPath, strings.TrimSuffix(asset, ".gofile") + ".go")
		} else {
			new_asset_path = filepath.Join(bs.CurrentDirPath, asset)
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

		err = t.Execute(f, bs)
		if err != nil {
			log.Println("Failed to execute template:", err)
			return err
		}
	}
	return nil
}

func (bs *bootstrapper) Bootstrap() bool {
	log.Println("Bootstrapping a codon project in golang ...")

	if err := bs.UpdateCurrentDirPath(); err != nil {
		return false
	}

	if err := bs.process_templates(); err != nil {
		return false
	}

	return true
}

var Bootstrapper = bootstrapper{}
