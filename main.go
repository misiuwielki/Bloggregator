package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/misiuwielki/Bloggregator/internal/config"
	"github.com/misiuwielki/Bloggregator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name        string
	description string
	arguments   []string
}

type commands struct {
	lst map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.lst[cmd.name]
	if !ok {
		return fmt.Errorf("uknown command %s", cmd.name)
	}
	return fn(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.lst[name] = f
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("postgres", cfg.Url)
	dbQueries := database.New(db)
	s := state{
		config: &cfg,
		db:     dbQueries}
	cmds := commands{
		lst: map[string]func(*state, command) error{
			"login":     handlerLogin,
			"register":  handlerCreateUser,
			"reset":     handlerReset,
			"users":     handlerGetAllUsers,
			"agg":       handlerAggregate,
			"addfeed":   middleWareLoggedIn(handlerAddFeed),
			"feeds":     handlerGetAllFeeds,
			"follow":    middleWareLoggedIn(handlerFollowFeed),
			"following": middleWareLoggedIn(handlerFollowsForUser),
			"unfollow":  middleWareLoggedIn(handlerUnfollowFeed),
		},
	}
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}
	cmd := command{
		name:      args[1],
		arguments: args[2:],
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}

}
