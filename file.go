package gentronics

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gobuffalo/velvet"
	"github.com/pkg/errors"
)

type ShouldFunc func(Data) bool

type File struct {
	Path          string
	Template      string
	TemplateFuncs template.FuncMap
	Permission    os.FileMode
	Should        ShouldFunc
}

func (f *File) Run(rootPath string, data Data) error {
	if !f.Should(data) {
		return nil
	}

	path, err := f.render(f.Path, data)
	if err != nil {
		return err
	}

	body, err := f.render(f.Template, data)
	if err != nil {
		return err
	}

	return f.save(rootPath, path, body)
}

func (f *File) save(rootPath, path, body string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(filepath.Join(rootPath, dir), 0755)
	if err != nil {
		return err
	}

	odir := filepath.Join(rootPath, path)
	fmt.Printf("--> %s\n", odir)

	ff, err := os.Create(odir)
	if err != nil {
		return err
	}

	_, err = ff.WriteString(body)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (f *File) render(s string, data Data) (string, error) {
	t, err := velvet.Parse(s)
	if err != nil {
		return "", err
	}
	t.Helpers.AddMany(f.TemplateFuncs)
	return t.Exec(data.ToVelvet())
}

func NewFile(path string, t string) *File {
	return &File{
		Path:          path,
		Template:      t,
		TemplateFuncs: map[string]interface{}{},
		Permission:    0664,
		Should: func(data Data) bool {
			return true
		},
	}
}
