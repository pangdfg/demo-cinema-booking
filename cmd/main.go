package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	redis "github.com/pangdfg/cinema/internal/adapters"
	"github.com/pangdfg/cinema/internal/booking"
)

func main() {
	r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

	r.Static("/v1", "./static")
	rdb, err := redis.NewClient("localhost:6379")
	if err != nil {
		log.Fatalf("redis error: %v", err)
	}
	store := booking.NewRedisStore(rdb)
	svc := booking.NewService(store)
	h := booking.NewHandler(svc)

	r.GET("/movies", listMovies)

	r.GET("/movies/:movieID/seats", h.ListSeats)
	r.POST("/movies/:movieID/seats/:seatID/hold", h.HoldSeat)

	r.PUT("/sessions/:sessionID/confirm", h.ConfirmSession)
	r.DELETE("/sessions/:sessionID", h.ReleaseSession)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

    r.Run() 
}


var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}


func listMovies(c *gin.Context) {
	c.JSON(http.StatusOK, movies)
}


type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}