package router

import (
	"net/http"
	"strconv"

	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/diskutil"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/gin-gonic/gin"
)

func Load(middleware ...gin.HandlerFunc) (http.Handler, *eventutil.Engine) {

	// Gin engine and middlewares.
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(middleware...)

	// Command routes.
	g.POST("/user", server.CreateUser)
	g.GET("/user/:id", server.GetUser)
	g.DELETE("/user/:id", server.DeleteUser)
	g.PUT("/user/:id", server.UpdateUser)

	// Health check.
	g.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "\nPONG!")
	})
	// Disk check.
	g.GET("/disk", func(c *gin.Context) {
		disk := diskutil.Usage("/")

		freeMB := strconv.Itoa(int(disk.Free)) + " MB"
		freeGB := strconv.Itoa(int(disk.FreeGB())) + " GB"
		usedMB := strconv.Itoa(int(disk.Used)) + " MB"
		usedGB := strconv.Itoa(int(disk.UsedGB())) + " GB"
		usedPercentage := int(disk.UsedPercentage())
		usedPercentageString := strconv.Itoa(usedPercentage)
		freePercentage := 100 - usedPercentage

		status := http.StatusOK
		text := "OK"

		switch {
		case freePercentage < 5:
			status = http.StatusTooManyRequests
			text = "WARNING"
		case freePercentage < 10:
			status = http.StatusServiceUnavailable
			text = "CRITICAL"
		}

		message := "\nDISK " + text + " - free space: " + freeMB + "(" + freeGB + ") / " + usedMB + "(" + usedGB + ")" + " | Used: " + usedPercentageString + "%"

		c.String(status, message)
	})

	// Event handlers.
	e := eventutil.New(g)
	e.POST("/event-store/user.create/", "user.create", server.CreateUser)

	//e.GET("/kit_metrics", server.)

	return e.Gin, e
}
