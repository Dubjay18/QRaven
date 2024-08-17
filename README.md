# QRaven API

This is a sample QRaven API built with Go.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version X.X.X)
- PostgreSQL
- Redis

### Installing

A step by step series of examples that tell you how to get a development environment running.

1. Clone the repository
```bash
git clone https://github.com/Dubjay18/qraven.git
```

2. Navigate to the project directory
```bash
cd qraven
```
3. Install dependencies
```bash
go mod download
```
4. Create a `.env` file in the root directory and add the following environment variables
```bash
PORT=8080
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_HOST=your_db_host
DB_PORT=your_db_port
REDIS_HOST=your_redis_host
REDIS_PORT=your_redis_port
```

5. Run the application
```bash
go run main.go
```

## Built with
- [Go](https://golang.org/) - The programming language used
- [Gin](https://) - The web framework used
- [Gorm](https://) - The ORM used
- [PostgreSQL](https://) - The database used 
- [Redis](https://) - The caching server used
  