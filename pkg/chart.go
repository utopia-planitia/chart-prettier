package prettier

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type Chart struct {
	manifests []Manifest
}

func (c *Chart) LoadChart(appFs afero.Fs, path string) error {
	addManifests := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isManifestFile(info.Name()) {
			content, err := afero.ReadFile(appFs, path)
			if err != nil {
				return err
			}

			c.AddManifests(string(content))
		}

		return nil
	}

	err := afero.Walk(appFs, path, addManifests)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chart) DeleteFiles(appFs afero.Fs, path string) error {
	deleteManifests := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isManifestFile(info.Name()) {
			appFs.Remove(path)
		}

		return nil
	}

	err := afero.Walk(appFs, path, deleteManifests)
	if err != nil {
		return err
	}

	return nil
}

func isManifestFile(path string) bool {
	for _, ext := range []string{".yaml", ".yml", ".yaml.tpl", ".yml.tpl"} {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}

func (c *Chart) WriteOut(appFs afero.Fs, path string) error {
	unique, err := uniqueNames(c.manifests)
	if err != nil {
		return err
	}

	for name, manifest := range unique {
		ext := ".yaml"
		if strings.Contains(manifest.Yaml, "{{") {
			ext = ".yaml.tpl"
		}

		filename := filepath.Join(path, name+ext)
		content := manifest.Yaml

		err := afero.WriteFile(appFs, filename, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chart) AddManifests(yaml string) error {
	manifests, err := SplitManifests(yaml)
	if err != nil {
		return err
	}

	c.manifests = append(c.manifests, manifests...)

	return nil
}

func uniqueNames(manifests []Manifest) (map[string]Manifest, error) {
	kindAsName := func(m Manifest) string { return m.Kind }
	nameInName := func(m Manifest) string { return m.Kind + "-" + m.Metadata.Name }
	namespaceInName := func(m Manifest) string { return m.Kind + "-" + m.Metadata.Name + "-" + m.Metadata.Namespace }

	kindNamed, collisions := filterUniqueNames(manifests, kindAsName)
	nameNamed, collisions := filterUniqueNames(collisions, nameInName)
	namespaceNamed, collisions := filterUniqueNames(collisions, namespaceInName)

	if len(collisions) > 0 {
		return map[string]Manifest{}, fmt.Errorf("not all manifests are unique using kind, name, and namespace as identifier")
	}

	all := map[string]Manifest{}
	for name, manifest := range kindNamed {
		all[strings.ToLower(name)] = manifest
	}
	for name, manifest := range nameNamed {
		all[strings.ToLower(name)] = manifest
	}
	for name, manifest := range namespaceNamed {
		all[strings.ToLower(name)] = manifest
	}

	return all, nil
}

func filterUniqueNames(manifests []Manifest, name func(Manifest) string) (map[string]Manifest, []Manifest) {
	uniqueNamed := map[string]Manifest{}

	named := map[string][]Manifest{}

	for _, manifest := range manifests {
		name := name(manifest)
		named[name] = append(named[name], manifest)
	}

	collisions := []Manifest{}

	for name, manifests := range named {
		if len(manifests) > 1 {
			collisions = append(collisions, manifests...)
			continue
		}

		uniqueNamed[name] = manifests[0]
	}

	return uniqueNamed, collisions
}
