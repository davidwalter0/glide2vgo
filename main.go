package main // import "github.com/davidwalter0/glide2vgo"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/davidwalter0/go-cfg"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Overwrite bool `json:"overwrite" doc:"replace an existing file with new output"`
	VGoDivert bool `json:"vgo-divert" doc:"use default name go.mod+base package name: go.mod.glide2vgo"`
}

var app = &Config{}

func init() {
	var err error
	if err = cfg.Process("", app); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

const VGoMod = "go.mod"

type Dependency struct {
	Package string
	Version string
}

type GlideConfig struct {
	Package string       `json:"package"`
	Import  []Dependency `json:"import"`
}

// GlideConfigFromYAML returns an instance of GlideConfig from YAML
func GlideConfigFromYAML(yml []byte) (*GlideConfig, error) {
	cfg := &GlideConfig{}
	err := yaml.Unmarshal([]byte(yml), &cfg)
	return cfg, err
}

func main() {
	var goModText = ""
	var err error
	var b []byte
	var c *GlideConfig
	if b, err = ioutil.ReadFile("glide.yaml"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if c, err = GlideConfigFromYAML(b); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	goModText += fmt.Sprintf("module %s\n", c.Package)
	goModText += fmt.Sprintf("require (\n")
	var dep Dependency
	for _, dep = range c.Import {
		goModText += fmt.Sprintf("\t%s %s\n", dep.Package, dep.Version)
	}
	goModText += fmt.Sprintf(")\n\n")

	var packageName = filepath.Base(c.Package)
	var VGoModDivert = fmt.Sprintf("go.mod.%s", packageName)

	if app.VGoDivert {
		if err = ioutil.WriteFile(VGoModDivert, []byte(goModText), 0666); err != nil {
			log.Fatal(err)
		}
	} else if _, err := os.Stat(VGoMod); os.IsNotExist(err) || app.Overwrite {
		// If go.mod doesn't exist or overwrite, ok to write it
		if err = ioutil.WriteFile(VGoMod, []byte(goModText), 0666); err != nil {
			log.Fatal(err)
		}
	}
}
