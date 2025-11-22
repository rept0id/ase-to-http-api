package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func tokenMiddleware(ENV_TOKEN string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var EXPECTED_TOKEN string

		var givenToken string

		/*** * * ***/

		EXPECTED_TOKEN = ENV_TOKEN

		givenToken = c.GetHeader("Authorization")

		/*** * * ***/

		if givenToken == "" {
			c.AbortWithError(
				http.StatusUnauthorized,
				errors.New("Authorization header is required"),
			)

			return
		}

		/*** * * ***/

		if givenToken != "Bearer "+EXPECTED_TOKEN {
			c.AbortWithError(
				http.StatusUnauthorized,
				errors.New("Invalid token"),
			)

			return
		}

		/*** * * ***/

		c.Next()
	}
}

func main() {
	var ENV_TOKEN string

	var ginEngine *gin.Engine

	/*** * * ****/

	ENV_TOKEN = os.Getenv("TOKEN")
	if len(ENV_TOKEN) < 16 {
		panic("TOKEN is too weak")
	}

	ginEngine = gin.Default()
	ginEngine.Use(tokenMiddleware(ENV_TOKEN))

	/*** * * ***/

	ginEngine.GET("/ase", func(c *gin.Context) {
		var err error

		var qIp string
		var qPort string
		var qPortInt int

		var aseRes TypeAseRes

		var res struct {
			TypeRes

			Data     TypeAseRes     `json:"Data"`
			Metadata map[string]any `json:"Metadata"`
		}

		/*** * * ***/

		res.Metadata = make(map[string]any)

		/*** * * ***/

		qIp = c.Query("ip")
		qPort = c.Query("port")
		qPortInt, err = strconv.Atoi(qPort)
		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				errors.New("Invalid port parameter"),
			)
		}

		// aseRes
		aseRes, err = ase(qIp, qPortInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res)
		}

		// res
		res.Data = aseRes

		/*** * * ***/

		c.JSON(http.StatusOK, res)
	})

	/*** * * ***/

	ginEngine.Run(":80")
}
