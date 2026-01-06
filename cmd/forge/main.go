package main

import (
	"fmt"
	"os"

	"github.com/ramakanth98/incident-forge/internal/orchestrator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: forge <investigate> <path-to-incident-bundle>")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "investigate":
		if len(os.Args) < 3 {
			fmt.Println("usage: forge investigate ./testdata/incidents/incident-001")
			os.Exit(2)
		}
		path := os.Args[2]

		if err := orchestrator.RunInvestigate(path); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(2)
	}
}
