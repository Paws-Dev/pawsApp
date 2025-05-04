package app

import (
	"fmt"
	"reflect"
	"slices"
)

type Application struct {
	deps     []func() *Dependency
	List     map[string]any
	Cfg      *Config
	init     map[string]func(cfg *Config, list map[string]any) (any, error)
	start    map[string]func(cfg *Config, dep any) error
	depLists [][]string
	depSeq   []string
}

type Dependency struct {
	name  string
	deps  []string
	init  func(cfg *Config, list map[string]any) (any, error)
	start func(cfg *Config, dep any) error
}

func NewDependency(name string, deps []string,
	init func(cfg *Config, list map[string]any) (any, error),
	start func(cfg *Config, dep any) error) *Dependency {
	return &Dependency{
		name:  name,
		deps:  deps,
		init:  init,
		start: start,
	}
}

func New() *Application {
	return &Application{
		deps:     make([]func() *Dependency, 0),
		List:     make(map[string]any),
		Cfg:      NewConfig(),
		init:     make(map[string]func(cfg *Config, list map[string]any) (any, error)),
		start:    make(map[string]func(cfg *Config, dep any) error),
		depLists: [][]string{},
		depSeq:   []string{},
	}
}

func (a *Application) Register(dep func() *Dependency) {
	a.deps = append(a.deps, dep)
}

func (a *Application) Start() {
	for _, dep := range a.deps {
		dependency := dep()
		fmt.Printf("Registering depndency %s with dependencies %s\n", dependency.name, dependency.deps)
		if len(dependency.deps) != 0 {
			depList := slices.Concat([]string{dependency.name}, dependency.deps)
			a.depLists = append(a.depLists, depList)
		} else {
			a.depLists = append(a.depLists, []string{dependency.name})
		}
		a.init[dependency.name] = dependency.init
		if dependency.start != nil {
			a.start[dependency.name] = dependency.start
		}
	}
	BuildInitSeq(a.depLists, &a.depSeq)
	for _, name := range a.depSeq {
		fmt.Printf("Initializing depndency %s\n", name)
		dep, err := a.init[name](a.Cfg, a.List)
		if err != nil {
			fmt.Printf("Error nitializing depndency %s\n", name)
			panic(err)
		}
		a.List[name] = dep
	}
	for name, start := range a.start {
		fmt.Printf("Starting component %s\n", name)
		err := start(a.Cfg, a.List[name])
		if err != nil {
			fmt.Printf("Error starting component %s\n", name)
			panic(err)
		}
	}
}

func BuildInitSeq(depLists [][]string, depSeq *[]string) {
	for len(depLists) != 0 {
		removed := 0
		for i, d := range depLists {
			if len(d) == 1 {
				*depSeq = append(*depSeq, d[0])
				depLists = slices.Delete(depLists, i, i+1)
				removed++
			}
		}
		for i, d := range depLists {
			for j, k := range d {
				if j != 0 && slices.Contains(*depSeq, k) {
					(depLists)[i] = slices.Delete((depLists)[i], j, j+1)
					removed++
				}
			}
		}
		if removed == 0 {
			fmt.Println("Init sequence build failed")
			fmt.Println("Initable dependencies: ", *depSeq)
			fmt.Println("UnInitable dependencies: ", depLists)
			panic("Dependencies initialization error, initialization sequence contains cycles or unresolvable dependencies")
		}
	}
	fmt.Println("Dependency initialization sequence build ok:", *depSeq)
}

func As[T any](dep any) T {
	component, ok := dep.(T)
	if !ok {
		argumentType := reflect.TypeOf(dep)
		var t T
		targetType := reflect.TypeOf(t)
		panic(fmt.Sprintf("Type conversion error during init: cannot convert from %v to %v",
			argumentType, targetType))
	}
	return component
}
