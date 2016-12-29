package gentronics

import (
	"os"

	"github.com/gobuffalo/velvet"
	"github.com/pkg/errors"
)

type Data map[string]interface{}

func (d Data) ToVelvet() *velvet.Context {
	return velvet.NewContextWith(d)
}

type RunFn func(string, Data) error

type Runnable interface {
	Run(string, Data) error
}

type Generator struct {
	Runners []Runnable
	Should  ShouldFunc
}

func New() *Generator {
	return &Generator{
		Runners: []Runnable{},
	}
}

func (g *Generator) Add(r Runnable) {
	g.Runners = append(g.Runners, r)
}

func (g *Generator) Run(rootPath string, data Data) error {
	if g.Should != nil {
		b := g.Should(data)
		if !b {
			return nil
		}
	}
	err := os.MkdirAll(rootPath, 0755)
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.Chdir(rootPath)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, r := range g.Runners {
		err := r.Run(rootPath, data)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
