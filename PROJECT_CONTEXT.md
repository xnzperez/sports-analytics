Sports Analytics (Imagine Cup 2025)

## Resumen Ejecutivo

Sports Analytics es una plataforma Full-Stack orientada a la gestión responsable del bankroll y al análisis predictivo aplicado a apuestas deportivas (especial foco en e-Sports como LoL, Valorant y CS2). Resuelve dos problemas principales del ecosistema de apuestas: (1) falta de visibilidad y métricas accionables sobre el rendimiento de una estrategia individual (win rate, ROI, profit, distribución por deporte) y (2) ausencia de gestión automatizada y segura de las liquidaciones y del ledger financiero. La plataforma permite simular, ejecutar y liquidar apuestas, proporciona recomendaciones "AI Tips" y herramientas de control de riesgo para ayudar al usuario a ser más rentable y responsable.

## Arquitectura del Sistema

- Frontend: Single Page Application desarrollada con React + TypeScript y empaquetada con Vite. Consume la API REST del backend mediante Axios y presenta dashboards con Recharts y componentes (Lucide React). El estado global es manejado por Zustand con persistencia local (auth-storage).
- Backend: API RESTful escrita en Go (Golang) usando Fiber (v2) como framework HTTP y GORM como ORM para acceso transaccional a PostgreSQL. Código modularizado en paquetes internos: auth, betting, market, analytics, platform/database y worker.
- Worker: Un componente interno (goroutine) que actúa como scheduler para procesar apuestas pendientes en background (StartScheduler). Este worker consulta apuestas pendientes, decide el resultado (actualmente por simulación determinista) y resuelve apuestas llamando a los servicios de negocio.
- Integración y despliegue: El backend incluye un Dockerfile y está pensado para desplegarse en Azure Container Apps; la base de datos objetivo es Azure Database for PostgreSQL (Flexible Server). El frontend puede desplegarse como app estática en Vercel o servir archivos estáticos detrás de un CDN.
- Comunicación: Frontend ↔ Backend via HTTPS REST (Axios). El backend usa transacciones ACID en la base de datos para preservar la integridad del ledger y evitar race conditions.

Evidencia en el repositorio:

- backend/cmd/api/main.go — definición de rutas y arranque del worker.
- backend/internal/worker/resolver.go — scheduler que ejecuta processPendingBets cada 10 segundos.
- backend/internal/platform/database/postgres.go — conexión GORM a PostgreSQL y pool config.
- backend/Dockerfile — instrucciones para construir imagen del backend.
- frontend/src/lib/axios.ts — configuración de Axios y uso de VITE_API_URL.

## Stack Tecnológico

- Backend:
  - Lenguaje: Go (go 1.25.x según go.mod)
  - Framework HTTP: github.com/gofiber/fiber/v2
  - ORM: gorm.io/gorm con driver postgres (pgx)
  - Contenerización: Dockerfile multi-stage (build + final)
- Frontend:
  - Vite (vite) para bundling y dev server
  - React 19 + TypeScript
  - UI: Tailwind CSS, Lucide React (iconos), componentes propios
  - State: Zustand (+ persist middleware)
  - HTTP client: Axios (configurado en src/lib/axios.ts con interceptores)
  - Charts: Recharts
  - Validación: zod, react-hook-form

## Lógica de Negocio y Worker

- Flujo de creación de apuesta (resumen técnico):

  1. El frontend POST /api/bets llamará al handler del servicio de apuestas con PlaceBetRequest.
  2. En backend/internal/betting/service.go PlaceBet ejecuta una transacción: obtiene saldo del usuario con row-lock (GetUserBalanceForUpdate), verifica fondos, descuenta el stake y crea la apuesta + entrada de ledger (Transaction) en la misma transacción.
  3. La apuesta queda con estado "pending" hasta su resolución.

- Worker (Auto-Settlement):

  - backend/internal/worker/resolver.go inicia StartScheduler que corre un ticker cada 10 segundos y ejecuta processPendingBets.
  - processPendingBets obtiene apuestas pendientes (GetPendingBets), extrae detalles JSON y determina el ganador (actualmente simulateWinner con hashing determinista del ID para simulación).
  - Para cada apuesta, llama service.ResolveBet(betID, outcome) que delega a repo.ResolveBet para actualizar estado, registrar transacciones de pago y ajustar balances atómicamente.

- Cálculo de estadísticas (Win Rate, Profit, ROI):
  - backend/internal/betting/service.go GetUserStats y GetUserDashboardStats recogen estadísticas "raw" del repositorio y calculan métricas derivadas:
    - WinRate = (Total Won / (Total Won + Total Lost)) \* 100 (evitando dividir por 0).
    - TotalWagered y TotalReturned son valores traídos del repo; NetProfit = TotalReturned - TotalWagered.
    - ROI = (NetProfit / TotalWagered) \* 100 cuando TotalWagered > 0.
  - Dashboard agrega métricas por deporte, calcula profit por apuesta como (stake \* odds - stake) cuando WON y -stake cuando LOST.

## Componentes Clave

- Backend (paquetes principales):

  - auth: registro, login, middleware JWT (auth/middleware.go protege rutas, valida JWT_SECRET).
  - betting: entidades Bet, Transaction; lógica de apuestas y stats; repositorios con transacciones.
  - market: sincronización/obtención de mercados externos (endpoints públicos para listar mercados).
  - platform/database: conexión a PostgreSQL via GORM, configuración del pool.
  - worker: lógica de resolución automática y simulación de partidos.

