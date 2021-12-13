package view

import (
	"bytes"
)

// Render uses the data parameter to render a view.
// It returns the resulting bytes or an error.
func (v *View) Render(data interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	if t, err := v.Load(); err != nil {
		return nil, err
	} else if err = t.ExecuteTemplate(b, v.Name, data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
