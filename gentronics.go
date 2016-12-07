package gentronics

import "github.com/pkg/errors"

type Data map[string]interface{}

type Runnable interface {
	Run(string, Data) error
}

type Generator struct {
	Runners []Runnable
}

func New() *Generator {
	return &Generator{
		Runners: []Runnable{},
	}
}

func (g *Generator) Add(r Runnable) {
	g.Runners = append(g.Runners, r)
}

func (g *Generator) Run(dir string, data Data) error {
	for _, r := range g.Runners {
		err := r.Run(dir, data)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
