package app

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
)

type AppStarter interface {
	InitDeps()
	StartApp() error
	RegisterDep(dep Dependency, name string, depNames []string)
}

type Dependency interface {
	Init(deps map[string]Dependency)
	Start() error
}

func As[T any](dep Dependency) T {
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

type App struct {
	Deps     map[string]Dependency
	depLists [][]string
	depSeq   []string
}

func (a *App) RegisterDep(dep Dependency, name string, depNames []string) {
	fmt.Printf("Registering dep %s with dependencies %s\n", name, depNames)
	a.Deps[name] = dep
	depList := slices.Concat([]string{name}, depNames)
	a.depLists = append(a.depLists, depList)
}

func BuildInitSeq(depLists [][]string, depSeq *[]string) {
	for len(depLists) != 0 {
		length := len(*depSeq)
		for i, d := range depLists {
			if len(d) == 1 {
				*depSeq = append(*depSeq, d[0])
				depLists = slices.Delete(depLists, i, i+1)
			}

		}

		for i, d := range depLists {
			for j, k := range d {
				if j != 0 && slices.Contains(*depSeq, k) {
					(depLists)[i] = slices.Delete((depLists)[i], j, j+1)
				}
			}
		}

		if len(*depSeq) <= length {
			fmt.Println("Initialization sequence build failed")
			fmt.Println("Initializable dependencies: ", depSeq)
			fmt.Println("UnInitializable dependencies: ", depLists)
			panic("Initialization error, invalid initialization sequence")
		}
	}
	fmt.Println("Dependency initialization sequence build ok:", depSeq)

}

func (a *App) InitDeps() {
	BuildInitSeq(a.depLists, &a.depSeq)
	for _, name := range a.depSeq {
		fmt.Printf("Initializing dependency %s\n", name)
		a.Deps[name].Init(a.Deps)
	}
}

func (a *App) StartApp() error {
	for name, component := range a.Deps {
		fmt.Printf("Starting component %s\n", name)
		err := component.Start()
		if err != nil {
			return errors.New("Component " + name + " failed to start: " + err.Error())
		}
	}
	return nil
}

func NewApp() *App {
	return &App{
		Deps:     make(map[string]Dependency),
		depLists: [][]string{},
		depSeq:   []string{},
	}
}
