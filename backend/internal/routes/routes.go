package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yayasan/erp-backend/internal/handler"
	"github.com/yayasan/erp-backend/internal/middleware"
)

type Router struct {
	authHandler      *handler.AuthHandler
	userHandler      *handler.UserHandler
	branchHandler    *handler.BranchHandler
	roleHandler      *handler.RoleHandler
	accountHandler   *handler.AccountHandler
	journalHandler   *handler.JournalHandler
	budgetHandler    *handler.BudgetHandler
	reportHandler    *handler.ReportHandler
	studentHandler   *handler.StudentHandler
	paymentHandler   *handler.PaymentHandler
	employeeHandler  *handler.EmployeeHandler
	assetHandler     *handler.AssetHandler
	inventoryHandler *handler.InventoryHandler
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	branchHandler *handler.BranchHandler,
	roleHandler *handler.RoleHandler,
	accountHandler *handler.AccountHandler,
	journalHandler *handler.JournalHandler,
	budgetHandler *handler.BudgetHandler,
	reportHandler *handler.ReportHandler,
	studentHandler *handler.StudentHandler,
	paymentHandler *handler.PaymentHandler,
	employeeHandler *handler.EmployeeHandler,
	assetHandler *handler.AssetHandler,
	inventoryHandler *handler.InventoryHandler,
) *Router {
	return &Router{
		authHandler:      authHandler,
		userHandler:      userHandler,
		branchHandler:    branchHandler,
		roleHandler:      roleHandler,
		accountHandler:   accountHandler,
		journalHandler:   journalHandler,
		budgetHandler:    budgetHandler,
		reportHandler:    reportHandler,
		studentHandler:   studentHandler,
		paymentHandler:   paymentHandler,
		employeeHandler:  employeeHandler,
		assetHandler:     assetHandler,
		inventoryHandler: inventoryHandler,
	}
}

