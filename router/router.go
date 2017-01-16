package router

import (
	"fmt"
	"net/http"

	"github.com/TeaMeow/KitSvc/server"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
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
		message := "OK"
		c.String(http.StatusOK, "\n"+message)
	})
	// Disk check.
	g.GET("/disk", func(c *gin.Context) {
		u, _ := disk.Usage("/")

		usedMB := int(u.Used) / MB
		usedGB := int(u.Used) / GB
		totalMB := int(u.Total) / MB
		totalGB := int(u.Total) / GB
		usedPercent := int(u.UsedPercent)

		status := http.StatusOK
		text := "OK"

		if usedPercent >= 95 {
			status = http.StatusOK
			text = "CRITICAL"
		} else if usedPercent >= 90 {
			status = http.StatusTooManyRequests
			text = "WARNING"
		}

		message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)

		c.String(status, "\n"+message)
	})
	// Disk check.
	g.GET("/cpu", func(c *gin.Context) {
		cores, _ := cpu.Counts(false)

		a, _ := load.Avg()
		l1 := a.Load1
		l5 := a.Load5
		l15 := a.Load15

		status := http.StatusOK
		text := "OK"

		if l5 >= float64(cores-1) {
			status = http.StatusOK
			text = "CRITICAL"
		} else if l5 >= float64(cores-2) {
			status = http.StatusTooManyRequests
			text = "WARNING"
		}

		message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)

		c.String(status, "\n"+message)
	})

	// Event handlers.
	e := eventutil.New(g)
	e.POST("/event-store/user.create/", "user.create", server.CreateUser)

	//e.GET("/kit_metrics", server.)

	return e.Gin, e
}
