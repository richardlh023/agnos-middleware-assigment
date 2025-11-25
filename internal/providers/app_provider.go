package providers

import (
	"agnos-middleware/internal/configs"
	"fmt"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	Config *configs.ApplicationConfig
}

func NewApp(config *configs.ApplicationConfig) *App {
	fmt.Println("Initializing App...")

	router := gin.Default()

	app := &App{
		Router: router,
		Config: config,
	}

	app.InitializeMiddlewares()

	return app
}

func (a *App) InitializeMiddlewares() {
	fmt.Println("Setting up middlewares...")
}

func (a *App) InitializeControllers(controllers []interface{}) {
	fmt.Println("Setting up controllers...")
	a.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Agnos Middleware API is running",
		})
	})
}

func (a *App) Listen() {
	port := a.Config.App.Port
	fmt.Printf("Server listening on port %s\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)

	// Start server (like: this.app.listen(port))
	if err := a.Router.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
