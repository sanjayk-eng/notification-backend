package query

import (
	"database/sql"
	"fmt"
	"log"
	"sanjay/api/util"
)

type QueryMeth struct {
	DB *sql.DB
}

func GetQueryHandler(db *sql.DB) *QueryMeth {
	return &QueryMeth{
		DB: db,
	}
}
func (q *QueryMeth) CheckUserWithCountryCode(phone string, countryCode string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1 AND country_code=$2)`

	err := q.DB.QueryRow(query, phone, countryCode).Scan(&exists)
	if err != nil {
		log.Printf("CheckUserWithCountryCode error: %v", err)
		return false
	}
	return exists
}

func (q *QueryMeth) StoredPhoneWithCountryCode(countryCode, phone string) error {
	fmt.Println("country-code", countryCode, "phone", phone)
	if exists := q.CheckUserWithCountryCode(phone, countryCode); exists {
		fmt.Println(exists)
		return nil
	}
	insertQuery := `
		INSERT INTO users (country_code, phone, created_at)
		VALUES ($1, $2, NOW())
	`
	_, err := q.DB.Exec(insertQuery, countryCode, phone)
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		return err
	}
	return nil
}
func (q *QueryMeth) GetDataPhone(phone string) (*util.AllDetails, error) {
	user := util.AllDetails{}
	query := `
		SELECT phone, country_code, fullName, profilePhoto
		FROM users
		WHERE phone = $1 
	`
	err := q.DB.QueryRow(query, phone).Scan(&user.Phone, &user.CountryCode, &user.FullName, &user.ProfilePhoto)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("GetUserByPhoneAndCountryCode error: %v", err)
	}
	fmt.Println("phone-user", user.ProfilePhoto)

	return &user, nil
}
func (q *QueryMeth) UpdateUserProfile(phone, profilePhoto, fullName string) error {
	query := `
		UPDATE users
		SET 
			fullName = $1,
			profilePhoto = $2,
			updated_at = NOW()
		WHERE phone = $3
	`
	_, err := q.DB.Exec(query, fullName, profilePhoto, phone)
	if err != nil {
		return fmt.Errorf("UpdateUserProfile error: %v", err)
	}

	return nil
}

func (q *QueryMeth) GetAllUsers() ([]util.AllDetails, error) {
	query := `
        SELECT phone, country_code, fullName, profilePhoto
        FROM users
    `

	rows, err := q.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers query error: %v", err)
	}
	defer rows.Close()

	var users []util.AllDetails

	for rows.Next() {
		var user util.AllDetails
		if err := rows.Scan(&user.Phone, &user.CountryCode, &user.FullName, &user.ProfilePhoto); err != nil {
			return nil, fmt.Errorf("GetAllUsers scan error: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers rows error: %v", err)
	}

	return users, nil
}
