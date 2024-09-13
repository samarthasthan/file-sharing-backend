-- name: CreateAccount :exec
INSERT INTO Users (UserID, FirstName, LastName, Email, Password)
VALUES ($1, $2, $3, $4, $5);