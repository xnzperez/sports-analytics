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

	// 👇 NUEVO IMPORT: El Worker Automático
	"github.com/xnzperez/sports-analytics-backend/internal/worker"

	// --- SWAGGER IMPORTS ---
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/xnzperez/sports-analytics-backend/docs"
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
// @description "Inserte el token JWT con el prefijo Bearer. Ejemplo: 'Bearer eyJhbGci...'"
func main() {
	// 1. Cargar Variables de Entorno
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  No se encontró archivo .env, usando variables del sistema")
	}

	// 2. Conectar a Base de Datos
	database.Connect()

	// Migrar la Nueva Tabla
	database.Instance.AutoMigrate(&auth.User{}, &betting.Bet{}, &betting.Transaction{}, &market.Match{})

	// 3. Inicializar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sports Analytics API v1",
	})

	// 4. Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// CORS Configurado explícitamente para tu Frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 5. INICIALIZACIÓN DE HANDLERS
	authHandler := auth.NewHandler(database.Instance)
	bettingHandler := betting.NewHandler(database.Instance)
	marketHandler := market.NewHandler(database.Instance)

	// 👇 AQUÍ ARRANCAMOS EL MOTOR AUTOMÁTICO 👇
	// Le pasamos el servicio (usando el método GetService que creamos en el paso anterior)
	// Esto inicia el proceso en segundo plano sin detener el servidor.
	worker.StartScheduler(bettingHandler.GetService())

	// 6. RUTA DE DOCUMENTACIÓN (SWAGGER)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// 7. DEFINICIÓN DE RUTAS

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, _ := database.Instance.DB()
		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "DB Disconnected"})
		}
		return c.JSON(fiber.Map{"status": "ok", "message": "Systems Operational 🚀"})
	})

	// Rutas Públicas (Auth)
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// --- RUTAS DE MARKET PÚBLICAS (TEMPORALES) ---
	app.Post("/api/test-sync", marketHandler.SyncMarketsHandler)
	app.Get("/api/markets", marketHandler.ListMarketsHandler)
	app.Post("/api/admin/resolve", bettingHandler.SettleMatchHandler)

	// --- RUTAS PROTEGIDAS (API) ---
	api := app.Group("/api", auth.Protected())

	// User Routes
	api.Get("/me", authHandler.GetMe)

	// Betting Routes (Apuestas)
	api.Post("/bets", bettingHandler.PlaceBet)
	api.Get("/bets", bettingHandler.GetBetsHandler)
	api.Patch("/bets/:id/resolve", bettingHandler.ResolveBetHandler)

	// Analytics & Ledger
	api.Get("/stats", bettingHandler.GetStatsHandler)
	api.Get("/transactions", bettingHandler.GetTransactionsHandler)

	// Market Routes (Partidos)
	api.Post("/admin/sync", marketHandler.SyncMarketsHandler)

	// 8. Arrancar Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Servidor corriendo en puerto " + port)
	log.Fatal(app.Listen(":" + port))
}
