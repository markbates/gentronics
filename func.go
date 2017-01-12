package gentronics

type Func struct {
	Should ShouldFunc
	Runner RunFn
}

func (f *Func) Run(rootPath string, data Data) error {
	if f.Should != nil && !f.Should(data) {
		return nil
	}
	return f.Runner(rootPath, data)
}
