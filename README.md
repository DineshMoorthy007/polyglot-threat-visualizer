<div align="center">

# Polyglot Threat Visualizer
**A Multi-Stack Cyber Range for Real-Time Vulnerability and Defense Visualization**

![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Java](https://img.shields.io/badge/Java-ED8B00?style=for-the-badge&logo=openjdk&logoColor=white)
![Spring Boot](https://img.shields.io/badge/Spring_Boot-6DB33F?style=for-the-badge&logo=spring&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-005C84?style=for-the-badge&logo=mysql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)

</div>

---

## Executive Summary

The Polyglot Threat Visualizer is a containerized monorepo designed to demonstrate critical database vulnerabilities and their respective mitigations in real-time. By implementing vulnerable endpoints alongside secure, industry-standard alternatives, this project provides a tangible environment for security education, testing, and architectural review.

---

## System Architecture

The architecture consists of a React frontend interfacing with two distinct backend microservices (Java and Go), both of which interact with a shared relational database. 

```mermaid
graph TD
    Client[React Dashboard] 

    subgraph Service Layer
        JavaAPI[Java Spring Boot Service<br>Port: 8080]
        GoAPI[Go Gin Service<br>Port: 8081]
    end
    
    subgraph Data Layer
        DB[(Shared MySQL Database<br>Port: 3307)]
    end

    Client -->|HTTP POST/PUT| JavaAPI
    Client -->|HTTP POST/PUT| GoAPI

    JavaAPI <-->|JDBC / Hibernate| DB
    GoAPI <-->|database/sql / GORM| DB
```

---

## The Threat Landscape

The system explicitly contrasts insecure implementations with secure paradigms. The defenses are activated globally via the application's internal state mechanism.

| Technology Stack | Demonstrated Vulnerability | Applied Defense Paradigm |
| :--- | :--- | :--- |
| **Java Spring Boot** | **SQL Injection (SQLi)** (Raw String Concatenation)<br>**Wipeout** (Unauthenticated `TRUNCATE`) | **Spring Data JPA ORM** (Parameterized Queries)<br>**Role-Based Access Control** (Exception Handling) |
| **Go (Gin + GORM)** | **IDOR** (Insecure Direct Object Reference)<br>**Denial of Service** (Unthrottled Insertion) | **JWT Ownership Validation** (Record Verification)<br>**Idempotency Keys** & Strict **Rate Limiting** |

---

## Attack & Defense Workflows

The application operates in two primary states: **Vulnerable** (Shield Inactive) and **Protected** (Shield Active). The following sequence diagram illustrates the behavioral difference across the system.

```mermaid
sequenceDiagram
    participant Dashboard as React Dashboard
    participant API as Backend Service (Java/Go)
    participant DB as MySQL Database

    rect rgb(240, 240, 240)
    note right of Dashboard: Initialization
    Dashboard->>API: POST /toggle-shield (Set State)
    API-->>Dashboard: State Synchronized
    end

    rect rgb(255, 240, 240)
    note right of Dashboard: Scenario A: Shield Inactive (Vulnerable)
    Dashboard->>API: Malicious Request (SQLi / DoS / Wipeout)
    API->>DB: Execute Raw/Unsafe Query
    DB-->>API: Data Corrupted / Compromised
    API-->>Dashboard: Attack Successful (UI Glitch)
    end

    rect rgb(240, 255, 240)
    note right of Dashboard: Scenario B: Shield Active (Protected)
    Dashboard->>API: Malicious Request (SQLi / DoS / Wipeout)
    API-->>API: Validate Input / Check Rate Limits
    API-xDB: Database Execution Halted
    API-->>Dashboard: Threat Blocked (Security Exception)
    end
```

---

## Prerequisites & Installation

The entire infrastructure is automated via Docker. To run the environment, ensure you have **Docker** and **Docker Compose** installed.

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/polyglot-threat-visualizer.git

# 2. Navigate to the project directory
cd polyglot-threat-visualizer

# 3. Build and initialize the orchestration
docker compose up --build -d
```

*Note: The MySQL container exposes port `3307` to the host machine to prevent conflicts with pre-existing local database instances.*

---

## Demonstration Guide

To evaluate the system's capabilities, follow this structured walkthrough:

1. **Access the Interface:** Navigate to `http://localhost:3000` to load the Threat Visualizer dashboard.
2. **Execute Attacks (Vulnerable State):** Initiate an attack by clicking any of the **Launch Attack** actions (e.g., Wipeout, SQL Injection, Duplication). Observe the resulting data corruption and the visual feedback (shaking/glitching animations).
3. **Engage Defenses:** Press the hidden keybind `Alt + X`. This keystroke broadcasts a state-change request to both backends, instantly engaging the security protocols without requiring a page reload or visible UI toggle.
4. **Verify Mitigations (Protected State):** Repeat the previously successful attacks. The backends will intercept the payloads, safely reject the requests, and the dashboard will display a "Threat Blocked" notification, confirming system integrity.
