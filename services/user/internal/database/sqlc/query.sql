-- name: RegisterUser :exec
INSERT INTO Users (UserID, FirstName, LastName, Email, Password)
VALUES ($1, $2, $3, $4, $5);

-- name: GetPasswordByEmail :one
SELECT Password FROM Users
WHERE Email = $1;