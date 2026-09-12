package connect

import (
	"errors"
	"fmt"
	"os"

	"github.com/adhir2099/GO-GOLANG-API/structures"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

const engineSQL = "mysql"

var ErrNotFound = errors.New("user not found")

func InitializeDB() error {
	connection, err := ConnectORM(DBString())
	if err != nil {
		return err
	}
	db = connection
	return nil
}

func CloseConnection() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func ConnectORM(connectionString string) (*gorm.DB, error) {
	connection, err := gorm.Open(engineSQL, connectionString)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}
	return connection, nil
}

func GetUser(id uint) (structures.User, error) {
	var user structures.User
	if err := db.First(&user, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return structures.User{}, ErrNotFound
		}
		return structures.User{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func CreateUser(user structures.User) (structures.User, error) {
	if err := db.Create(&user).Error; err != nil {
		return structures.User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func UpdateUser(id uint, updates map[string]any) (structures.User, error) {
	user, err := GetUser(id)
	if err != nil {
		return structures.User{}, err
	}
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return structures.User{}, fmt.Errorf("update user: %w", err)
	}
	return GetUser(id)
}

func DeleteUser(id uint) error {
	user, err := GetUser(id)
	if err != nil {
		return err
	}
	if err := db.Delete(&user).Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func DBString() string {
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		return dsn
	}
	return (&mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      envOrDefault("DB_HOST", "127.0.0.1:3306"),
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}).FormatDSN()
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
