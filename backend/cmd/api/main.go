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
	"github.com/xnzperez/sports-analytics-backend/internal/market"
	"github.com/xnzperez/sports-analytics-backend/internal/platform/database"
	"github.com/xnzperez/sports-analytics-backend/internal/worker"

	// --- SWAGGER IMPORTS ---
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/xnzperez/sports-analytics-backend/docs"
)

// @title           Sports Analytics API
// @version         1.0
// @description     API profesional para seguimiento de apuestas y estadísticas deportivas. Imagine Cup 2025 Edition.
// @contact.name    Carlos Pérez
// @contact.email   xnzperez.dev@gmail.com
// @host            sports-analytics-backend.azurewebsites.net
// @BasePath        /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// 1. Cargar Variables de Entorno
	// En producción (Azure Container Apps), el archivo .env no existe, se usan vars del sistema.
	if err := godotenv.Load(".env"); err != nil {
		log.Println("ℹ️  Info: No .env file found, using system environment variables (Cloud Mode)")
	}

	// 2. Conectar a Base de Datos
	database.Connect()

	// Migrar la Nueva Tabla (AutoMigrate es seguro si los structs están bien definidos)
	database.Instance.AutoMigrate(&auth.User{}, &betting.Bet{}, &betting.Transaction{}, &market.Match{})

	// 3. Inicializar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sports Analytics API v1",
	})

	// 4. Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// CORS: Vital para que tu Frontend en Vercel pueda hablar con el Backend en Azure
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://sports-analytics-eight.vercel.app",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
	}))

	// 5. INICIALIZACIÓN DE HANDLERS
	authHandler := auth.NewHandler(database.Instance)
	bettingHandler := betting.NewHandler(database.Instance)
	marketHandler := market.NewHandler(database.Instance)

	// 🔄 MOTOR AUTOMÁTICO (WORKER)
	// Inicia el proceso en segundo plano para resolver apuestas y simular partidos.
	worker.StartScheduler(bettingHandler.GetService())

	// 6. RUTA DE DOCUMENTACIÓN (SWAGGER)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// 7. DEFINICIÓN DE RUTAS

	// --- Health Check (Vital para Azure Container Apps) ---
	// Azure llamará a esto para saber si tu contenedor está vivo.
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := database.Instance.DB()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "DB Connection Fail"})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "DB Ping Fail"})
		}
		return c.JSON(fiber.Map{"status": "ok", "message": "Systems Operational 🚀"})
	})

	// --- Grupo de Autenticación (Público) ---
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// --- Grupo de API (Público / Mixto) ---
	apiPublic := app.Group("/api")
	apiPublic.Get("/markets", marketHandler.ListMarketsHandler) // El frontend necesita ver partidos sin login a veces, o puedes protegerlo.

	// --- RUTAS PROTEGIDAS (Requieren Token JWT) ---
	api := app.Group("/api", auth.Protected())

	// Perfil
	api.Get("/me", authHandler.GetMe)

	// Apuestas
	api.Post("/bets", bettingHandler.PlaceBet)
	api.Get("/bets", bettingHandler.GetBetsHandler)
	api.Patch("/bets/:id/resolve", bettingHandler.ResolveBetHandler)

	// Finanzas & Stats
	api.Get("/stats", bettingHandler.GetStatsHandler)
	api.Get("/transactions", bettingHandler.GetTransactionsHandler)

	// Admin (Protegido)
	// Eliminamos /sync-ahora público. Usamos este endpoint seguro si necesitamos forzar.
	api.Post("/admin/sync", marketHandler.SyncMarketsHandler)
	api.Post("/admin/resolve", bettingHandler.SettleMatchHandler)

	// 8. Arrancar Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Servidor corriendo en puerto " + port)
	log.Fatal(app.Listen(":" + port))
}
