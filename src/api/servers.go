package api

/**
 * servers.go - /servers rest api implementation
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notional-labs/gobetween/src/config"
	"github.com/notional-labs/gobetween/src/manager"
	"github.com/notional-labs/gobetween/src/stats"
)

/**
 * Attaches /servers handlers
 */
func attachServers(app *gin.RouterGroup) {
	/**
	 * Find all current configured servers
	 */
	app.GET("/servers", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, manager.All())
	})

	/**
	 * Find server by name
	 */
	app.GET("/servers/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.IndentedJSON(http.StatusOK, manager.Get(name))
	})

	/**
	 * Delete server by name
	 */
	app.DELETE("/servers/:name", func(c *gin.Context) {
		name := c.Param("name")
		manager.Delete(name) //nolint:errcheck
		c.IndentedJSON(http.StatusOK, nil)
	})

	/**
	 * Create new server with name :name
	 */
	app.POST("/servers/:name", func(c *gin.Context) {
		name := c.Param("name")

		cfg := config.Server{}
		if err := c.BindJSON(&cfg); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := manager.Create(name, cfg); err != nil {
			c.IndentedJSON(http.StatusConflict, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, nil)
	})

	/**
	 * Get server stats
	 */
	app.GET("/servers/:name/stats", func(c *gin.Context) {
		name := c.Param("name")
		c.IndentedJSON(http.StatusOK, stats.GetStats(name))
	})
}
