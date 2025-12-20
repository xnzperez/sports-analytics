package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	// Módulos internos
	"github.com/xnzperez/sports-analytics-backend/internal/auth"
	"github.com/xnzperez/sports-analytics-backend/internal/betting"
	"github.com/xnzperez/sports-analytics-backend/internal/platform/database"

	// --- SWAGGER IMPORTS ---
	fiberSwagger "github.com/swaggo/fiber-swagger"        // fiber-swagger middleware
	_ "github.com/xnzperez/sports-analytics-backend/docs" // <--- Importante: Esto carga los archivos generados

	"github.com/joho/godotenv"
)

// @title           Sports Analytics API
// @version         1.0
// @description     API profesional para seguimiento de apuestas y estadísticas deportivas.
// @contact.name    Carlos Pérez
// @contact.email   xnzperez.dev@gmail.com
// @host            localhost:3000
// @BasePath        /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY1MTQ2NDgsInVzZXJfaWQiOiJkYTE4ZGFkZC1lMmVhLTRhMjAtYTNjOS1jNzUxZWQ1ZjFiNWYiLCJ1c2VybmFtZSI6InhuenBlcmV6In0.QqWkJ1FvpMFO8aBf1WJblPYj_ZSkceDSSs4jhwBi7Wg"
func main() {
	// 1. Cargar Variables de Entorno
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  No se encontró archivo .env, usando variables del sistema")
	}

	// 2. Conectar a Base de Datos
	database.Connect()
	database.Migrate()

	// 3. Inicializar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sports Analytics API v1",
	})

	// 4. Middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// 5. INICIALIZACIÓN DE HANDLERS
	authHandler := auth.NewHandler(database.Instance)
	bettingHandler := betting.NewHandler(database.Instance)

	// 6. RUTA DE DOCUMENTACIÓN (SWAGGER)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// 7. DEFINICIÓN DE RUTAS
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, _ := database.Instance.DB()
		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "DB Disconnected"})
		}
		return c.JSON(fiber.Map{"status": "ok", "message": "Systems Operational 🚀"})
	})

	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	api := app.Group("/api", auth.Protected())
	api.Get("/me", authHandler.GetMe)

	// Betting Routes
	api.Post("/bets", bettingHandler.PlaceBet)
	api.Patch("/bets/:id/resolve", bettingHandler.ResolveBetHandler)
	api.Get("/bets", bettingHandler.GetBetsHandler)

	// Analytics & Ledger
	api.Get("/stats", bettingHandler.GetStatsHandler)
	api.Get("/transactions", bettingHandler.GetTransactionsHandler)

	// 8. Arrancar Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Servidor corriendo en puerto " + port)
	log.Fatal(app.Listen(":" + port))
}
