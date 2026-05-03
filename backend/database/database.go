package database

import (
	"fmt"
	"log"
	"time"

	"documentManage/config"
	"documentManage/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.DocumentType{},
		&models.Document{},
		&models.Favorite{},
		&models.News{},
		&models.Banner{},
		&models.Announcement{},
	)
	if err != nil {
		log.Fatalf("自动迁移表结构失败: %v", err)
	}
	log.Println("表结构迁移完成")

	initSuperAdmin()
}

func initSuperAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)
	if count > 0 {
		return
	}

	superAdmin := &models.User{
		Username: "admin",
		Password: "admin123",
		RealName: "超级管理员",
		Role:     models.RoleSuperAdmin,
		Status:   1,
	}
	superAdmin.HashPassword()

	if err := DB.Create(superAdmin).Error; err != nil {
		log.Printf("创建超级管理员失败: %v", err)
		return
	}
	log.Println("超级管理员创建成功, 账号: admin, 密码: admin123")
}
