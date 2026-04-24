package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/misiuwielki/Bloggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username wasn't passed")
	}
	username := cmd.arguments[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not registered")
		} else {
			log.Printf("error on login: %v", err)
			return fmt.Errorf("couldn't login")
		}
	}
	s.config.SetUser(user.Name)
	fmt.Printf("%s was set as current user\n", user.Name)
	return nil
}

func handlerCreateUser(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username wasn't passed")
	}
	username := cmd.arguments[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: username,
	})
	if err != nil {
		log.Printf("error on creating user: %s", err)
		return fmt.Errorf("couldn't create user")
	}
	s.config.SetUser(user.Name)
	fmt.Printf("%s account was created and logged in\n", user.Name)
	log.Printf("%v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetDatabase(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerGetAllUsers(s *state, cmd command) error {
	userS, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Printf("error on getting users: %s", err)
		return fmt.Errorf("couldn't process request")
	}
	for _, user := range userS {
		if user.Name == s.config.Current_username {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

func middleWareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.Current_username)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
