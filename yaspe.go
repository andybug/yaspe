package main

import "fmt"
import "os"
//import "github.com/garyburd/redigo/redis"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: yaspe <command>")
		os.Exit(1)
	}
}
