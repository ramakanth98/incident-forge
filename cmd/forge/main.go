package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ramakanth98/incident-forge/internal/orchestrator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: forge <command> [args]")
		fmt.Println("commands: investigate")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "investigate":
		investigateCmd(os.Args[2:])
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(2)
	}
}

func investigateCmd(args []string) {
	fs := flag.NewFlagSet("investigate", flag.ContinueOnError)

	maxEvidence := fs.Int("max-evidence", 500, "max evidence items each agent can scan")
	outDir := fs.String("out", "", "output directory (default: <bundle>/out)")

	if err := fs.Parse(args); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(2)
	}

	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Println("usage: forge investigate [--max-evidence 500] [--out ./out] <bundlePath>")
		os.Exit(2)
	}

	bundlePath := rest[0]

	if err := orchestrator.RunInvestigate(bundlePath, orchestrator.Budgets{MaxEvidencePerAgent: *maxEvidence}, *outDir); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
