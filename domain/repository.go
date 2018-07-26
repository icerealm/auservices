package domain

import "database/sql"

//Repository represent proxy to find record in database.
type Repository struct {
	dbConn *sql.DB
}

//FindCategoryByName find category by name
func (r *Repository) FindCategoryByName(cName string, userID string) (*Category, error) {
	query := `
	SELECT id, category_nm, category_desc, category_type, user_id 
	FROM CATEGORY
	WHERE category_nm = $1 AND user_id = $2;
	`
	var category Category
	err := r.dbConn.QueryRow(query, cName, userID).Scan(&category.id,
		&category.categoryNm,
		&category.categoryDesc,
		&category.categoryType,
		&category.userID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

//FindAllCategories find all categories of each user id
func (r *Repository) FindAllCategories(userID string) ([]*Category, error) {
	query := `SELECT id, category_nm, category_desc, category_type, user_id 
	FROM category
	WHERE user_id = $1`

	rows, err := r.dbConn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resultList []*Category
	for rows.Next() {
		var category Category
		rows.Scan(&category.id,
			&category.categoryNm,
			&category.categoryDesc,
			&category.categoryType,
			&category.userID)
		resultList = append(resultList, &category)
	}
	return resultList, nil
}