func (r *Router) Setup(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Yayasan As-Salam Joglo ERP API is running",
			"version": "1.0.0",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth endpoints (authenticated)
			authProtected := protected.Group("/auth")
			{
				authProtected.GET("/me", r.authHandler.GetMe)
				authProtected.POST("/change-password", r.authHandler.ChangePassword)
				authProtected.POST("/logout", r.authHandler.Logout)
			}

			// User profile endpoints
			profile := protected.Group("/profile")
			{
				profile.PUT("", r.userHandler.UpdateProfile)
			}

			// User management endpoints (requires permissions)
			users := protected.Group("/users")
			users.Use(middleware.RequirePermission("users.view")) // DIPERBAIKI
			{
				users.GET("", r.userHandler.GetAll)
				users.GET("/:id", r.userHandler.GetByID)

				users.POST("", middleware.RequirePermission("users.create"), r.userHandler.Create)     // DIPERBAIKI
				users.PUT("/:id", middleware.RequirePermission("users.update"), r.userHandler.Update)   // DIPERBAIKI
				users.DELETE("/:id", middleware.RequirePermission("users.delete"), r.userHandler.Delete) // DIPERBAIKI
			}

			// Branch management endpoints
			branches := protected.Group("/branches")
			branches.Use(middleware.RequirePermission("branches.view")) // DIPERBAIKI
			{
				branches.GET("", r.branchHandler.GetAll)
				branches.GET("/active", r.branchHandler.GetAllActive)
				branches.GET("/:id", r.branchHandler.GetByID)

				branches.POST("", middleware.RequirePermission("branches.create"), r.branchHandler.Create) // DIPERBAIKI
				branches.PUT("/:id", middleware.RequirePermission("branches.update"), r.branchHandler.Update) // DIPERBAIKI
				branches.DELETE("/:id", middleware.RequirePermission("branches.delete"), r.branchHandler.Delete) // DIPERBAIKI
			}

			// Role management endpoints
			roles := protected.Group("/roles")
			roles.Use(middleware.RequirePermission("roles.view")) // DIPERBAIKI
			{
				roles.GET("", r.roleHandler.GetAll)
				roles.GET("/:id", r.roleHandler.GetByID)

				roles.POST("", middleware.RequirePermission("roles.create"), r.roleHandler.Create) // DIPERBAIKI
				roles.PUT("/:id", middleware.RequirePermission("roles.update"), r.roleHandler.Update) // DIPERBAIKI
				roles.DELETE("/:id", middleware.RequirePermission("roles.delete"), r.roleHandler.Delete) // DIPERBAIKI
			}

			// Account (COA) management endpoints
			accounts := protected.Group("/accounts")
			accounts.Use(middleware.RequirePermission("accounts.view")) // DIPERBAIKI
			{
				accounts.GET("", r.accountHandler.GetAll)
				accounts.GET("/tree", r.accountHandler.GetTree)
				accounts.GET("/detail", r.accountHandler.GetDetailAccounts)
				accounts.GET("/category/:category", r.accountHandler.GetByCategory)
				accounts.GET("/code/:code", r.accountHandler.GetByCode)
				accounts.GET("/:id", r.accountHandler.GetByID)

				accounts.POST("", middleware.RequirePermission("accounts.create"), r.accountHandler.Create)       // DIPERBAIKI
				accounts.POST("/bulk-import", middleware.RequirePermission("accounts.create"), r.accountHandler.BulkImport) // DIPERBAIKI
				accounts.PUT("/:id", middleware.RequirePermission("accounts.update"), r.accountHandler.Update)     // DIPERBAIKI
				accounts.DELETE("/:id", middleware.RequirePermission("accounts.delete"), r.accountHandler.Delete) // DIPERBAIKI
			}

			// Journal Entry endpoints
			journals := protected.Group("/journals")
			journals.Use(middleware.RequirePermission("journals.view")) // DIPERBAIKI
			{
				journals.GET("", r.journalHandler.GetAll)
				journals.GET("/status/:status", r.journalHandler.GetByStatus)
				journals.GET("/:id", r.journalHandler.GetByID)

				journals.POST("", middleware.RequirePermission("journals.create"), r.journalHandler.Create) // DIPERBAIKI
				journals.PUT("/:id", middleware.RequirePermission("journals.update"), r.journalHandler.Update) // DIPERBAIKI
				journals.DELETE("/:id", middleware.RequirePermission("journals.delete"), r.journalHandler.Delete) // DIPERBAIKI

				journals.POST("/:id/submit", middleware.RequirePermission("journals.submit"), r.journalHandler.SubmitForReview) // DIPERBAIKI
				journals.POST("/:id/review", middleware.RequirePermission("journals.review"), r.journalHandler.Review) // DIPERBAIKI
				journals.POST("/:id/post", middleware.RequirePermission("journals.post"), r.journalHandler.Post)     // DIPERBAIKI
				journals.POST("/:id/unpost", middleware.RequirePermission("journals.post"), r.journalHandler.Unpost) // DIPERBAIKI
			}

			// Budget endpoints
			budgets := protected.Group("/budgets")
			budgets.Use(middleware.RequirePermission("budgets.view")) // DIPERBAIKI
			{
				budgets.GET("", r.budgetHandler.GetAll)
				budgets.GET("/:id", r.budgetHandler.GetByID)

				budgets.POST("", middleware.RequirePermission("budgets.create"), r.budgetHandler.Create) // DIPERBAIKI
				budgets.PUT("/:id", middleware.RequirePermission("budgets.update"), r.budgetHandler.Update) // DIPERBAIKI
				budgets.DELETE("/:id", middleware.RequirePermission("budgets.delete"), r.budgetHandler.Delete) // DIPERBAIKI

				budgets.POST("/vs-actual", middleware.RequirePermission("budgets.view"), r.budgetHandler.GetBudgetVsActual) // DIPERBAIKI
			}

			// Report endpoints
			reports := protected.Group("/reports")
			reports.Use(middleware.RequirePermission("reports.view")) // DIPERBAIKI
			{
				reports.POST("/trial-balance", r.reportHandler.GetTrialBalance)
				reports.POST("/balance-sheet", r.reportHandler.GetBalanceSheet)
				reports.POST("/income-statement", r.reportHandler.GetIncomeStatement)
				reports.POST("/general-ledger", r.reportHandler.GetGeneralLedger)
			}

			// Student endpoints
			students := protected.Group("/students")
			students.Use(middleware.RequirePermission("students.view")) // DIPERBAIKI
			{
				students.GET("", r.studentHandler.GetAll)
				students.GET("/search", r.studentHandler.Search)
				students.GET("/statistics", r.studentHandler.GetStatistics)
				students.GET("/:id", r.studentHandler.GetByID)

				students.POST("", middleware.RequirePermission("students.create"), r.studentHandler.Create) // DIPERBAIKI
				students.PUT("/:id", middleware.RequirePermission("students.update"), r.studentHandler.Update) // DIPERBAIKI
				students.DELETE("/:id", middleware.RequirePermission("students.delete"), r.studentHandler.Delete) // DIPERBAIKI
			}

			// Payment endpoints
			payments := protected.Group("/payments")
			payments.Use(middleware.RequirePermission("payments.view")) // DIPERBAIKI
			{
				payments.GET("", r.paymentHandler.GetAllPayments)
				payments.GET("/:id", r.paymentHandler.GetPaymentByID)

				payments.POST("", middleware.RequirePermission("payments.create"), r.paymentHandler.CreatePayment) // DIPERBAIKI
				payments.POST("/:id/post", middleware.RequirePermission("payments.post"), r.paymentHandler.PostPayment) // DIPERBAIKI
				payments.DELETE("/:id", middleware.RequirePermission("payments.delete"), r.paymentHandler.DeletePayment) // DIPERBAIKI
			}

			// Invoice endpoints
			invoices := protected.Group("/invoices")
			invoices.Use(middleware.RequirePermission("invoices.view")) // DIPERBAIKI
			{
				invoices.GET("", r.paymentHandler.GetAllInvoices)
				invoices.GET("/overdue", r.paymentHandler.GetOverdueInvoices)
				invoices.GET("/student/:student_id", r.paymentHandler.GetInvoicesByStudent)
				invoices.GET("/:id", r.paymentHandler.GetInvoiceByID)

				invoices.POST("", middleware.RequirePermission("invoices.create"), r.paymentHandler.CreateInvoice) // DIPERBAIKI
				invoices.PUT("/:id", middleware.RequirePermission("invoices.update"), r.paymentHandler.UpdateInvoice) // DIPERBAIKI
				invoices.DELETE("/:id", middleware.RequirePermission("invoices.delete"), r.paymentHandler.DeleteInvoice) // DIPERBAIKI
			}

			// Employee endpoints
			employees := protected.Group("/employees")
			employees.Use(middleware.RequirePermission("employees.view")) // DIPERBAIKI
			{
				employees.GET("", r.employeeHandler.GetAll)
				employees.GET("/search", r.employeeHandler.Search)
				employees.GET("/teachers", r.employeeHandler.GetTeachers)
				employees.GET("/statistics", r.employeeHandler.GetStatistics)
				employees.GET("/:id", r.employeeHandler.GetByID)
				employees.GET("/:employee_id/payrolls", r.employeeHandler.GetEmployeePayrolls)

				employees.POST("", middleware.RequirePermission("employees.create"), r.employeeHandler.Create) // DIPERBAIKI
				employees.PUT("/:id", middleware.RequirePermission("employees.update"), r.employeeHandler.Update) // DIPERBAIKI
				employees.DELETE("/:id", middleware.RequirePermission("employees.delete"), r.employeeHandler.Delete) // DIPERBAIKI
			}

			// Payroll endpoints
			payrolls := protected.Group("/payrolls")
			payrolls.Use(middleware.RequirePermission("payrolls.view")) // DIPERBAIKI
			{
				payrolls.GET("", r.employeeHandler.GetAllPayrolls)
				payrolls.GET("/period/:period", r.employeeHandler.GetPayrollsByPeriod)
				payrolls.GET("/:id", r.employeeHandler.GetPayrollByID)

				payrolls.POST("", middleware.RequirePermission("payrolls.create"), r.employeeHandler.CreatePayroll)       // DIPERBAIKI
				payrolls.POST("/bulk", middleware.RequirePermission("payrolls.create"), r.employeeHandler.GenerateBulkPayroll) // DIPERBAIKI
				payrolls.POST("/:id/process", middleware.RequirePermission("payrolls.process"), r.employeeHandler.ProcessPayroll) // DIPERBAIKI
				payrolls.DELETE("/:id", middleware.RequirePermission("payrolls.delete"), r.employeeHandler.DeletePayroll) // DIPERBAIKI
			}

			// Asset endpoints
			assets := protected.Group("/assets")
			assets.Use(middleware.RequirePermission("assets.view")) // DIPERBAIKI
			{
				assets.GET("", r.assetHandler.GetAll)
				assets.GET("/search", r.assetHandler.Search)
				assets.GET("/:id", r.assetHandler.GetByID)
				assets.GET("/:id/depreciation", r.assetHandler.CalculateDepreciation)

				assets.POST("", middleware.RequirePermission("assets.create"), r.assetHandler.Create) // DIPERBAIKI
				assets.PUT("/:id", middleware.RequirePermission("assets.update"), r.assetHandler.Update) // DIPERBAIKI
				assets.DELETE("/:id", middleware.RequirePermission("assets.delete"), r.assetHandler.Delete) // DIPERBAIKI

				// Maintenance
				assets.POST("/maintenance", middleware.RequirePermission("assets.create"), r.assetHandler.CreateMaintenance) // DIPERBAIKI

				// Transfer
				assets.POST("/transfer", middleware.RequirePermission("assets.create"), r.assetHandler.CreateTransfer) // DIPERBAIKI
				assets.POST("/transfer/:id/approve", middleware.RequirePermission("assets.approve"), r.assetHandler.ApproveTransfer) // DIPERBAIKI
			}

			// Inventory endpoints
			inventory := protected.Group("/inventory")
			inventory.Use(middleware.RequirePermission("inventory.view")) // DIPERBAIKI
			{
				// Items
				inventory.GET("/items", r.inventoryHandler.GetAllItems)
				inventory.GET("/items/low-stock", r.inventoryHandler.GetLowStockItems)
				inventory.GET("/items/search", r.inventoryHandler.SearchItems)
				inventory.GET("/items/:id", r.inventoryHandler.GetItemByID)
				inventory.POST("/items", middleware.RequirePermission("inventory.create"), r.inventoryHandler.CreateItem) // DIPERBAIKI
				inventory.PUT("/items/:id", middleware.RequirePermission("inventory.update"), r.inventoryHandler.UpdateItem) // DIPERBAIKI
				inventory.DELETE("/items/:id", middleware.RequirePermission("inventory.delete"), r.inventoryHandler.DeleteItem) // DIPERBAIKI

				// Transactions
				inventory.POST("/stock-in", middleware.RequirePermission("inventory.create"), r.inventoryHandler.CreateStockIn) // DIPERBAIKI
				inventory.POST("/stock-out", middleware.RequirePermission("inventory.create"), r.inventoryHandler.CreateStockOut) // DIPERBAIKI
				inventory.POST("/adjustment", middleware.RequirePermission("inventory.create"), r.inventoryHandler.CreateAdjustment) // DIPERBAIKI
				inventory.GET("/items/:item_id/transactions", r.inventoryHandler.GetTransactionHistory)

				// Stock Opname
				inventory.POST("/opname", middleware.RequirePermission("inventory.create"), r.inventoryHandler.CreateStockOpname) // DIPERBAIKI
				inventory.POST("/opname/:id/approve", middleware.RequirePermission("inventory.approve"), r.inventoryHandler.ApproveStockOpname) // DIPERBAIKI
				inventory.POST("/opname/:id/process", middleware.RequirePermission("inventory.approve"), r.inventoryHandler.ProcessStockOpname) // DIPERBAIKI
			}
		}
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Route not found",
			"error":   "The requested endpoint does not exist",
		})
	})
}