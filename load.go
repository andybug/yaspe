package main

import "fmt"
//import "github.com/garyburd/redigo/redis"
import "io/ioutil"
import "os"
import "regexp"

func loadData(args []string) error {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "yaspe load <dir>")
	}

	dir := args[0]

	// get the list of available seasons
	seasons := getSeasons(dir)

	// for each season...
	for _, s := range seasons {
		fmt.Println("season", s)
		rounds := getRounds(dir, s)
		for _, r := range rounds {
			fmt.Println("round", r)
		}
	}

	return nil
}

func getSeasons(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	r, err := regexp.Compile("^\\d{4}$")
	if err != nil {
		panic(err)
	}

	var seasons []string
	for _, f := range files {
		match := r.MatchString(f.Name())
		if f.IsDir() && match {
			seasons = append(seasons, f.Name())
		}
	}

	return seasons
}

func getRounds(dir string, season string) []string {
	files, err := ioutil.ReadDir(dir + "/" + season)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	r, err := regexp.Compile("^round(\\d{2}).json$")
	if err != nil {
		panic(err)
	}

	var rounds []string
	for _, f := range files {
		match := r.MatchString(f.Name())
		if match {
			submatch := r.FindStringSubmatch(f.Name())
			rounds = append(rounds, submatch[1])
		}
	}

	return rounds
}

func readTeams(dir string, season string) error {
	return nil
}

func readRound(dir string, season string) error {
	return nil
}
