package pawsStart

import (
	"errors"
	"fmt"
	"slices"
)

type ApplicationStarter interface {
	InitializeDependencies()
	StartApplication() error
	RegisterDependency(component Dependency, name string, dependencyNames []string)
}

type Dependency interface {
	Initialize(dependencies map[string]Dependency)
	Start() error
}

func As[T any](dependency Dependency) T {
	component, ok := dependency.(T)
	if !ok {
		panic("Type conversion error during initialization")
	}
	return component
}

type Application struct {
	Dependencies       map[string]Dependency
	dependencyLists    [][]string
	dependencySequence []string
}

func (a *Application) RegisterDependency(dependency Dependency, name string, dependencyNames []string) {
	fmt.Printf("Registering dependency %s with dependencies %s\n", name, dependencyNames)
	a.Dependencies[name] = dependency
	dependencyList := slices.Concat([]string{name}, dependencyNames)
	a.dependencyLists = append(a.dependencyLists, dependencyList)
}

func BuildInitializationSequence(dependencyLists [][]string, dependencySequence *[]string) {
	for len(dependencyLists) != 0 {
		length := len(*dependencySequence)
		for i, d := range dependencyLists {
			if len(d) == 1 {
				*dependencySequence = append(*dependencySequence, d[0])
				dependencyLists = slices.Delete(dependencyLists, i, i+1)
			}

		}

		for i, d := range dependencyLists {
			for j, k := range d {
				if j != 0 && slices.Contains(*dependencySequence, k) {
					(dependencyLists)[i] = slices.Delete((dependencyLists)[i], j, j+1)
				}
			}
		}

		if len(*dependencySequence) <= length {
			fmt.Println("Initialization sequence build failed")
			fmt.Println("Initializable dependencies: ", *dependencySequence)
			fmt.Println("UnInitializable dependencies: ", dependencyLists)
			panic("Initialization error, invalid initialization sequence")
		}
	}
	fmt.Println("Dependency initialization sequence build ok:", *dependencySequence)

}

func (a *Application) InitializeDependencies() {
	BuildInitializationSequence(a.dependencyLists, &a.dependencySequence)
	for _, name := range a.dependencySequence {
		fmt.Printf("Initializing dependency %s\n", name)
		a.Dependencies[name].Initialize(a.Dependencies)
	}
}

func (a *Application) StartApplication() error {

	for name, component := range a.Dependencies {
		fmt.Printf("Starting component %s\n", name)
		err := component.Start()
		if err != nil {
			return errors.New("Component " + name + " failed to start: " + err.Error())
		}
	}
	return nil
}

func NewApplication() *Application {
	return &Application{
		Dependencies:       make(map[string]Dependency),
		dependencyLists:    [][]string{},
		dependencySequence: []string{},
	}
}