- Frontend (features):
  - features/auth: formularios, servicios (auth.service.ts), estado (auth.store.ts) con persistencia local.
  - features/betting: PlaceBetModal, BetHistory, BetDetailsModal — componentes para crear y listar apuestas.
  - features/analytics: DashboardCharts, consumo del endpoint /api/stats para KPI y AI Tip.
  - state global: useAuthStore (Zustand persist) para token y perfil; axios interceptor agrega Authorization header.

## Potencial de IA

- Implementación actual de "AI Tips":
  - backend/internal/analytics/advisor.go expone GenerateSmartTip(stats StatsInput) que aplica reglas heurísticas (umbral de mínimo 5 apuestas, detección de varianza, recomendaciones de stake) y devuelve un AdvisorResult con Message y Level. Esta función implementa:
    - Mensajes de aprendizaje (pocos datos).
    - Detección de paradojas (ej. WinRate alto pero Profit negativo) con recomendaciones de reducción de stake.
    - Sugerencia de stake basada en una heurística (ej. suggestedStake = Bankroll \* 0.02, similar a una Kelly simplificada).
- Extensiones posibles (alto impacto para Imagine Cup):
  - Integración con Azure OpenAI (LLMs) para generar explicaciones personalizadas y planes de ajuste de estrategia en lenguaje natural.
  - Modelos ML/Time Series que predigan valor esperado (EV) por mercado, usando features como historial de apuestas, cuotas, mercados y métricas por equipo/ligas.
  - Sistema de "Contra-fraude" con modelos de anomalías (detectan patrones de abuso o manipulación) y recomendaciones de juego responsable (auto-límites, cooling-off).

Cómo ayuda al usuario a ser más rentable y responsable:

- Métricas accionables (WinRate, ROI, Profit por deporte) para identificar mercados rentables.
- Recomendaciones de stake y alertas de varianza evitan sobre-apostar y promueven gestión de bankroll.
- Ledger y transacciones atómicas permiten auditoría y trazabilidad financiera.

## Escalabilidad

- Por qué Go y Azure los considero adecuados:
  - Go (Golang) aporta alta concurrencia (goroutines), bajo uso de memoria y tiempos de arranque rápidos, ideal para microservicios y workers de alta frecuencia.
  - Fiber es un framework ligero y rápido para HTTP en Go, con baja latencia.
  - Docker + Azure Container Apps permite escalar réplicas del backend horizontalmente; el worker puede ejecutarse en réplicas separadas o como Job/Service independiente.
  - GORM + pooling (SetMaxOpenConns/SetMaxIdleConns) y uso de drivers nativos (pgx) ayudan en cargas altas de conexiones.
- Recomendaciones para escalar a nivel global:
  - Separar el worker del API (deploy como servicio de background / Azure Container Instances o Azure Functions Durable / Azure WebJobs) para evitar interferencia en latencia de respuesta.
  - Añadir una cola (Azure Service Bus / RabbitMQ / Azure Queue) para tareas de resolución y permitir procesamiento idempotente y escalable.
  - Cachear lecturas frecuentes (Azure Cache for Redis) para dashboards con alta demanda.
  - Implementar sharding/particionado para el ledger si la carga lo exige; usar read replicas para reporting.

## Seguridad, Observabilidad y Buenas Prácticas

- Seguridad:
  - Autenticación basada en JWT (middleware en auth/middleware.go). JWT_SECRET debe almacenarse en variables de entorno del entorno de ejecución.
  - Operaciones monetarias envueltas en transacciones DB para evitar race conditions. Uso de bloqueos SELECT FOR UPDATE al leer saldo.
  - No almacenar secretos en el repositorio; .env sólo para desarrollo.
- Observabilidad:
  - Logs SQL activados en producción para depuración (logger.Info en GORM). Añadir integraciones con Azure Monitor / Application Insights para métricas, traces y alertas.
- Pruebas y migraciones:
  - AutoMigrate se usa en main.go para crear tablas iniciales; en producción usar migraciones controladas y versionadas.

## Despliegue y Runbook (resumen práctico)

1. Variables de entorno mínimas (ejemplos):
   - DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE
   - JWT_SECRET
   - PORT (opcional, default 3000)
   - VITE_API_URL para frontend en dev
2. Backend (local):
   - cd backend && go mod tidy && go run cmd/api/main.go
3. Backend (docker):
   - docker build -t sports-analytics-backend:latest -f backend/Dockerfile backend
   - docker run -e DB\_\* -e JWT_SECRET -p 3000:3000 sports-analytics-backend:latest
4. Frontend (local):
   - cd frontend && pnpm install && pnpm dev
5. Despliegue en Azure:
   - Empaquetar imagen Docker y enviar a ACR (Azure Container Registry).
   - Crear Azure Container App apuntando a la imagen y configurar variables de entorno y conexiones a Azure Database for PostgreSQL.
   - Para frontend, usar Vercel o Azure Static Web Apps + CDN.

## Análisis de archivos (consultados para verificar la documentación)

Se inspeccionaron los archivos clave del workspace para asegurar precisión en la documentación:

- README.md
- backend/go.mod
- backend/Dockerfile
- backend/cmd/api/main.go
- backend/internal/worker/resolver.go
- backend/internal/betting/service.go
- backend/internal/analytics/advisor.go
- backend/internal/platform/database/postgres.go
- frontend/package.json
- frontend/src/lib/axios.ts
- frontend/src/features/auth/auth.store.ts
- frontend/src/features/auth/services/auth.service.ts
- frontend/src/pages/DashboardPage.tsx
- frontend/src/features/analytics/components/DashboardCharts.tsx
- frontend/src/components y features relevantes (estructura de carpetas analizada)
