package db

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/weeb-vip/auth/config"
)

const (
	maxIdleConns = 10
	maxOpenConns = 100
)

type SafeDBService struct {
	mu sync.Mutex
	db DB //nolint
}

type dbConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  bool
}

type Service struct {
	db *gorm.DB
}

var dbservice = SafeDBService{ // nolint
	db: nil,
}

func (service *Service) setupSQLDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func (service *Service) connect(config dbConfig) *gorm.DB {
	sslmode := "disable"
	if config.sslmode {
		sslmode = "enable"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.host,
			config.port,
			config.user,
			config.password,
			config.dbname,
			sslmode,
		),
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	service.setupSQLDB(db)

	service.db = db

	return db
}

func NewDBService() DB { //nolint
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	temp := &Service{}
	temp.connect(dbConfig{
		host:     cfg.DBConfig.Host,
		port:     cfg.DBConfig.Port,
		user:     cfg.DBConfig.User,
		password: cfg.DBConfig.Password,
		dbname:   cfg.DBConfig.DB,
		sslmode:  cfg.DBConfig.SSL,
	})

	dbservice.SetDB(temp)

	return dbservice.GetDB()
}

func (service *Service) GetDB() *gorm.DB {
	return service.db
}

func (service *SafeDBService) GetDB() DB {
	service.mu.Lock()
	defer service.mu.Unlock()

	return service.db
}

func (service *SafeDBService) SetDB(db DB) {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.db = db
}

func GetDBService() DB {
	if dbservice.GetDB() != nil {
		return dbservice.GetDB()
	}

	return NewDBService()
}
