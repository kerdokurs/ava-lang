package main

type Environment[T any] struct {
	envs []map[string]*T
}

func NewEnvironment[T any]() *Environment[T] {
	envs := make([]map[string]*T, 0)

	env := &Environment[T]{
		envs: envs,
	}
	env.EnterBlock()

	return env
}

func (e *Environment[T]) EnterBlock() {
	e.envs = append(e.envs, make(map[string]*T))
}

func (e *Environment[T]) ExitBlock() {
	e.envs = e.envs[:len(e.envs)-1]
}

func (e *Environment[T]) Declare(variable string) {
	e.DeclareAssign(variable, nil)
}

func (e *Environment[T]) Assign(variable string, value T) {
	env := e.findEnv(variable)
	(*env)[variable] = &value
}

func (e *Environment[T]) DeclareAssign(variable string, value *T) {
	e.envs[len(e.envs)-1][variable] = value
}

func (e *Environment[T]) Get(variable string) *T {
	return e.envs[len(e.envs)-1][variable]
}

func (e *Environment[T]) findEnv(variable string) *map[string]*T {
	k := len(e.envs) - 1
	for k >= 0 {
		if _, ok := e.envs[k][variable]; ok {
			return &e.envs[k]
		}
		k--
	}

	return nil
}
