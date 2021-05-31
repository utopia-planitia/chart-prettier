package prettier

import (
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

	for _, doc := range chunks {
		doc = strings.TrimSpace(doc)

		if doc == "" {
			continue
		}

		manifest, err := NewManifest(doc)
		if err == io.EOF {
			continue
		}
		if err != nil {
			return []Manifest{}, fmt.Errorf("parse yaml failed: %v: %v", err, doc)
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func NewManifest(yml string) (Manifest, error) {
	d := yaml.NewDecoder(strings.NewReader(stripDown(yml)))
	d.KnownFields(true)
	m := Manifest{}

	err := d.Decode(&m) // invalid yaml returns empty string as a kind
	if err != nil {
		return Manifest{}, err
	}

	m.Yaml = strings.TrimSpace(yml)

	return m, nil
}

func stripDown(ymlIn string) string {
	keepLine := func(line string) bool {
		lowerCase := strings.ToLower(line)
		return strings.HasPrefix(lowerCase, "kind: ") || strings.HasPrefix(lowerCase, "metadata:") || strings.HasPrefix(lowerCase, "  name:") || strings.HasPrefix(lowerCase, "  namespace:")
	}

	buf := strings.Builder{}

	for _, line := range strings.Split(ymlIn, "\n") {
		if keepLine(line) {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
