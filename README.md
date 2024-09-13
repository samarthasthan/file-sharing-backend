# File Sharing & Management System Architecture

```
Name - Samarth Asthan
Reg.No - 21BRS1248
Mail - samarthasthan27@gmail.com
```

## What I Planned - Let's See the Final One

This document outlines the architecture for a File Sharing and Management System that I plan to implement. The system will support multiple users, file uploads, metadata storage, and notifications, while focusing on scalability and performance optimization. The architecture makes use of Kafka, gRPC, and several other technologies to ensure decoupled communication, efficient storage, and fast file search capabilities.

## Architecture Overview

![App Screenshot](./assets/trademarkia_architecture.png)

The system is designed to handle user interactions for uploading, managing, and sharing files. It will use an API Gateway for routing requests and Kafka for event-driven communication. Inter-service communication will be optimized using gRPC. The architecture will ensure file storage, metadata indexing, caching, and notifications, with monitoring and logging integrated for system reliability.

## Technologies Planned

- **Go** for backend service development.
- **gRPC** for fast, low-latency inter-service communication.
- **Kafka** for event-driven architecture and decoupled communication.
- **PostgreSQL** for metadata storage.
- **Elasticsearch** for file metadata search and indexing.
- **Redis** for caching.
- **Docker** for containerization of services.
- **S3** for file storage (can be replaced by local storage).
- **Grafana, Prometheus, Loki, Zipkin** for monitoring, logging, and tracing.

## Core Components

### 1. **Users**

- Users interact with the system through the API Gateway. They can upload, manage, and share files.

### 2. **API Gateway**

- Acts as the main interface between the users and the backend services.
- Routes user requests to appropriate services.

### 3. **Kafka (Message Broker)**

- Used for asynchronous communication between different components.
- Handles events like file upload, processing, and notifications.
- Facilitates decoupled communication, ensuring system scalability.

### 4. **Storage Service**

- Files are stored in an S3 bucket (or local storage if S3 is not available).
- Metadata of the files is stored in PostgreSQL for persistence.

### 5. **Metadata and Search**

- File metadata is stored in PostgreSQL for quick retrieval.
- Elasticsearch is used for indexing and searching files based on their metadata.
- This enables fast search functionality within the system.

### 6. **Cache**

- Redis is used as a caching layer to store frequently accessed file metadata.
- Improves the system's performance by reducing the load on the database.

### 7. **File Processing**

- Once files are uploaded, they are processed (e.g., compressed, resized, etc.).
- The processed files are then stored in the S3 bucket.

### 8. **Deletion Service**

- Handles requests for file deletion.
- Ensures proper cleanup of files from both the storage (S3) and the metadata store (PostgreSQL).

### 9. **Notification Service**

- Sends notifications to users based on events like file uploads or deletions.
- Integrates with Kafka to listen for file-related events and triggers notifications accordingly.

### 10. **Monitoring, Logging & Tracing**

- **Grafana, Prometheus**: Used for system monitoring, performance tracking, and alerting.
- **Loki**: Provides centralized logging to track events and errors across services.
- **Zipkin**: Ensures distributed tracing, allowing insights into request flows and latency within the system.

## System Flow

1. Users send file upload requests via the **API Gateway**.
2. The request is routed to the **File Process** service using **gRPC**, which uploads the file to **S3** and stores the metadata in **PostgreSQL**.
3. The metadata is indexed in **Elasticsearch** for easy searching.
4. **Kafka** manages the communication between services, allowing events like **File Upload** and **File Deletion** to be handled asynchronously.
5. **Redis** caches frequently requested metadata to improve performance.
6. The **Notification Service** listens to Kafka and sends out notifications when necessary (e.g., file upload complete).
7. All events, requests, and performance metrics are tracked using **Grafana**, **Prometheus**, **Loki**, and **Zipkin**.
