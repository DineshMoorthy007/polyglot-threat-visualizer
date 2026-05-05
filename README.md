# Polyglot Threat Visualizer

Welcome to the **Polyglot Threat Visualizer**! This project is a monorepo consisting of multiple backend microservices designed to demonstrate common web vulnerabilities and how to mitigate them using a dynamic "Shield" mechanism.

## Project Structure

The repository currently contains the following services:

### 1. Java Spring Boot Backend (`backend-java`)
A Spring Boot 3 application connected to a MySQL database using Spring Data JPA and raw JdbcTemplate.
*   **Location:** `/backend-java`
*   **Technologies:** Java, Spring Boot, JPA, Hibernate, MySQL, Lombok.
*   **Threats Demonstrated:**
    *   **SQL Injection (SQLi):** `POST /api/insert` executes raw SQL strings without parameterization when the shield is off.
    *   **Dangerous Operations:** `POST /api/purge` runs a `TRUNCATE TABLE` on user data when the shield is off.
*   **Protection:** When the Shield is toggled ON, the service blocks dangerous operations (returning 403 Forbidden) and routes inserts through secure Spring Data JPA repositories.

### 2. Go Gin Backend (`backend-go`)
A lightweight Go backend utilizing the Gin framework and GORM.
*   **Location:** `/backend-go`
*   **Technologies:** Go 1.21, Gin, GORM, MySQL.
*   **Threats Demonstrated:**
    *   **Denial of Service (DoS) / Duplication:** `POST /api/go/data` accepts unthrottled rapid insertions leading to duplicate data and potential DoS when the shield is off.
    *   **Insecure Direct Object Reference (IDOR):** `PUT /api/go/data/:id` allows any user to overwrite arbitrary records simply by modifying the ID parameter when the shield is off.
*   **Protection:** When the Shield is toggled ON, the service enforces an `Idempotency-Key` check, sets a rate limit (5 req/sec), and strictly verifies ownership via an `Authorization` header before updating records via GORM.

## The Shield Mechanism

Each service contains a `/toggle` endpoint that activates or deactivates the "Shield". 

*   **Shield OFF (Vulnerable State):** The APIs simulate poorly written code, executing raw SQL, skipping authorization checks, and lacking rate limits.
*   **Shield ON (Secure State):** The APIs enforce industry best practices (parameterized queries, token ownership validation, rate limiting, and idempotency).

## Getting Started

Both services expect a MySQL instance. The connection can be configured via environment variables (e.g., in a local `.env` file).

*For detailed setup instructions, navigate into the respective backend directories.*
