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
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Branch struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Name       string
	CategoryID uuid.NullUUID
	Category   *Category
}

func NewBranch(name string, categoryID uuid.NullUUID) *Branch {
	return &Branch{
		ID:         uuid.New(),
		Name:       name,
		CategoryID: categoryID,
	}
}

type Category struct {
	ID   uuid.UUID `gorm:"primaryKey"`
	Name string
}

func NewCategory(name string) *Category {
	return &Category{
		ID:   uuid.New(),
		Name: name,
	}
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=search_poc port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var branches []*Branch
	db.Model(&Branch{}).Preload("Category").Find(&branches)

	// // category := NewCategory("categoria1")
	// // db.Create(category)

	// branch := NewBranch("branch", uuid.NullUUID{})
	// db.Create(branch)

	//users := []User{
	// 	{FullName: "Jailton Angelo Teixeira Junior", Email: "jailton.junior94@outlook.com"},
	// 	{FullName: "Stefany Kelly Lima Teixeira", Email: "stefany.teixeira@outlook.com"},
	// 	{FullName: "Antony Lima Teixeira", Email: "antony.teixeira@outlook.com"},
	// }
	// db.Create(&users)

	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))
	router.Get("/branches", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		var users []Branch

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
