package prettier

import (
	"io"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
	Yaml string
}

func SplitManifests(yml string) []Manifest {
	yml = strings.TrimSpace(yml)

	if strings.HasPrefix(yml, "---\n") {
		yml = "\n" + yml
	}

	if strings.HasSuffix(yml, "\n---") {
		yml += "\n"
	}

	separator := regexp.MustCompile("\n---\n")
	chunks := separator.Split(yml, -1)
	manifests := []Manifest{}

	for _, doc := range chunks {
		doc = strings.TrimSpace(doc)

		if doc == "" {
			continue
		}

		manifest, err := parseManifest(doc)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Panicf("parse yaml failed: %v: %v", err, doc)
		}

		manifests = append(manifests, manifest)
	}

	return manifests
}

func parseManifest(yml string) (Manifest, error) {
	d := yaml.NewDecoder(strings.NewReader(yml))
	m := Manifest{}

	err := d.Decode(&m) // invalid yaml returns empty string as a kind
	if err != nil {
		return Manifest{}, err
	}

	m.Yaml = yml

	return m, nil
}
