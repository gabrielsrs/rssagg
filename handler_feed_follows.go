package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabrielsrs/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprint("Error parsing JSON: ", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	respondWIthJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}

	respondWIthJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")

	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Couldn't parse feed follows Id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWIthJSON(w, 400, fmt.Sprintf("Couldn't delete feed follows: %v", err))
		return
	}

	respondWIthJSON(w, 204, struct{}{})
}
