package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/database"
	"github.com/yayasan/erp-backend/internal/handler"
	"github.com/yayasan/erp-backend/internal/middleware"
	"github.com/yayasan/erp-backend/internal/repository"
	"github.com/yayasan/erp-backend/internal/routes"
	"github.com/yayasan/erp-backend/internal/service"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Check if seed command
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		log.Println("Running database seeder...")
		if err := runSeeder(); err != nil {
			log.Fatal("Failed to seed database:", err)
		}
		log.Println("✅ Database seeded successfully!")
		return
	}

	// Initialize Gin
	if config.GlobalConfig.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	branchRepo := repository.NewBranchRepository(database.DB)
	roleRepo := repository.NewRoleRepository(database.DB)
	accountRepo := repository.NewAccountRepository(database.DB)
	journalRepo := repository.NewJournalRepository(database.DB)
	budgetRepo := repository.NewBudgetRepository(database.DB)
	fiscalYearRepo := repository.NewFiscalYearRepository(database.DB)
	studentRepo := repository.NewStudentRepository(database.DB)
	parentRepo := repository.NewParentRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)
	invoiceRepo := repository.NewInvoiceRepository(database.DB)
	employeeRepo := repository.NewEmployeeRepository(database.DB)
	payrollRepo := repository.NewPayrollRepository(database.DB)
	assetRepo := repository.NewAssetRepository(database.DB)
	inventoryRepo := repository.NewInventoryRepository(database.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo)
	userService := service.NewUserService(userRepo, roleRepo)
	branchService := service.NewBranchService(branchRepo)
	roleService := service.NewRoleService(roleRepo)
	accountService := service.NewAccountService(accountRepo)
	journalService := service.NewJournalService(journalRepo, accountRepo, branchRepo)
	budgetService := service.NewBudgetService(database.DB, budgetRepo, accountRepo, fiscalYearRepo)
	reportService := service.NewReportService(database.DB, accountRepo, journalRepo)
	studentService := service.NewStudentService(studentRepo, parentRepo, branchRepo)
	paymentService := service.NewPaymentService(paymentRepo, invoiceRepo, branchRepo, studentRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo, studentRepo, branchRepo)
	employeeService := service.NewEmployeeService(employeeRepo, branchRepo)
	payrollService := service.NewPayrollService(payrollRepo, employeeRepo, branchRepo)
	assetService := service.NewAssetService(assetRepo, branchRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, branchRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	branchHandler := handler.NewBranchHandler(branchService)
	roleHandler := handler.NewRoleHandler(roleService)
	accountHandler := handler.NewAccountHandler(accountService)
	journalHandler := handler.NewJournalHandler(journalService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	reportHandler := handler.NewReportHandler(reportService)
	studentHandler := handler.NewStudentHandler(studentService)
	paymentHandler := handler.NewPaymentHandler(paymentService, invoiceService)
	employeeHandler := handler.NewEmployeeHandler(employeeService, payrollService)
	assetHandler := handler.NewAssetHandler(assetService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	// Setup routes
	appRouter := routes.NewRouter(
		authHandler,
		userHandler,
		branchHandler,
		roleHandler,
		accountHandler,
		journalHandler,
		budgetHandler,
		reportHandler,
		studentHandler,
		paymentHandler,
		employeeHandler,
		assetHandler,
		inventoryHandler,
	)
	appRouter.Setup(router)

	// Create server
	addr := fmt.Sprintf("%s:%s", config.GlobalConfig.Server.Host, config.GlobalConfig.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.GlobalConfig.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.GlobalConfig.Server.IdleTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on %s", addr)
		log.Printf("📝 Environment: %s", config.GlobalConfig.App.Env)
		log.Printf("🌐 API Base URL: http://%s/api/v1", addr)
		log.Printf("❤️  Health Check: http://%s/health", addr)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited gracefully")
}

func runSeeder() error {
	// TODO: Implement database seeder
	// This will be created in the next step
	log.Println("⚠️  Database seeder not yet implemented")
	log.Println("Please run database/seeds/seeder.go manually")
	return nil
}
