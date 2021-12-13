// Package view implements a handler for Go HTML Templates.
// A view is the highest level handler.
// It is responsible for loading all the templates used in the view.
// It also handles the http requests (which it should not).
package view

import (
	"html/template"
	"log"
	"net/http"
)

// View implements a view.
// It stores the templates to use for rendering.
type View struct {
	Name  string
	Yield string // this is the optional override for content

	// path and list both hold templates files to load.
	// we load from the path first, then from the list.
	// that allows us to override templates.
	// when two templates have the same name, the template
	// library discards the first one loaded, keeping only
	// the last template with that name.
	layoutFile          string
	pathToTemplateFiles string
	listOfTemplateFiles []string

	// templateSet holds the set of parsed templates.
	templateSet *template.Template
}

// New creates a new view.
// It will load all the template listOfTemplateFiles along with an optional list of listOfTemplateFiles.
// The intent of that list is to "specialize" the view by providing templates that have a common name but different view logic.
// I think that we pretty much avoid using them.
func New(name string, layoutFile string, yield string, pathToTemplateFiles string, listOfSpecializedTemplates ...string) (*View, error) {
	log.Printf("[views] creating new view: %q\n", name)
	v := &View{
		Name:                name,
		Yield:               yield,
		layoutFile:          layoutFile,
		pathToTemplateFiles: pathToTemplateFiles,
	}
	v.listOfTemplateFiles = append(v.listOfTemplateFiles, listOfSpecializedTemplates...)
	return v, nil
}

// Load loads all template files for the view.
// It starts by loading all files found in pathToTemplateFiles.
// It then loads the optional listOfTemplateFiles to customize.
//
// In a production environment, we would cache the loads.
// For now, though, we reload every time.
// This allows us to easily test changes to templates.
func (v *View) Load() (*template.Template, error) {
	log.Println(v)
	return v.load()
}

func (v *View) Handler() http.HandlerFunc {
	type item struct {
		Label    string
		Link     string
		Children []*item
	}
	var ctx struct {
		Title   string
		SideBar struct {
			LeftMenu  []*item
			RightMenu []*item
		}
	}
	ctx.Title = "Far Horizons"
	ctx.SideBar.LeftMenu = append(ctx.SideBar.LeftMenu, &item{Label: "First Page!", Link: "/"})
	ctx.SideBar.LeftMenu = append(ctx.SideBar.LeftMenu, &item{Label: "Second page", Link: "/"})
	ctx.SideBar.LeftMenu = append(ctx.SideBar.LeftMenu, &item{
		Label: "Third page with subs",
		Link:  "/",
		Children: []*item{
			&item{Label: "First subpage", Link: "/"},
			&item{Label: "Second subpage", Link: "/"},
		},
	})
	ctx.SideBar.LeftMenu = append(ctx.SideBar.LeftMenu, &item{Label: "Fourth page", Link: "/"})
	ctx.SideBar.RightMenu = append(ctx.SideBar.RightMenu, &item{Label: "Sixth page", Link: "/"})
	ctx.SideBar.RightMenu = append(ctx.SideBar.RightMenu, &item{Label: "Seventh page", Link: "/"})
	ctx.SideBar.RightMenu = append(ctx.SideBar.RightMenu, &item{Label: "Another one", Link: "/"})
	ctx.SideBar.RightMenu = append(ctx.SideBar.RightMenu, &item{Label: "The last one", Link: "/"})

	return func(w http.ResponseWriter, r *http.Request) {
		if data, err := v.Render(&ctx); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write(data)
		}
	}
}
