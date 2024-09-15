-- name: UploadFileByEmail :exec
INSERT INTO Files (FileID,UserID, FileName, FileSize, FileType, StorageLocation, UploadDate, ExpiresAt)
VALUES (
    $1,
    (SELECT UserID FROM Users WHERE Email = $2),
    $3, $4, $5, $6, $7, $8
);


-- name: GetFilesByUser :many
SELECT FileID, FileName, FileSize, FileType, StorageLocation, UploadDate, IsProcessed, ExpiresAt
FROM Files
WHERE UserID = (SELECT UserID FROM Users WHERE Email = $1)  ORDER BY UploadDate DESC;

-- name: GetFileByID :one
SELECT FileID, UserID, FileName, FileSize, FileType, StorageLocation, UploadDate, IsProcessed, ExpiresAt, UpdatedAt
FROM Files
WHERE FileID = $1;

-- name: MarkFileAsProcessed :exec
UPDATE Files
SET IsProcessed = TRUE, UpdatedAt = CURRENT_TIMESTAMP
WHERE FileID = $1;

-- name: DeleteFile :exec
DELETE FROM Files
WHERE FileID = $1;

-- name: GetExpiredFiles :many
SELECT FileID, UserID, FileName, FileSize, FileType, StorageLocation, UploadDate, IsProcessed, ExpiresAt, UpdatedAt
FROM Files
WHERE ExpiresAt < CURRENT_TIMESTAMP AND IsProcessed = FALSE;

-- name: DeleteExpiredFiles :exec
DELETE FROM Files
WHERE ExpiresAt < CURRENT_TIMESTAMP;

-- name: SearchFiles :many
SELECT FileID, UserID, FileName, FileSize, FileType, StorageLocation, UploadDate, IsProcessed, ExpiresAt, UpdatedAt
FROM Files
WHERE FileName ILIKE '%' || $1 || '%' OR 
      UploadDate::date = $2 OR
      FileType ILIKE '%' || $3 || '%';
