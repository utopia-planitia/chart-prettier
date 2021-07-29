package prettier

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
	Yaml string
}

type List struct {
	Items []interface{}
}

func SplitManifests(yml string) ([]Manifest, error) {
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

	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		extracted, err := manifestsFromChunk(chunk)
		if err != nil {
			return []Manifest{}, err
		}

		manifests = append(manifests, extracted...)
	}

	return manifests, nil
}

func manifestsFromChunk(chunk string) ([]Manifest, error) {
	manifest, err := NewManifest(chunk)
	if err == io.EOF {
		return []Manifest{}, nil
	}
	if err != nil {
		return []Manifest{}, fmt.Errorf("parse yaml failed: %v: %v", err, chunk)
	}

	if strings.ToLower(manifest.Kind) != "list" {
		return []Manifest{manifest}, nil
	}

	ls := List{}

	err = yaml.Unmarshal([]byte(chunk), &ls)
	if err != nil {
		return []Manifest{}, err
	}

	manifests := []Manifest{}

	for _, element := range ls.Items {
		var buf bytes.Buffer

		enc := yaml.NewEncoder(&buf)

		enc.SetIndent(2)

		err := enc.Encode(element)
		if err != nil {
			return []Manifest{}, err
		}

		manifest, err := NewManifest(buf.String())
		if err != nil {
			return []Manifest{}, err
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func NewManifest(yml string) (Manifest, error) {
	d := yaml.NewDecoder(strings.NewReader(excludeTemplates(yml)))
	m := Manifest{}

	err := d.Decode(&m) // invalid yaml returns empty string as a kind
	if err != nil {
		return Manifest{}, err
	}

	m.Yaml = strings.TrimSpace(yml)

	return m, nil
}

func excludeTemplates(ymlIn string) string {
	buf := strings.Builder{}

	for _, line := range strings.Split(ymlIn, "\n") {
		if strings.Contains(line, "{{") {
			continue
		}

		if strings.Contains(line, "}}") {
			continue
		}

		buf.WriteString(line + "\n")
	}

	return buf.String()
}
