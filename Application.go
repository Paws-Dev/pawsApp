package pawsInit

import (
	"errors"
	"fmt"
	"slices"
)

type ApplicationStarter interface {
	InitializeComponents()
	StartComponents() error
	RegisterComponent(component Component, componentName string, dependencyNames []string)
}

type Component interface {
	Initialize(components map[string]Component)
	Start() error
}

type Application struct {
	Components         map[string]Component
	dependencyLists    [][]string
	dependencySequence []string
}

func (a *Application) RegisterComponent(component Component, componentName string, dependencyNames []string) {
	fmt.Printf("Registering component %s with dependencies %s\n", componentName, dependencyNames)
	a.Components[componentName] = component
	dependencyList := slices.Concat([]string{componentName}, dependencyNames)
	a.dependencyLists = append(a.dependencyLists, dependencyList)
}

func BuildInitializationSequence(dependencyLists *[][]string, dependencySequence *[]string) {
	for len(*dependencyLists) != 0 {
		length := len(*dependencySequence)
		for i, d := range *dependencyLists {
			if len(d) == 1 {
				*dependencySequence = append(*dependencySequence, d[0])
				*dependencyLists = slices.Delete(*dependencyLists, i, i+1)
			}

		}

		for i, d := range *dependencyLists {
			for j, k := range d {
				if j != 0 && slices.Contains(*dependencySequence, k) {
					(*dependencyLists)[i] = slices.Delete((*dependencyLists)[i], j, j+1)
				}
			}
		}

		if len(*dependencySequence) <= length {
			fmt.Println("Initialization sequence build failed")
			fmt.Println("Initializable components: ", *dependencySequence)
			fmt.Println("UnInitializable components: ", *dependencyLists)
			panic("Initialization error, invalid initialization sequence")
		}
	}
	fmt.Println("Dependency initialization sequence build ok:", *dependencySequence)

}

func (a *Application) InitializeComponents() {
	BuildInitializationSequence(&a.dependencyLists, &a.dependencySequence)
	for _, name := range a.dependencySequence {
		fmt.Printf("Initializing component %s\n", name)
		a.Components[name].Initialize(a.Components)
	}
}

func (a *Application) StartComponents() error {

	for name, component := range a.Components {
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
		Components:         make(map[string]Component),
		dependencyLists:    [][]string{},
		dependencySequence: []string{},
	}
}
