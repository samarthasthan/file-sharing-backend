CREATE TABLE Users (
    UserID CHAR(36) PRIMARY KEY,                -- Unique ID for each user, typically a UUID
    FirstName VARCHAR(255) NOT NULL,            -- User's first name
    LastName VARCHAR(255) NOT NULL,             -- User's last name
    Email VARCHAR(255) NOT NULL UNIQUE,         -- Unique email address for each user
    IsVerified BOOLEAN DEFAULT FALSE,           -- Flag to indicate if the user's email is verified
    Password VARCHAR(255) NOT NULL,             -- User's hashed password
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the user was created
    UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the user was last updated
    DeletedAt TIMESTAMP NULL                   -- Timestamp for when the user was deleted (soft delete)
    
);


CREATE TABLE Files (
    FileID CHAR(36) PRIMARY KEY,                 -- Unique ID for each file
    UserID CHAR(36) NOT NULL,                   -- Foreign key referencing the user
    FileName VARCHAR(255) NOT NULL,             -- Original file name
    FileSize BIGINT NOT NULL,                   -- Size of the file in bytes
    FileType VARCHAR(50) NOT NULL,              -- MIME type of the file      
    StorageLocation TEXT NOT NULL,             -- Storage location (S3 bucket or local path)
    UploadDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the file was uploaded
    IsProcessed BOOLEAN DEFAULT FALSE,          -- Flag indicating if the file has been processed
    ExpiresAt TIMESTAMP NULL,                   -- Timestamp when the file will expire
    UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_user
        FOREIGN KEY (UserID)
        REFERENCES Users(UserID)
        ON DELETE CASCADE                    -- Automatically delete files if the user is deleted
);
