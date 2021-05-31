package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	prettier "github.com/utopia-planitia/chart-prettier/pkg"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	yml, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	manifests, err := prettier.SplitManifests(string(yml))
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		_, err = fmt.Printf("---\n%s\n", manifest.Yaml)
		if err != nil {
			return err
		}
	}

	return nil
}
