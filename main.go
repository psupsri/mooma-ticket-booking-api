package main

import (
	"log"
	"mooma-api/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	c := config.Load()
	db, err := config.ConnectDatabase(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	_ = db

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.Run()
}
