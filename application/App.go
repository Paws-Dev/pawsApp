package app

import (
	"fmt"
	"reflect"
	"slices"
)

type Application struct {
	dep      map[string]any
	init     map[string]func(name string, dep map[string]any) (func(any) error, error)
	start    map[string]func(any) error
	depLists [][]string
	depSeq   []string
}

func New() *Application {
	return &Application{
		dep:      make(map[string]any),
		init:     make(map[string]func(name string, dep map[string]any) (func(any) error, error)),
		start:    make(map[string]func(any) error),
		depLists: [][]string{},
		depSeq:   []string{},
	}
}

func (s *Application) InitComponent(name string, init func(name string, dep map[string]any) (func(any) error, error), dependencies ...string) {
	fmt.Printf("Registering dep %s with dependencies %s\n", name, dependencies)
	s.init[name] = init
	if len(dependencies) != 0 {
		depList := slices.Concat([]string{name}, dependencies)
		s.depLists = append(s.depLists, depList)
	} else {
		s.depLists = append(s.depLists, []string{name})
	}

}

func (s *Application) Start() {
	BuildInitSeq(s.depLists, &s.depSeq)
	for _, name := range s.depSeq {
		fmt.Printf("Initializing dependency %s\n", name)
		start, err := s.init[name](name, s.dep)
		if err != nil {
			fmt.Printf("Error nitializing dependency %s\n", name)
			panic(err)
		}
		if start != nil {
			s.start[name] = start
		}
	}
	for name, start := range s.start {
		fmt.Printf("Starting dependency %s\n", name)
		err := start(s.dep[name])
		if err != nil {
			fmt.Printf("Error starting dependency %s\n", name)
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
			fmt.Println("Initable components: ", *depSeq)
			fmt.Println("UnInitable components: ", depLists)
			panic("Components initialization error, initialization sequence contains cycles or unresolvable dependencies")
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
