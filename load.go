package main

import "errors"
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
	seasons, err := getSeasons(dir)
	if err != nil {
		return err
	}

	// for each season...
	for _, s := range seasons {
		// read teams into redis
		err := readTeams(dir, s)
		if err != nil {
			return err
		}

		// get list of rounds
		rounds, err := getRounds(dir, s)
		if err != nil {
			return err
		}

		for _, r := range rounds {
			fmt.Println("round", r)
		}
	}

	return nil
}

func getSeasons(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getSeasons: " + err.Error())
		return nil, errors.New("getSeasons: failed to read directory " + dir)
	}

	r, err := regexp.Compile("^\\d{4}$")
	if err != nil {
		fmt.Fprintln(os.Stderr, "getSeasons: " + err.Error())
		return nil, errors.New("getSeasons: regexp compilation failed")
	}

	var seasons []string
	for _, f := range files {
		match := r.MatchString(f.Name())
		if f.IsDir() && match {
			seasons = append(seasons, f.Name())
		}
	}

	return seasons, nil
}

func getRounds(dir string, season string) ([]string, error) {
	files, err := ioutil.ReadDir(dir + "/" + season)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getRounds: " + err.Error())
		return nil, errors.New("getRounds: failed to read directory " + dir)
	}

	r, err := regexp.Compile("^round(\\d{2}).json$")
	if err != nil {
		fmt.Fprintln(os.Stderr, "getRounds: " + err.Error())
		return nil, errors.New("getRounds: regexp compilation failed")
	}

	var rounds []string
	for _, f := range files {
		match := r.MatchString(f.Name())
		if match {
			submatch := r.FindStringSubmatch(f.Name())
			rounds = append(rounds, submatch[1])
		}
	}

	return rounds, nil
}

func readTeams(dir string, season string) error {
	return nil
}

func readRound(dir string, season string) error {
	return nil
}
