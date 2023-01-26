package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"pcivault-go-vue/pcivault"
	"time"
)

var tokens []pcivault.TokenData

func init() {
	http.DefaultClient.Timeout = time.Second * 5

	var err error
	tokens, err = pcivault.GetVaultTokens()
	if err != nil {
		log.Errorf("failed to get tokens from PCI Vault")
	}
}

func main() {
	server := echo.New()
	/*
	 * Typically you would have some kind of middleware to authenticate requests from your front-end.
	 * I'm leaving this out for the sake of brevity.
	 */
	server.GET("/get-secret", func(c echo.Context) error {
		ce, err := pcivault.GetCaptureEndpoint()
		if err != nil {
			err = fmt.Errorf("could not get capture endpoint: %w", err)
			log.Error(err)
			return err
		}
		return c.JSON(http.StatusOK, ce)
	})
	server.GET("/get-tokens", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"tokens": tokens})
	})
	server.POST("/stripe", func(c echo.Context) error {
		var td pcivault.TokenData
		err := c.Bind(&td)
		if err != nil {
			return err
		}

		err = pcivault.PostToStripe(td)
		return err
	})

	err := server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("Server quit unexpectedly")
	}
}
