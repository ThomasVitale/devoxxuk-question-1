package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Answers struct {
	Player        string `json:"player"`
	SessionId     string `json:"sessionId"`
	OptionA       bool   `json:"optionA"`
	OptionB       bool   `json:"optionB"`
	OptionC       bool   `json:"optionC"`
	OptionD       bool   `json:"optionD"`
	RemainingTime int    `json:"remainingTime"`
}

type GameScore struct {
	Player     string
	SessionId  string
	Time       time.Time
	Level      string
	LevelScore int
}

var redisHost = os.Getenv("REDIS_HOST") // <hostname>:<port>
var redisPassword = os.Getenv("REDIS_PASSWORD")
var gameEventingEnabled = os.Getenv("GAME_EVENTING_ENABLED")
var sink = os.Getenv("GAME_EVENTING_BROKER_URI")
var cloudEventsEnabled bool = false

func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	if gameEventingEnabled != "" && gameEventingEnabled != "false" {
		cloudEventsEnabled = true
	}

	points := 0
	var answers Answers

	err := json.NewDecoder(req.Body).Decode(&answers)
	if err != nil {
		log.Println("Error while deserializing Answers: ", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if answers.OptionA == true {
		points = 5
	}
	if answers.OptionB == true {
		points = 0
	}
	if answers.OptionC == true {
		points = 0
	}
	if answers.OptionD == true {
		points = 3 // Bonus option
	}

	points += answers.RemainingTime

	score := GameScore{
		Player:     answers.Player,
		SessionId:  answers.SessionId,
		Level:      "devoxxuk-question-1",
		LevelScore: points,
		Time:       time.Now(),
	}
	scoreJson, err := json.Marshal(score)
	if err != nil {
		log.Println("Error while serializing GameScore: ", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = client.RPush("score-"+answers.SessionId, string(scoreJson)).Err()
	if err != nil {
		log.Println("Error while pushing GameScore to Redis: ", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if cloudEventsEnabled {
		emitCloudEvent(scoreJson)
	}

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(res, string(scoreJson))
}

func emitCloudEvent(gs []byte) error {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// Create an Event.
	event := cloudevents.NewEvent()
	newUUID, _ := uuid.NewUUID()
	event.SetID(newUUID.String())
	event.SetTime(time.Now())
	event.SetSource("devoxxuk-question-1")
	event.SetType("GameScoreEvent")
	event.SetData(cloudevents.ApplicationJSON, gs)

	log.Printf("Emitting an Event: %s to SINK: %s", event, sink)

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), sink)

	// Send that Event.
	result := c.Send(ctx, event)
	if result != nil {
		log.Printf("Resutl: %s", result)
		if cloudevents.IsUndelivered(result) {
			log.Printf("failed to send, %v", result)
			return result
		}
	}
	return nil
}
