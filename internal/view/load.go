package view

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// load will load all template files for the view.
// It starts by loading all files found in pathToTemplateFiles.
// It then loads the optional listOfTemplateFiles to customize.
func (v *View) load() (*template.Template, error) {
	fmt.Println("load ----------------------------------------------------------")
	var files []string
	err := filepath.Walk(v.pathToTemplateFiles, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".gohtml" {
			// avoid files in the yield path
			if !strings.Contains(path, "yield") {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// sort the list of templates files help with reproducible behavior.
	// the driving example is two templates in the path that have the same
	// name. we'd like to ensure that they're always loaded in the same order.
	sort.Strings(files)

	// append the list of listOfTemplateFiles that customize the layout for this view
	files = append(files, v.listOfTemplateFiles...)

	// the optional yield file
	if v.Yield != "" {
		files = append(files, v.pathToTemplateFiles+v.Yield)
	}

	// the layout file
	files = append(files, v.layoutFile)

	return template.ParseFiles(files...)
}
