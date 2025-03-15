package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/matthieukhl/align-back/config"
	"github.com/matthieukhl/align-back/internal/handler"
	"github.com/matthieukhl/align-back/internal/repository"
	"github.com/matthieukhl/align-back/internal/service"
	"github.com/matthieukhl/align-back/pkg/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

func main() {
	cmd := &cobra.Command{
		Use:   "align-back",
		Short: "Pilates Management System API Server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/config.yaml)")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() {
	// Initialize config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.DB.ConnMaxLifetime) * time.Second)

	log.Info().Msg("Connected to database")

	// Initialize repositories
	clientRepo := repository.NewClientRepository(db)
	packageRepo := repository.NewPackageRepository(db)
	classRepo := repository.NewClassRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	billingRepo := repository.NewBillingRepository(db)

	// Initialize services
	clientService := service.NewClientService(clientRepo)
	packageService := service.NewPackageService(packageRepo)
	classService := service.NewClassService(classRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, classRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, scheduleRepo, clientRepo)
	billingService := service.NewBillingService(billingRepo, clientRepo, packageRepo)

	// Initialize handlers
	clientHandler := handler.NewClientHandler(clientService)
	packageHandler := handler.NewPackageHandler(packageService)
	classHandler := handler.NewClassHandler(classService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	billingHandler := handler.NewBillingHandler(billingService)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("API is running"))
		})

		// Clients endpoints
		r.Route("/clients", func(r chi.Router) {
			r.Get("/", clientHandler.GetAll)
			r.Post("/", clientHandler.Create)
			r.Get("/{id}", clientHandler.GetByID)
			r.Put("/{id}", clientHandler.Update)
			r.Delete("/{id}", clientHandler.Delete)
			r.Put("/low-group-credit", clientHandler.GetLowGroupCredits)
			r.Put("/low-private-credits", clientHandler.GetLowPrivateCredits)
		})

		// Packages endpoints
		r.Route("/packages", func(r chi.Router) {
			r.Get("/", packageHandler.GetAll)
			r.Post("/", packageHandler.Create)
			r.Get("/{id}", packageHandler.GetByID)
			r.Put("/{id}", packageHandler.Update)
			r.Delete("/{id}", packageHandler.Delete)
		})

		// Classes endpoints
		r.Route("/classes", func(r chi.Router) {
			r.Get("/", classHandler.GetAll)
			r.Post("/", classHandler.Create)
			r.Get("/{id}", classHandler.GetByID)
			r.Put("/{id}", classHandler.Update)
			r.Delete("/{id}", classHandler.Delete)
		})

		// Schedule endpoints
		r.Route("/schedule", func(r chi.Router) {
			r.Get("/", scheduleHandler.GetAll)
			r.Post("/", scheduleHandler.Create)
			r.Get("/{id}", scheduleHandler.GetByID)
			r.Put("/{id}", scheduleHandler.Update)
			r.Delete("/{id}", scheduleHandler.Delete)
			r.Get("/date/{date}", scheduleHandler.GetByDate)
			r.Get("/week/{date}", scheduleHandler.GetByWeek)
		})

		// Appointments endpoints
		r.Route("/appointments", func(r chi.Router) {
			r.Get("/", appointmentHandler.GetAll)
			r.Post("/", appointmentHandler.Create)
			r.Get("/{id}", appointmentHandler.GetByID)
			r.Put("/{id}", appointmentHandler.Update)
			r.Delete("/{id}", appointmentHandler.Delete)
			r.Get("/client/{clientId}", appointmentHandler.GetByClientID)
			r.Get("/schedule/{scheduleId}", appointmentHandler.GetByScheduleID)
		})

		// Billing endpoints
		r.Route("/billings", func(r chi.Router) {
			r.Get("/", billingHandler.GetAll)
			r.Post("/", billingHandler.Create)
			r.Get("/{id}", billingHandler.GetByID)
			r.Put("/{id}", billingHandler.Update)
			r.Delete("/{id}", billingHandler.Delete)
			r.Get("/client/{clientId}", billingHandler.GetByClientID)
			r.Get("/recent", billingHandler.GetRecent)
		})

		// Dashboard data endpoint
		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			// Implement dashboard data aggregation here
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Dashboard data endpoint"}`))
		})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Info().Msgf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
