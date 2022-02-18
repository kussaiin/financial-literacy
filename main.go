package main

import (
	"github.com/gin-gonic/gin"
	"github.com/user/financial-literacy/routes"
	"log"
	"os"
)

func main()  {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.Default()
	routes.HandleRequests(router)
	log.Fatal(router.Run(":" + port))
}