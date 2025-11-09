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
			users.Use(middleware.PermissionMiddleware("users.view"))
			{
				users.GET("", r.userHandler.GetAll)
				users.GET("/:id", r.userHandler.GetByID)
				
				users.POST("", middleware.PermissionMiddleware("users.create"), r.userHandler.Create)
				users.PUT("/:id", middleware.PermissionMiddleware("users.update"), r.userHandler.Update)
				users.DELETE("/:id", middleware.PermissionMiddleware("users.delete"), r.userHandler.Delete)
			}

			// Branch management endpoints
			branches := protected.Group("/branches")
			branches.Use(middleware.PermissionMiddleware("branches.view"))
			{
				branches.GET("", r.branchHandler.GetAll)
				branches.GET("/active", r.branchHandler.GetAllActive)
				branches.GET("/:id", r.branchHandler.GetByID)
				
				branches.POST("", middleware.PermissionMiddleware("branches.create"), r.branchHandler.Create)
				branches.PUT("/:id", middleware.PermissionMiddleware("branches.update"), r.branchHandler.Update)
				branches.DELETE("/:id", middleware.PermissionMiddleware("branches.delete"), r.branchHandler.Delete)
			}

			// Role management endpoints
			roles := protected.Group("/roles")
			roles.Use(middleware.PermissionMiddleware("roles.view"))
			{
				roles.GET("", r.roleHandler.GetAll)
				roles.GET("/:id", r.roleHandler.GetByID)
				
				roles.POST("", middleware.PermissionMiddleware("roles.create"), r.roleHandler.Create)
				roles.PUT("/:id", middleware.PermissionMiddleware("roles.update"), r.roleHandler.Update)
				roles.DELETE("/:id", middleware.PermissionMiddleware("roles.delete"), r.roleHandler.Delete)
			}

			// Account (COA) management endpoints
			accounts := protected.Group("/accounts")
			accounts.Use(middleware.PermissionMiddleware("accounts.view"))
			{
				accounts.GET("", r.accountHandler.GetAll)
				accounts.GET("/tree", r.accountHandler.GetTree)
				accounts.GET("/detail", r.accountHandler.GetDetailAccounts)
				accounts.GET("/category/:category", r.accountHandler.GetByCategory)
				accounts.GET("/code/:code", r.accountHandler.GetByCode)
				accounts.GET("/:id", r.accountHandler.GetByID)
				
				accounts.POST("", middleware.PermissionMiddleware("accounts.create"), r.accountHandler.Create)
				accounts.POST("/bulk-import", middleware.PermissionMiddleware("accounts.create"), r.accountHandler.BulkImport)
				accounts.PUT("/:id", middleware.PermissionMiddleware("accounts.update"), r.accountHandler.Update)
				accounts.DELETE("/:id", middleware.PermissionMiddleware("accounts.delete"), r.accountHandler.Delete)
			}

			// Journal Entry endpoints
			journals := protected.Group("/journals")
			journals.Use(middleware.PermissionMiddleware("journals.view"))
			{
				journals.GET("", r.journalHandler.GetAll)
				journals.GET("/status/:status", r.journalHandler.GetByStatus)
				journals.GET("/:id", r.journalHandler.GetByID)
				
				journals.POST("", middleware.PermissionMiddleware("journals.create"), r.journalHandler.Create)
				journals.PUT("/:id", middleware.PermissionMiddleware("journals.update"), r.journalHandler.Update)
				journals.DELETE("/:id", middleware.PermissionMiddleware("journals.delete"), r.journalHandler.Delete)
				
				journals.POST("/:id/submit", middleware.PermissionMiddleware("journals.submit"), r.journalHandler.SubmitForReview)
				journals.POST("/:id/review", middleware.PermissionMiddleware("journals.review"), r.journalHandler.Review)
				journals.POST("/:id/post", middleware.PermissionMiddleware("journals.post"), r.journalHandler.Post)
				journals.POST("/:id/unpost", middleware.PermissionMiddleware("journals.post"), r.journalHandler.Unpost)
			}

			// Budget endpoints
			budgets := protected.Group("/budgets")
			budgets.Use(middleware.PermissionMiddleware("budgets.view"))
			{
				budgets.GET("", r.budgetHandler.GetAll)
				budgets.GET("/:id", r.budgetHandler.GetByID)
				
				budgets.POST("", middleware.PermissionMiddleware("budgets.create"), r.budgetHandler.Create)
				budgets.PUT("/:id", middleware.PermissionMiddleware("budgets.update"), r.budgetHandler.Update)
				budgets.DELETE("/:id", middleware.PermissionMiddleware("budgets.delete"), r.budgetHandler.Delete)
				
				budgets.POST("/vs-actual", middleware.PermissionMiddleware("budgets.view"), r.budgetHandler.GetBudgetVsActual)
			}

			// Report endpoints
			reports := protected.Group("/reports")
			reports.Use(middleware.PermissionMiddleware("reports.view"))
			{
				reports.POST("/trial-balance", r.reportHandler.GetTrialBalance)
				reports.POST("/balance-sheet", r.reportHandler.GetBalanceSheet)
				reports.POST("/income-statement", r.reportHandler.GetIncomeStatement)
				reports.POST("/general-ledger", r.reportHandler.GetGeneralLedger)
			}

			// Student endpoints
			students := protected.Group("/students")
			students.Use(middleware.PermissionMiddleware("students.view"))
			{
				students.GET("", r.studentHandler.GetAll)
				students.GET("/search", r.studentHandler.Search)
				students.GET("/statistics", r.studentHandler.GetStatistics)
				students.GET("/:id", r.studentHandler.GetByID)
				
				students.POST("", middleware.PermissionMiddleware("students.create"), r.studentHandler.Create)
				students.PUT("/:id", middleware.PermissionMiddleware("students.update"), r.studentHandler.Update)
				students.DELETE("/:id", middleware.PermissionMiddleware("students.delete"), r.studentHandler.Delete)
			}

			// Payment endpoints
			payments := protected.Group("/payments")
			payments.Use(middleware.PermissionMiddleware("payments.view"))
			{
				payments.GET("", r.paymentHandler.GetAllPayments)
				payments.GET("/:id", r.paymentHandler.GetPaymentByID)
				
				payments.POST("", middleware.PermissionMiddleware("payments.create"), r.paymentHandler.CreatePayment)
				payments.POST("/:id/post", middleware.PermissionMiddleware("payments.post"), r.paymentHandler.PostPayment)
				payments.DELETE("/:id", middleware.PermissionMiddleware("payments.delete"), r.paymentHandler.DeletePayment)
			}

			// Invoice endpoints
			invoices := protected.Group("/invoices")
			invoices.Use(middleware.PermissionMiddleware("invoices.view"))
			{
				invoices.GET("", r.paymentHandler.GetAllInvoices)
				invoices.GET("/overdue", r.paymentHandler.GetOverdueInvoices)
				invoices.GET("/student/:student_id", r.paymentHandler.GetInvoicesByStudent)
				invoices.GET("/:id", r.paymentHandler.GetInvoiceByID)
				
				invoices.POST("", middleware.PermissionMiddleware("invoices.create"), r.paymentHandler.CreateInvoice)
				invoices.PUT("/:id", middleware.PermissionMiddleware("invoices.update"), r.paymentHandler.UpdateInvoice)
				invoices.DELETE("/:id", middleware.PermissionMiddleware("invoices.delete"), r.paymentHandler.DeleteInvoice)
			}

			// Employee endpoints
			employees := protected.Group("/employees")
			employees.Use(middleware.PermissionMiddleware("employees.view"))
			{
				employees.GET("", r.employeeHandler.GetAll)
				employees.GET("/search", r.employeeHandler.Search)
				employees.GET("/teachers", r.employeeHandler.GetTeachers)
				employees.GET("/statistics", r.employeeHandler.GetStatistics)
				employees.GET("/:id", r.employeeHandler.GetByID)
				employees.GET("/:employee_id/payrolls", r.employeeHandler.GetEmployeePayrolls)
				
				employees.POST("", middleware.PermissionMiddleware("employees.create"), r.employeeHandler.Create)
				employees.PUT("/:id", middleware.PermissionMiddleware("employees.update"), r.employeeHandler.Update)
				employees.DELETE("/:id", middleware.PermissionMiddleware("employees.delete"), r.employeeHandler.Delete)
			}

			// Payroll endpoints
			payrolls := protected.Group("/payrolls")
			payrolls.Use(middleware.PermissionMiddleware("payrolls.view"))
			{
				payrolls.GET("", r.employeeHandler.GetAllPayrolls)
				payrolls.GET("/period/:period", r.employeeHandler.GetPayrollsByPeriod)
				payrolls.GET("/:id", r.employeeHandler.GetPayrollByID)
				
				payrolls.POST("", middleware.PermissionMiddleware("payrolls.create"), r.employeeHandler.CreatePayroll)
				payrolls.POST("/bulk", middleware.PermissionMiddleware("payrolls.create"), r.employeeHandler.GenerateBulkPayroll)
				payrolls.POST("/:id/process", middleware.PermissionMiddleware("payrolls.process"), r.employeeHandler.ProcessPayroll)
				payrolls.DELETE("/:id", middleware.PermissionMiddleware("payrolls.delete"), r.employeeHandler.DeletePayroll)
			}

			// Asset endpoints
			assets := protected.Group("/assets")
			assets.Use(middleware.PermissionMiddleware("assets.view"))
			{
				assets.GET("", r.assetHandler.GetAll)
				assets.GET("/search", r.assetHandler.Search)
				assets.GET("/:id", r.assetHandler.GetByID)
				assets.GET("/:id/depreciation", r.assetHandler.CalculateDepreciation)
				
				assets.POST("", middleware.PermissionMiddleware("assets.create"), r.assetHandler.Create)
				assets.PUT("/:id", middleware.PermissionMiddleware("assets.update"), r.assetHandler.Update)
				assets.DELETE("/:id", middleware.PermissionMiddleware("assets.delete"), r.assetHandler.Delete)
				
				// Maintenance
				assets.POST("/maintenance", middleware.PermissionMiddleware("assets.create"), r.assetHandler.CreateMaintenance)
				
				// Transfer
				assets.POST("/transfer", middleware.PermissionMiddleware("assets.create"), r.assetHandler.CreateTransfer)
				assets.POST("/transfer/:id/approve", middleware.PermissionMiddleware("assets.approve"), r.assetHandler.ApproveTransfer)
			}

			// Inventory endpoints
			inventory := protected.Group("/inventory")
			inventory.Use(middleware.PermissionMiddleware("inventory.view"))
			{
				// Items
				inventory.GET("/items", r.inventoryHandler.GetAllItems)
				inventory.GET("/items/low-stock", r.inventoryHandler.GetLowStockItems)
				inventory.GET("/items/search", r.inventoryHandler.SearchItems)
				inventory.GET("/items/:id", r.inventoryHandler.GetItemByID)
				inventory.POST("/items", middleware.PermissionMiddleware("inventory.create"), r.inventoryHandler.CreateItem)
				inventory.PUT("/items/:id", middleware.PermissionMiddleware("inventory.update"), r.inventoryHandler.UpdateItem)
				inventory.DELETE("/items/:id", middleware.PermissionMiddleware("inventory.delete"), r.inventoryHandler.DeleteItem)
				
				// Transactions
				inventory.POST("/stock-in", middleware.PermissionMiddleware("inventory.create"), r.inventoryHandler.CreateStockIn)
				inventory.POST("/stock-out", middleware.PermissionMiddleware("inventory.create"), r.inventoryHandler.CreateStockOut)
				inventory.POST("/adjustment", middleware.PermissionMiddleware("inventory.create"), r.inventoryHandler.CreateAdjustment)
				inventory.GET("/items/:item_id/transactions", r.inventoryHandler.GetTransactionHistory)
				
				// Stock Opname
				inventory.POST("/opname", middleware.PermissionMiddleware("inventory.create"), r.inventoryHandler.CreateStockOpname)
				inventory.POST("/opname/:id/approve", middleware.PermissionMiddleware("inventory.approve"), r.inventoryHandler.ApproveStockOpname)
				inventory.POST("/opname/:id/process", middleware.PermissionMiddleware("inventory.approve"), r.inventoryHandler.ProcessStockOpname)
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
