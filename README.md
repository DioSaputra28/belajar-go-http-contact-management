# Contact Management API

REST API untuk manajemen kontak dengan fitur user authentication, contact management, dan address management.

## Tech Stack

- **Go** - Programming language
- **MySQL** - Database
- **Docker** - Containerization
- **Swagger** - API Documentation

## Prerequisites

Sebelum memulai, pastikan kamu sudah install:

- [Docker](https://docs.docker.com/get-docker/) (version 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.0+)

## Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd <project-directory>
```

### 2. Setup Environment Variables

Copy file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Edit file `.env` sesuai kebutuhan:

```env
# Application
APP_PORT=8080
HOST_PORT=8080

# Database Configuration
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_secure_password
DB_NAME=contact_management

# MySQL Host Port
MYSQL_HOST_PORT=3306
```

**Penting:** Ganti `DB_PASSWORD` dengan password yang aman!

### 3. Build dan Run dengan Docker Compose

```bash
docker-compose up --build
```

Atau run di background:

```bash
docker-compose up -d --build
```

### 4. Cek Status Container

```bash
docker-compose ps
```

Output yang diharapkan:

```
NAME                        STATUS              PORTS
contact-management-app      Up                  0.0.0.0:8080->8080/tcp
contact-management-db       Up (healthy)        0.0.0.0:3306->3306/tcp
```

### 5. Akses Aplikasi

- **API Base URL:** `http://localhost:8080`
- **Swagger Documentation:** `http://localhost:8080/swagger/index.html`

## Docker Commands

### Start Services

```bash
docker-compose up
```

### Start Services (Background)

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### Stop Services dan Hapus Volumes

```bash
docker-compose down -v
```

### Rebuild Services

```bash
docker-compose up --build
```

### Lihat Logs

```bash
# Semua services
docker-compose logs

# Specific service
docker-compose logs app
docker-compose logs mysql

# Follow logs (real-time)
docker-compose logs -f app
```

### Restart Services

```bash
docker-compose restart
```

### Akses Container Shell

```bash
# App container
docker-compose exec app sh

# MySQL container
docker-compose exec mysql bash
```

## Database Setup

### Akses MySQL dari Host

Jika kamu ingin akses MySQL dari host machine (menggunakan MySQL client atau GUI tools):

```bash
mysql -h 127.0.0.1 -P 3306 -u root -p
```

Masukkan password sesuai `DB_PASSWORD` di file `.env`.

### Import Database Schema

Jika kamu punya file SQL schema:

```bash
docker-compose exec -T mysql mysql -u root -p${DB_PASSWORD} ${DB_NAME} < schema.sql
```

## API Endpoints

### User Management

- `POST /user` - Register user baru
- `POST /login` - Login user
- `GET /user` - Get current user (requires auth)
- `GET /user/:id` - Get user by ID (requires auth)
- `PUT /user/:id` - Update user (requires auth)

### Contact Management

- `POST /contact` - Create contact (requires auth)
- `GET /contact` - Get all contacts (requires auth)
- `GET /contact/:id` - Get contact by ID (requires auth)
- `PUT /contact/:id` - Update contact (requires auth)
- `DELETE /contact/:id` - Delete contact (requires auth)

### Address Management

- `POST /address/` - Create address (requires auth)
- `GET /address/:contactId` - Get addresses by contact (requires auth)
- `GET /address/:contactId/:addressId` - Get specific address (requires auth)
- `PUT /address/:contactId/:addressId` - Update address (requires auth)
- `DELETE /address/:contactId/:addressId` - Delete address (requires auth)

## Development

### Run Tanpa Docker (Local Development)

1. Install dependencies:

```bash
go mod download
```

2. Setup `.env` file dengan `DB_HOST=localhost`

3. Pastikan MySQL sudah running di local

4. Run aplikasi:

```bash
go run .
```

### Update Dependencies

```bash
go get -u
go mod tidy
```

## Troubleshooting

### Port Sudah Digunakan

Jika port 8080 atau 3306 sudah digunakan, ubah di file `.env`:

```env
HOST_PORT=3000        # Ganti ke port lain
MYSQL_HOST_PORT=3307  # Ganti ke port lain
```

### Container Tidak Bisa Connect ke MySQL

Tunggu beberapa detik sampai MySQL container healthy. Cek dengan:

```bash
docker-compose logs mysql
```

### Reset Database

Hapus volume dan restart:

```bash
docker-compose down -v
docker-compose up -d
```

### Permission Denied Error

Pastikan Docker daemon running dan user kamu punya permission:

```bash
sudo usermod -aG docker $USER
```

Logout dan login kembali.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Port aplikasi di dalam container | `8080` |
| `HOST_PORT` | Port yang di-expose ke host | `8080` |
| `DB_HOST` | MySQL host (gunakan `mysql` untuk Docker) | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL username | `root` |
| `DB_PASSWORD` | MySQL password | - |
| `DB_NAME` | Database name | `contact_management` |
| `MYSQL_HOST_PORT` | MySQL port di host | `3306` |

## Project Structure

```
.
├── Dockerfile              # Docker image configuration
├── docker-compose.yaml     # Docker Compose configuration
├── .env                    # Environment variables (jangan commit!)
├── .env.example           # Environment variables template
├── main.go                # Application entry point
├── koneksi.go             # Database connection
├── user.go                # User handlers
├── contact.go             # Contact handlers
├── address.go             # Address handlers
├── middleware.go          # Authentication middleware
├── docs/                  # Swagger documentation
└── README.md              # This file
```

## License

[Your License Here]

## Contact

[Your Contact Information]
