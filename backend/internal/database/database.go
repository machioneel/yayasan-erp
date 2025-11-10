package database

import (
	"fmt"
	"log"
	"time"

	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes database connection
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	var gormLogger logger.Interface
	if config.GlobalConfig.IsDevelopment() {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established successfully")

	DB = db
	return db, nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// InitDB initializes database connection and runs migrations
func InitDB() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	config.GlobalConfig = cfg
	
	db, err := Connect(&cfg.Database)
	if err != nil {
		return err
	}
	
	DB = db
	
	// Run auto migrations
	log.Println("🔄 Running database migrations...")
	if err := AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("✅ Database migrations completed")
	
	return nil
}

// CloseDB closes database connection
func CloseDB() error {
	log.Println("Closing database connection...")
	return Close()
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Branch{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.Account{},
		&models.Journal{},
		&models.JournalLine{},
		&models.Budget{},
		&models.FiscalYear{},
		&models.Student{},
		&models.Parent{},
		&models.StudentParent{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.Payment{},
		&models.Employee{},
		&models.Payroll{},
		&models.Attendance{},
		&models.Asset{},
		&models.InventoryItem{},
		&models.InventoryCategory{},
		&models.StockTransaction{},
		&models.StockOpname{},
		&models.StockOpnameItem{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.Setting{},
		&models.AuditLog{},
		&models.Donor{},
	)
}
