package main

import "fmt"
import "os"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: yaspe <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "load":
		err := loadData(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
