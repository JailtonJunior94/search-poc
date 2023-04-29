package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jailtonjunior94/search-poc/pkg/responses"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       int `gorm:"primaryKey"`
	FullName string
	Email    string
	gorm.Model
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=search_poc port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})

	//users := []User{
	// 	{FullName: "Jailton Angelo Teixeira Junior", Email: "jailton.junior94@outlook.com"},
	// 	{FullName: "Stefany Kelly Lima Teixeira", Email: "stefany.teixeira@outlook.com"},
	// 	{FullName: "Antony Lima Teixeira", Email: "antony.teixeira@outlook.com"},
	// }
	// db.Create(&users)

	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))
	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		var users []User

		db.Where("LOWER(full_name) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(name))).Find(&users)
		responses.JSON(w, http.StatusOK, users)
	})

	server := http.Server{
		ReadTimeout:       time.Duration(30) * time.Second,
		ReadHeaderTimeout: time.Duration(30) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "8000"))
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
