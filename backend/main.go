package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "PostEaze backend running with Gin!",
        })
    })

    router.Run(":8080")
}
