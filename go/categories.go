package forum

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

type Category struct {
	Category_id string
	Name        string
	Description string
	Users       int
}

func GetAllCategories(db *sql.DB) []Category {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		rows.Scan(&category.Category_id, &category.Name, &category.Description, &category.Users)
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return categories
}

func CreateCategory(db *sql.DB, name string, description string) bool {

	if !isCategoryExists(db, name) {
		return false
	}

	u, err := uuid.NewV4()
	if err != nil {
		return false
	}
	_, err = db.Exec("INSERT INTO categories (category_id, name, description, users) VALUES (?, ?, ?, ?)", u.String(), name, description, 0)
	return err == nil
}

func GetCategoryById(db *sql.DB, id string) Category {
	var category Category
	err := db.QueryRow("SELECT * FROM categories WHERE category_id = ?", id).Scan(&category.Category_id, &category.Name, &category.Description, &category.Users)
	if err != nil {
		return Category{}
	}
	return category
}

func isCategoryExists(db *sql.DB, name string) bool {
	query := "SELECT COUNT(*) FROM categories WHERE name = ?"
	var count int
	err := db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}
