package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/v72/github"
	"github.com/liujed/dns01proxy/internal/gen"
)

var topDir = filepath.Join("..", "..", "..")

//go:generate go run .
func main() {
	err := mainWithErr()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func mainWithErr() error {
	gh := github.NewClient(nil)
	ctx := context.Background()

	// Get a list of repositories from the caddy-dns GitHub organization.
	fmt.Println("Finding caddy-dns repositories...")
	repos, err := getCaddyDNSRepos(ctx, gh)
	if err != nil {
		return fmt.Errorf("unable to get caddy-dns repositories: %w", err)
	}

	fmt.Printf("Downloading %d caddy-dns modules...\n", len(repos))
	err = goGetModules(repos)
	if err != nil {
		return fmt.Errorf("unable to 'go get' caddy-dns modules: %w", err)
	}

	fmt.Println("Generating builds...")
	maxLen := len(strconv.Itoa(len(repos)))
	for idx, repo := range repos {
		modPath := "github.com/caddy-dns/" + repo.GetName()

		// Print progress info.
		fmt.Printf("(%*d/%d) %s\n", maxLen, idx+1, len(repos), modPath)

		version, hash, err := getGoModuleVersion(modPath)
		if err != nil {
			return fmt.Errorf(
				"unable to get Go module version/hash for %q: %w",
				modPath,
				err,
			)
		}

		templateData := TemplateData{
			PackageName: repo.GetName(),
			ModPath:     modPath,
			Version:     fmt.Sprintf("%s (%s)", version, hash),
		}

		// Ensure the module's build directory exists.
		buildDir := filepath.Join(topDir, "internal", "builds", repo.GetName())
		err = os.MkdirAll(buildDir, 0755)
		if err != nil {
			return fmt.Errorf(
				"unable to create build directory for %q: %w",
				modPath,
				err,
			)
		}

		// Write the Makefile.
		f, err := os.Create(filepath.Join(buildDir, "Makefile"))
		if err != nil {
			return fmt.Errorf("unable to create Makefile for %q: %w", modPath, err)
		}
		func() {
			defer f.Close()
			err = generateMakefile(f, templateData)
		}()
		if err != nil {
			return fmt.Errorf("unable to write Makefile for %q: %w", modPath, err)
		}

		// Write the build's source.
		f, err = os.Create(filepath.Join(buildDir, "main.go"))
		if err != nil {
			return fmt.Errorf("unable to create main.go for %q: %w", modPath, err)
		}
		func() {
			defer f.Close()
			err = generateMain(f, templateData)
		}()
		if err != nil {
			return fmt.Errorf(
				"unable to write main.go for %q: %w",
				modPath,
				err,
			)
		}

		// Write the build metadata.
		f, err = os.Create(filepath.Join(buildDir, "metadata.json"))
		if err != nil {
			return fmt.Errorf(
				"unable to create build metadata for %q: %w",
				modPath,
				err,
			)
		}
		func() {
			defer f.Close()
			enc := json.NewEncoder(f)
			enc.SetIndent("", "  ")
			err = enc.Encode(gen.Build{
				ProjectURL:   *repo.HTMLURL,
				GoModPath:    modPath,
				GoModVersion: version,

				// XXX Assumes Caddy module name can be derived from GitHub repository
				// name.
				CaddyDocURL: fmt.Sprintf(
					"https://caddyserver.com/docs/modules/dns.providers.%s",
					repo.GetName(),
				),
			})
		}()
		if err != nil {
			return fmt.Errorf(
				"unable to write build metadata for %q: %w",
				modPath,
				err,
			)
		}
	}

	return nil
}

// Returns all public repositories in the caddy-dns GitHub organization that
// aren't archived, aren't empty, and aren't templates.
func getCaddyDNSRepos(
	ctx context.Context,
	gh *github.Client,
) ([]*github.Repository, error) {
	result := []*github.Repository{}
	for nextPage := 0; ; {
		repos, resp, err := gh.Repositories.ListByOrg(
			ctx,
			"caddy-dns",
			&github.RepositoryListByOrgOptions{
				Type: "public",
				Sort: "full_name",
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    nextPage,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			if repo.GetArchived() || repo.GetIsTemplate() || repo.GetSize() == 0 {
				continue
			}
			result = append(result, repo)
		}
		if resp.NextPage == 0 {
			return result, nil
		}

		nextPage = resp.NextPage
	}
}

// Runs 'go get' on the modules corresponding to the given Go modules.
func goGetModules(repos []*github.Repository) error {
	cmdLine := []string{"go", "get"}
	for _, repo := range repos {
		cmdLine = append(cmdLine, "github.com/caddy-dns/"+repo.GetName()+"@latest")
	}

	// Run `go get`.
	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
	cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	return err
}

// Returns the latest version and go-mod hash of the given Go module. For
// example, "v0.0.0" and "h1:abcd1234=".
//
// Depends on goGetModules being run.
func getGoModuleVersion(
	modPath string,
) (
	version string,
	hash string,
	err error,
) {
	// Use `go list` to get the module version.
	cmd := exec.Command("go", "list", "-f", "{{.Version}}", "-m", modPath)
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("error running 'go list': %w", err)
	}
	version = strings.TrimSpace(string(data))
	if version == "" {
		version = "unknown"
	}

	// Use `go list` to get the module hash.
	cmd = exec.Command("go", "list", "-f", "{{.Sum}}", "-m", modPath)
	cmd.Stderr = os.Stderr
	data, err = cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("error running 'go list': %w", err)
	}
	hash = strings.TrimSpace(string(data))
	if hash == "" {
		hash = "unknown"
	}

	return version, hash, nil
}

func generateMakefile(w io.Writer, templateData TemplateData) error {
	_, err := io.WriteString(w, `# DO NOT EDIT. Generated by internal/gen/builds/generator.go.

`)
	if err != nil {
		return fmt.Errorf("unable to write header: %w", err)
	}
	err = makefileTemplate.Execute(w, templateData)
	if err != nil {
		return fmt.Errorf("unable to execute Makefile template: %w", err)
	}

	return nil
}

func generateMain(w io.Writer, templateData TemplateData) error {
	buf := &bytes.Buffer{}
	_, err := buf.WriteString(`
// DO NOT EDIT. Generated by internal/gen/plugins/generator.go.
`)
	if err != nil {
		return fmt.Errorf("unable to write header to output buffer: %w", err)
	}
	err = mainTemplate.Execute(buf, templateData)
	if err != nil {
		return fmt.Errorf("unable to execute main.go template: %w", err)
	}

	// Format the result using gofmt.
	result, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("unable to format generated main.go: %w", err)
	}

	// Write the formatted result.
	_, err = w.Write(result)
	if err != nil {
		return fmt.Errorf("unable to write generated main.go:%w", err)
	}

	return nil
}
