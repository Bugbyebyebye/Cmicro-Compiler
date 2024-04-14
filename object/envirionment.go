package object

/**
 * @Description: 环境
 */

type Environment struct {
	store map[string]Object
	outer *Environment //外部环境,用于拓展当前环境
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// NewEnclosedEnvironment  创建闭包环境
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.store = make(map[string]Object)
	for key, val := range outer.store {
		env.store[key] = val
	}
	return env
}
