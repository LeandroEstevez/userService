package util

import (
	"net/http"
	"time"
	db "userMicroService/db/sqlc"

	"github.com/gin-gonic/gin"
)

const (
	YYYYMMDD = "2006-01-02"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ReturnStatusCode(response *http.Response) int {
	switch response.StatusCode {
	case 500:
		return http.StatusInternalServerError
	case 400:
		return http.StatusBadRequest
	case 404:
		return http.StatusNotFound
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func CreateRandomUser() db.User {
	return db.User{
		Username:       RandomString(6),
		HashedPassword: RandomString(6),
		FullName:       RandomFullName(),
		Email:          RandomEmail(),
		TotalExpenses:  RandomMoney(),
	}
}

func CreateRandomEntry(user db.User) db.Entry {
	date, _ := time.Parse(YYYYMMDD, "2022-12-11")

	return db.Entry{
		ID:      95,
		Owner:   user.Username,
		Name:    RandomString(6),
		DueDate: date,
		Amount:  5,
	}
}
