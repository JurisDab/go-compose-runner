package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed templates
var templatesFS embed.FS

const templatesRoot = "templates"

type Data struct {
	Name string
	Port int
}

func List() ([]string, error) {
	entries, err := fs.ReadDir(templatesFS, templatesRoot)
	if err != nil {
		return nil, fmt.Errorf("reading embedded templates: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func Generate(templateName, destDir string, data Data) error {
	srcRoot := path.Join(templatesRoot, templateName)
	if _, err := fs.Stat(templatesFS, srcRoot); err != nil {
		names, _ := List()
		return fmt.Errorf("unknown template %q (available: %s)", templateName, strings.Join(names, ", "))
	}

	if _, err := os.Stat(destDir); err == nil {
		return fmt.Errorf("destination %s already exists", destDir)
	}

	return fs.WalkDir(templatesFS, srcRoot, func(fsPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(filepath.FromSlash(srcRoot), filepath.FromSlash(fsPath))
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		content, err := templatesFS.ReadFile(fsPath)
		if err != nil {
			return err
		}

		if strings.HasSuffix(target, ".tmpl") {
			target = strings.TrimSuffix(target, ".tmpl")
			return renderFile(target, fsPath, content, data)
		}

		return os.WriteFile(target, content, 0o644)
	})
}

func renderFile(target, srcPath string, content []byte, data Data) error {
	tmpl, err := template.New(srcPath).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", srcPath, err)
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("rendering template %s: %w", srcPath, err)
	}
	return nil
}
