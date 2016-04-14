package main

import (
	"fmt"
    "github.com/gotstago/cards/statemachine"
    
)

// myStateObj is a simple object holding our various StateFn's for our state machine. We could have done this
// example with just functions instead of methods, but this is to show you can do the same with objects holding
// attributes you want to maintain through execution.
type StateObj struct{}

// PrintHello implements StateFn. This will be our starting state.
func (s StateObj) PrintHello() (statemachine.StateFn, error) {
	fmt.Println("Hello ")
	return s.PrintWorld, nil
}

// PrintWorld implements StateFn.
func (s StateObj) PrintWorld() (statemachine.StateFn, error) {
	fmt.Println("World")
	return nil, nil
}

func main() {
	so := StateObj{}

	// Creates a new statemachine executor that will start execution with myStateObj.PrintHello().
	exec := statemachine.New("helloWorld", so.PrintHello)

	// This begins execution and gets our final error state.
	if err := exec.Execute(); err != nil {
		// Do something with the error.
	}
}
