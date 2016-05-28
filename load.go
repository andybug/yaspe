package main

import "encoding/json"
import "errors"
import "fmt"
import "github.com/garyburd/redigo/redis"
import "io/ioutil"
import "os"
import "regexp"

func loadData(args []string) error {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "yaspe load <dir>")
	}

	dir := args[0]

	// establish redis connection
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Fprintln(os.Stderr, "loadData: " + err.Error())
		return errors.New("loadData: could not connect to redis server")
	}
	defer c.Close()
	c.Do("FLUSHDB")

	// get the list of available seasons
	seasons, err := getSeasons(dir)
	if err != nil {
		return err
	}

	// for each season...
	for _, s := range seasons {
		// read teams into redis
		err := readTeams(dir, s, c)
		if err != nil {
			return err
		}

		// get list of rounds
		rounds, err := getRounds(dir, s)
		if err != nil {
			return err
		}

		// load the games from each round
		for _, r := range rounds {
			err := readRound(dir, s, r, c)
			if err != nil {
				return err
			}
		}

		// add season to list
		c.Do("RPUSH", "seasons", s)
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

func readTeams(dir string, season string, c redis.Conn) error {
	type team struct {
		Uuid string
		Name string
	}

	path := dir + "/" + season + "/teams.json"

	// read json file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "readTeams: " + err.Error())
		return errors.New("readTeams: could not read file " + path)
	}

	// parse json file into list of teams
	teams := make([]team, 0)
	err = json.Unmarshal(file, &teams)
	if err != nil {
		panic(err)
	}

	// for each team, add it to redis and to the set of teams for the season
	for _, t := range teams {
		c.Send("SET", "team:" + t.Uuid, t.Name)
		c.Send("SADD", "teams:" + season, "team:" + t.Uuid)
	}
	c.Flush()

	// catch all of the responses
	for i := 0; i < len(teams); i++ {
		c.Receive()
		c.Receive()
	}

	return nil
}

func readRound(dir string, season string, round string, c redis.Conn) error {
	type teamscore struct {
		Uuid string
		Score int
	}

	type game struct {
		Date string
		Uuid string
		Home teamscore
		Away teamscore
		Neutral bool
	}

	path := dir + "/" + season + "/round" + round + ".json"

	// read json file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "readRound: " + err.Error())
		return errors.New("readRound: could not read file " + path)
	}

	// parse json file into list of games
	games := make([]game, 0)
	err = json.Unmarshal(file, &games)
	if err != nil {
		panic(err)
	}

	for _, g := range games {
		c.Send("HMSET", "game:" + g.Uuid,
			"date", g.Date,
			"away_team", "team:" + g.Away.Uuid,
			"away_score", g.Away.Score,
			"home_team", "team:" + g.Home.Uuid,
			"home_score", g.Home.Score,
			"neutral", g.Neutral)
		c.Send("RPUSH", "games:" + season + ":" + round, "game:" + g.Uuid)
	}
	c.Send("RPUSH", "games:" + season, "games:" + season + ":" + round)
	c.Flush()

	return nil
}
