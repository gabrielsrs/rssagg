package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabrielsrs/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWIthJSON(w, 201, databaseUserToUser(newUser))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWIthJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}

	respondWIthJSON(w, 200, databasePostsToPosts(posts))
}
