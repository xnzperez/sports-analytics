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

	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar Variables de Entorno
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  No se encontró archivo .env, usando variables del sistema")
	}

	// 2. Conectar a Base de Datos
	database.Connect()

	// 3. Ejecutar Migraciones
	database.Migrate()

	// 4. Inicializar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sports Analytics API v1",
	})

	// 5. Middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// 6. INICIALIZACIÓN DE HANDLERS
	// Aquí inyectamos la conexión a DB a cada módulo
	authHandler := auth.NewHandler(database.Instance)
	bettingHandler := betting.NewHandler(database.Instance)

	// 7. DEFINICIÓN DE RUTAS

	// --- A. Rutas Públicas ---
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, _ := database.Instance.DB()
		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "DB Disconnected"})
		}
		return c.JSON(fiber.Map{"status": "ok", "message": "Systems Operational 🚀"})
	})

	// Grupo Auth
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// --- B. Rutas Privadas / Protegidas (Requieren Token) ---
	// Usamos "api :=" una sola vez aquí para agrupar todo lo que requiere JWT
	api := app.Group("/api", auth.Protected())

	// Rutas de Usuario
	api.Get("/me", authHandler.GetMe)

	// Rutas de Apuestas (Betting)
	api.Post("/bets", bettingHandler.PlaceBet)
	api.Get("/bets", bettingHandler.GetBetsHandler)
	// --- NUEVA RUTA AGREGADA ---
	// Usamos PATCH porque estamos actualizando parcialmente el recurso (cambiando el estado)
	// :id es el parámetro que Fiber capturará
	api.Patch("/bets/:id/resolve", bettingHandler.ResolveBetHandler)
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
