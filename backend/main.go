package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"net/http"
	"pcivault-go-vue/pcivault"
)

func main() {
	server := echo.New()
	/*
	 * Typically you would have some kind of middleware to authenticate requests from your front-end.
	 * I'm leaving this out for the sake of brevity.
	 */
	server.GET("/get-secret", func(c echo.Context) error {
		ce, err := pcivault.GetCaptureEndpoint()
		if err != nil {
			err = errors.Wrap(err, "could not get capture endpoint")
			log.Error(err)
			return err
		}
		return c.JSON(http.StatusOK, ce)
	})

	err := server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("Server quit unexpectedly")
	}
}
