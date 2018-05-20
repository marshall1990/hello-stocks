package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting app server " + os.Args[0] + " 🎉")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	api := r.Group("/v1")

	api.GET("/quotes/:ticker", func(c *gin.Context) {
		apiUrl := "https://query1.finance.yahoo.com/v7/finance/quote?symbols="
		res, err := http.Get(apiUrl + c.Param("ticker"))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		c.String(http.StatusOK, string(body))
	})

	hostport := ":8080"
	if value := os.Getenv("HOSTPORT"); len(value) != 0 {
		hostport = value
	}

	server := http.Server{Addr: hostport, Handler: r}
	go func() {
		log.Printf("Listening on %s", hostport)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("🔥 fatal error: %s\n", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Bye 👋")
}
