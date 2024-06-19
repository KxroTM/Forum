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
	Image       string
	NbPosts     int
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
		rows.Scan(&category.Category_id, &category.Name, &category.Description, &category.Users, &category.Image)
		category.NbPosts = GetNumberPostsByCategory(db, category.Category_id)
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return categories
}

func CreateCategory(db *sql.DB, name, description, image string) bool {

	if !isCategoryExists(db, name) {
		return false
	}

	u, err := uuid.NewV4()
	if err != nil {
		return false
	}
	_, err = db.Exec("INSERT INTO categories (category_id, name, description, users, image) VALUES (?, ?, ?, ?, ?)", u.String(), name, description, 0, image)
	return err == nil
}

func GetCategoryById(db *sql.DB, id string) Category {
	rows, err := db.Query("SELECT * FROM categories WHERE category_id = ?", id)
	if err != nil {
		return Category{}
	}
	defer rows.Close()

	var category Category

	for rows.Next() {
		rows.Scan(&category.Category_id, &category.Name, &category.Description, &category.Users, &category.Image)
	}
	if err := rows.Err(); err != nil {
		return Category{}
	}

	category.NbPosts = GetNumberPostsByCategory(db, id)

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

func DeleteCategory(db *sql.DB, id string) bool {
	_, err := db.Exec("DELETE FROM categories WHERE category_id = ?", id)
	return err == nil
}

func GetFirst5Categories(db *sql.DB) []Category {
	rows, err := db.Query("SELECT * FROM categories LIMIT 5")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		rows.Scan(&category.Category_id, &category.Name, &category.Description, &category.Users, &category.Image)
		category.NbPosts = GetNumberPostsByCategory(db, category.Category_id)

		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil
	}

	return categories
}

func GetNumberPostsByCategory(db *sql.DB, id string) int {
	query := "SELECT COUNT(*) FROM posts WHERE category_id = ?"
	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
