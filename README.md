# Morphlink: A Smart URL Shortener

Morphlink is a full-stack, modular URL shortening service with a complete local development and observability environment. It includes real-time analytics, user authentication, and a developer dashboard for monitoring application health.

This project is built as a **modular monolith**, with clean separation between its core components:
*   **Auth Service:** Handles user registration, login, and JWT-based authentication.
*   **Link Manager:** Provides APIs for users to create and manage their links.
*   **High-Speed Redirector:** Efficiently redirects short URLs to their original destination.
*   **Analytics Engine:** Tracks link clicks in real-time and updates the user dashboard via WebSockets.

## Features

*   **User Authentication:** Secure user registration and login with JWT.
*   **URL Shortening:** Create unique, short IDs for any long URL.
*   **Real-time Analytics:** User dashboard updates instantly with click counts using WebSockets.
*   **Developer Observability:** Includes a full local monitoring stack with Prometheus and Grafana to visualize application traffic, CPU, and memory usage.
*   **Monorepo Structure:** Backend (Go) and Frontend (HTML/JS/CSS) code are managed in a single repository.
*   **Containerized Services:** Uses `docker-compose` to manage local PostgreSQL, Prometheus, and Grafana instances.

## Tech Stack

*   **Backend:** Go (Golang)
    *   **Router:** `chi`
    *   **Database:** PostgreSQL
    *   **Real-time:** `gorilla/websocket`
    *   **Observability:** `prometheus/client_golang`
*   **Frontend:** HTML5, CSS3, Vanilla JavaScript
*   **Database:** PostgreSQL
*   **Infrastructure:** Docker, Docker Compose
*   **Monitoring:** Prometheus, Grafana

## Prerequisites

Before you begin, ensure you have the following installed on your system:
*   **Go** (version 1.18 or higher)
*   **Docker Desktop** (must be running)

## Local Development Setup

Follow these steps to get the entire application and its monitoring dashboards running on your local machine.

### 1. Clone the Repository
Clone this monorepo to your local machine.
```bash
git clone https://github.com/AyushmanKS/appointy-task.git
```
cd appointy-task


### 2. Configure Environment Variables
#### The backend needs a configuration file to connect to the local database.
#### Navigate to the backend/ directory.
#### Create a new file named .env.
#### Copy and paste the following content into backend/.env:

### 3. Start the Infrastructure Services
### This command will start the PostgreSQL database, Prometheus, and Grafana containers in the background.
### From the root directory of the project (appointy-task/), run:

```bash
docker-compose up -d
```
