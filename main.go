package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Port       string `json:"port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
}

type Expense struct {
	ID         int       `db:"id" json:"id,omitempty"`
	Amount     float64   `db:"amount" json:"amount"`
	CategoryID int       `db:"category_id" json:"category_id"`
	Category   string    `db:"category" json:"category,omitempty"`
	Comment    string    `db:"comment" json:"comment"`
	Date       time.Time `db:"date" json:"date,omitempty"`
}

type Category struct {
	ID      int    `db:"id" json:"id,omitempty"`
	Name    string `db:"name" json:"name"`
	Deleted bool   `db:"deleted" json:"deleted,omitempty"`
}

var db *sqlx.DB

func main() {
	config := loadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/expenses", createExpense).Methods("POST")
	r.HandleFunc("/expenses/{id}", getExpenseByID).Methods("GET")
	r.HandleFunc("/expenses", getExpenses).Methods("GET")
	r.HandleFunc("/expenses/{id}", editExpense).Methods("PUT")
	r.HandleFunc("/expenses/{id}", deleteExpense).Methods("DELETE")

	r.HandleFunc("/categories", createCategory).Methods("POST")
	r.HandleFunc("/categories", getNotDeletedCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", getCategoryByID).Methods("GET")
	r.HandleFunc("/categories/{id}", editCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
}

func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Failed to open config file, using default values")
		return Config{
			Port:       "8080",
			DBUser:     "user",
			DBPassword: "password",
			DBHost:     "localhost",
			DBPort:     "3306",
			DBName:     "expenses_db",
		}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Failed to decode config file, using default values")
		return Config{
			Port:       "8080",
			DBUser:     "user",
			DBPassword: "password",
			DBHost:     "localhost",
			DBPort:     "3306",
			DBName:     "expenses_db",
		}
	}

	return config
}

func createExpense(w http.ResponseWriter, r *http.Request) {
	var e Expense
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %v", err), http.StatusBadRequest)
		return
	}
	result, err := db.NamedExec("INSERT INTO expenses (amount, category_id, comment, date) VALUES (:amount, :category_id, :comment, NOW())", &e)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert expense: %v", err), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	e.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func getExpenseByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var e Expense
	err := db.Get(&e, "SELECT id, amount, category_id, comment, date FROM expenses WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve expense: %v", err), http.StatusInternalServerError)
		return
	}
	e.Date = e.Date.UTC()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var c Category
	err := db.Get(&c, "SELECT id, name, deleted FROM categories WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve category: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func getNotDeletedCategories(w http.ResponseWriter, r *http.Request) {
	categories := []Category{}
	err := db.Select(&categories, "SELECT id, name FROM categories WHERE deleted = false")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve categories: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getExpenses(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("size")

	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "25"
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	offset := pageInt * pageSizeInt

	expenses := []Expense{}
	err := db.Select(&expenses, "SELECT e.id, e.amount, e.category_id, IFNULL(c.name, '') AS category, e.comment, e.date FROM expenses e LEFT JOIN categories c ON e.category_id = c.id ORDER BY e.date DESC LIMIT ? OFFSET ?", pageSizeInt, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve expenses: %v", err), http.StatusInternalServerError)
		return
	}
	for i := range expenses {
		expenses[i].Date = expenses[i].Date.UTC()
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve expenses: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func editExpense(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var e Expense
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = db.NamedExec("UPDATE expenses SET amount=:amount, category_id=:category_id WHERE id=:id", map[string]interface{}{
		"amount":      e.Amount,
		"category_id": e.CategoryID,
		"id":          id,
	})
	if err != nil {
		http.Error(w, "Failed to update expense", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func deleteExpense(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM expenses WHERE id=?", id)
	if err != nil {
		http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var c Category
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = db.NamedExec("INSERT INTO categories (name) VALUES (:name)", &c)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create category: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func editCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var c Category
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = db.NamedExec("UPDATE categories SET name=:name WHERE id=:id", map[string]interface{}{
		"name": c.Name,
		"id":   id,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update category: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := db.Exec("UPDATE categories SET deleted=true WHERE id=?", id)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
