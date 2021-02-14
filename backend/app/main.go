package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c9s/gomon/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Message godoc
type Message struct {
	User    string   `json:"user" binding:"required"`
	Message string   `json:"message" binding:"required"`
	Time    JSONTime `json:"time"`
}

// JSONTime godoc
type JSONTime time.Time

// MarshalJSON godoc
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2 Jan 15:04:05"))
	return []byte(stamp), nil
}

// Event godoc
type Event struct {
	Message      chan string
	NewClient    chan Client
	ClosedClient chan Client
	Clients      map[Client]bool
}

func (event *Event) listen() {
	for {
		select {
		case client := <-event.NewClient:
			event.Clients[client] = true
			log.Printf("Client added, total %d", len(event.Clients))

		case client := <-event.ClosedClient:
			delete(event.Clients, client)
			log.Printf("Client removed, total %d", len(event.Clients))

		case message := <-event.Message:
			for client := range event.Clients {
				client <- message
			}
		}
	}
}

// Client godoc
type Client chan string

func main() {
	if currentEnvironment, ok := os.LookupEnv("ENV"); ok {
		if currentEnvironment == "dev" {
			err := godotenv.Load("./.env")
			if err != nil {
				logger.Info("Can't load .env", err)
			}
		}
	}

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	event := &Event{
		Message:      make(chan string),
		NewClient:    make(chan Client),
		ClosedClient: make(chan Client),
		Clients:      make(map[Client]bool),
	}

	go event.listen()

	router.POST("/message", func(c *gin.Context) {
		var message Message
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		loc, _ := time.LoadLocation("Asia/Bangkok")
		message.Time = JSONTime(time.Now().In(loc))
		b, _ := json.Marshal(message)
		go func() {
			event.Message <- string(b)
		}()
	})

	router.GET("/message", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		client := make(Client)
		event.NewClient <- client

		defer func() {
			event.ClosedClient <- client
		}()

		go func() {
			<-c.Done()
			event.ClosedClient <- client
		}()

		c.Stream(func(w io.Writer) bool {
			if message, ok := <-client; ok {
				c.SSEvent("message", message)
				return true
			}
			return false
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)
}
