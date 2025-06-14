package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/liujed/caddy-dns01proxy/jsonutil"
	"github.com/liujed/dns01proxy/internal/gen"
)

var topDir = filepath.Join("..", "..", "..")

//go:generate go run .
func main() {
	err := mainWithErr()
	if err != nil {
		panic(err)
	}
}

func mainWithErr() error {
	templateData := MakeTemplateData()

	buildsDirPath := filepath.Join(topDir, "internal", "builds")

	entries, err := os.ReadDir(buildsDirPath)
	if err != nil {
		return fmt.Errorf("unable to list plugins: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		metadataFilePath := filepath.Join(
			buildsDirPath, entry.Name(), "metadata.json",
		)
		buildMeta, err := jsonutil.UnmarshalFromFile[gen.Build](metadataFilePath)
		if err != nil {
			return fmt.Errorf(
				"unable to read build metadata from %q: %w",
				metadataFilePath,
				err,
			)
		}
		templateData.Builds[entry.Name()] = buildMeta
	}

	// Write goreleaser.yml.
	f, err := os.Create(filepath.Join(topDir, ".goreleaser.yml"))
	if err != nil {
		return fmt.Errorf("unable to create goreleaser.yml: %w", err)
	}
	func() {
		defer f.Close()
		err = generateGoreleaser(f, templateData)
	}()
	if err != nil {
		return fmt.Errorf("unable to create goreleaser.yml: %w", err)
	}

	return nil
}

func generateGoreleaser(w io.Writer, templateData TemplateData) error {
	_, err := io.WriteString(w, `# DO NOT EDIT. Generated by internal/gen/release/generator.go.

`)
	if err != nil {
		return fmt.Errorf("unable to write to output buffer: %w", err)
	}

	err = goreleaserTemplate.Execute(w, templateData)
	if err != nil {
		return fmt.Errorf("unable to execute goreleaser template: %w", err)
	}

	return nil
}
