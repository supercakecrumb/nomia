# Nomia Production Deployment Guide

This guide provides step-by-step instructions for deploying the Nomia application to a production server at **nomia.vkiel.com**.

## Table of Contents

1. [Server Prerequisites](#server-prerequisites)
2. [Initial Server Setup](#initial-server-setup)
3. [Domain Configuration](#domain-configuration)
4. [Docker Installation](#docker-installation)
5. [Production Configuration](#production-configuration)
6. [Nginx Reverse Proxy Setup](#nginx-reverse-proxy-setup)
7. [SSL Certificate Setup](#ssl-certificate-setup)
8. [Deploy Application](#deploy-application)
9. [Database Setup](#database-setup)
10. [Data Import](#data-import)
11. [Post-Deployment Verification](#post-deployment-verification)
12. [Monitoring & Maintenance](#monitoring--maintenance)
13. [Backup Strategy](#backup-strategy)
14. [Update Procedures](#update-procedures)
15. [Troubleshooting](#troubleshooting)

---

## Server Prerequisites

### Recommended System Requirements

- **OS**: Ubuntu 22.04 LTS or Debian 12 (recommended)
- **CPU**: 2+ cores
- **RAM**: 4GB minimum (8GB recommended)
- **Disk**: 50GB minimum (SSD recommended)
- **Network**: Public IP address with open ports 80, 443

### Required Software

- Docker Engine 24.0+
- Docker Compose v2.20+
- Nginx 1.18+
- Certbot (Let's Encrypt client)
- Git

---

## Initial Server Setup

### 1. Update System Packages

```bash
# Update package list
sudo apt update

# Upgrade installed packages
sudo apt upgrade -y

# Install essential tools
sudo apt install -y curl wget git vim ufw
```

### 2. Configure Firewall

```bash
# Enable firewall
sudo ufw enable

# Allow SSH (important - do this first!)
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check firewall status
sudo ufw status
```

### 3. Create Deployment User (Optional but Recommended)

```bash
# Create nomia user
sudo adduser nomia

# Add to docker group (after Docker is installed)
sudo usermod -aG docker nomia

# Grant sudo privileges if needed
sudo usermod -aG sudo nomia

# Switch to nomia user
su - nomia
```

---

## Domain Configuration

### DNS Setup

Configure your domain registrar to point to your server:

**A Record:**
```
Host: nomia
Target: <YOUR_SERVER_IP>
TTL: 3600
```

**Verification:**
```bash
# Check DNS propagation
dig nomia.vkiel.com +short

# Or use nslookup
nslookup nomia.vkiel.com
```

Wait for DNS propagation (can take 5-60 minutes).

---

## Docker Installation

### Install Docker Engine

```bash
# Remove old versions
sudo apt remove docker docker-engine docker.io containerd runc

# Install prerequisites
sudo apt update
sudo apt install -y ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Set up repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Verify installation
docker --version
docker compose version

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Test Docker
sudo docker run hello-world
```

### Configure Docker for Non-Root User

```bash
# Add current user to docker group
sudo usermod -aG docker $USER

# Log out and back in for group changes to take effect
# Or run: newgrp docker

# Verify
docker ps
```

---

## Production Configuration

### 1. Create Deployment Directory

```bash
# Create app directory
sudo mkdir -p /opt/nomia
sudo chown $USER:$USER /opt/nomia
cd /opt/nomia

# Clone repository (if needed for configuration files)
git clone https://github.com/YOUR_USERNAME/nomia.git .

# Or just create necessary files manually
```

### 2. Create Production docker-compose.yml

Create `/opt/nomia/docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  backend:
    image: ghcr.io/supercakecrumb/nomia/backend:latest
    container_name: nomia-backend
    restart: unless-stopped
    ports:
      - "127.0.0.1:9090:8080"  # Map external 9090 to internal 8080
    environment:
      - FIXTURE_MODE=false
      - PORT=8080
      - LOG_LEVEL=info
      - DATABASE_URL=postgresql://postgres:${DB_PASSWORD}@postgres:5432/nomia?sslmode=disable
      - FRONTEND_URL=https://nomia.vkiel.com
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - nomia-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  frontend:
    image: ghcr.io/supercakecrumb/nomia/frontend:latest
    container_name: nomia-frontend
    restart: unless-stopped
    ports:
      - "127.0.0.1:3000:80"  # Internal only, served via Nginx
    networks:
      - nomia-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:80/"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:16-alpine
    container_name: nomia-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: nomia
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d:ro
      - ./backups:/backups
    networks:
      - nomia-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    # Note: Do NOT expose 5432 to host in production for security

volumes:
  postgres_data:
    driver: local

networks:
  nomia-network:
    driver: bridge
```

### 3. Create Environment File

Create `/opt/nomia/.env`:

```bash
# Database password (use a strong password!)
DB_PASSWORD=your_very_secure_password_here_min_20_chars

# Optional: Specify version tags instead of 'latest'
# BACKEND_TAG=v1.0.0
# FRONTEND_TAG=v1.0.0
```

**Generate a secure password:**

```bash
# Generate random password
openssl rand -base64 32

# Or use pwgen
sudo apt install pwgen
pwgen -s 32 1
```

### 4. Set Proper Permissions

```bash
# Protect environment file
chmod 600 /opt/nomia/.env

# Set ownership
sudo chown -R $USER:$USER /opt/nomia
```

---

## Nginx Reverse Proxy Setup

### 1. Install Nginx

```bash
sudo apt update
sudo apt install -y nginx

# Start and enable Nginx
sudo systemctl start nginx
sudo systemctl enable nginx

# Check status
sudo systemctl status nginx
```

### 2. Create Nginx Configuration

Create `/etc/nginx/sites-available/nomia`:

```nginx
# Upstream backend API
upstream nomia_backend {
    server 127.0.0.1:9090;
    keepalive 32;
}

# Upstream frontend
upstream nomia_frontend {
    server 127.0.0.1:3000;
    keepalive 32;
}

# Rate limiting
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=general_limit:10m rate=100r/s;

# HTTP - Redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name nomia.vkiel.com;

    # Allow Let's Encrypt validation
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    # Redirect all other traffic to HTTPS
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS - Main application
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name nomia.vkiel.com;

    # SSL certificates (will be configured by certbot)
    ssl_certificate /etc/letsencrypt/live/nomia.vkiel.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/nomia.vkiel.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_stapling on;
    ssl_stapling_verify on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logging
    access_log /var/log/nginx/nomia.access.log;
    error_log /var/log/nginx/nomia.error.log;

    # Client upload size
    client_max_body_size 100M;

    # Compression
    gzip on;
    gzip_vary on;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss;
    gzip_min_length 1000;

    # API endpoints
    location /api/ {
        limit_req zone=api_limit burst=20 nodelay;

        proxy_pass http://nomia_backend;
        proxy_http_version 1.1;
        
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support (if needed)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Buffering
        proxy_buffering off;
    }

    # Frontend - SPA routing
    location / {
        limit_req zone=general_limit burst=50 nodelay;

        proxy_pass http://nomia_frontend;
        proxy_http_version 1.1;
        
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Cache static assets
        location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
            proxy_pass http://nomia_frontend;
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Health check endpoint (bypass rate limiting)
    location /api/v1/health {
        proxy_pass http://nomia_backend;
        access_log off;
    }
}
```

### 3. Enable Site Configuration

```bash
# Create symbolic link
sudo ln -s /etc/nginx/sites-available/nomia /etc/nginx/sites-enabled/

# Remove default site
sudo rm /etc/nginx/sites-enabled/default

# Test configuration
sudo nginx -t

# Don't reload yet - SSL certificates need to be created first
```

### 4. Create Certbot Directory

```bash
sudo mkdir -p /var/www/certbot
sudo chown -R www-data:www-data /var/www/certbot
```

---

## SSL Certificate Setup

### 1. Install Certbot

```bash
# Install certbot and nginx plugin
sudo apt update
sudo apt install -y certbot python3-certbot-nginx

# Verify installation
certbot --version
```

### 2. Obtain SSL Certificate

**Important:** Before running this, ensure your domain DNS is properly configured and pointing to your server.

```bash
# Stop nginx temporarily for initial certificate
sudo systemctl stop nginx

# Obtain certificate using standalone mode
sudo certbot certonly --standalone \
    --preferred-challenges http \
    -d nomia.vkiel.com \
    --email your-email@example.com \
    --agree-tos \
    --no-eff-email

# Start nginx
sudo systemctl start nginx
```

**Alternative method (using webroot with nginx running):**

```bash
# Modify nginx config temporarily to allow certbot validation
# Then run:
sudo certbot certonly --webroot \
    -w /var/www/certbot \
    -d nomia.vkiel.com \
    --email your-email@example.com \
    --agree-tos \
    --no-eff-email
```

### 3. Configure Auto-Renewal

```bash
# Test renewal process (dry run)
sudo certbot renew --dry-run

# Certbot automatically sets up renewal via systemd timer
# Verify timer is active
sudo systemctl status certbot.timer

# Check renewal configuration
sudo cat /etc/letsencrypt/renewal/nomia.vkiel.com.conf
```

### 4. Reload Nginx with SSL

```bash
# Now that certificates exist, reload nginx
sudo nginx -t
sudo systemctl reload nginx

# Verify SSL is working
curl -I https://nomia.vkiel.com
```

---

## Deploy Application

### 1. Pull Docker Images

```bash
cd /opt/nomia

# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin

# Pull images
docker compose -f docker-compose.prod.yml pull

# Verify images
docker images | grep nomia
```

**Note:** For public images, login may not be required. If you need authentication, generate a GitHub Personal Access Token with `read:packages` permission.

### 2. Create Migrations Directory

```bash
# Copy migration files from repository
mkdir -p /opt/nomia/migrations

# If you cloned the repo:
cp -r migrations/* /opt/nomia/migrations/

# Otherwise, create migrations manually or download them
```

Create `/opt/nomia/migrations/001_initial_schema.sql` with the schema from the repository.

### 3. Start Services

```bash
cd /opt/nomia

# Start all services
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs -f

# Check individual service logs
docker compose -f docker-compose.prod.yml logs backend
docker compose -f docker-compose.prod.yml logs frontend
docker compose -f docker-compose.prod.yml logs postgres
```

### 4. Verify Services

```bash
# Check backend health
curl http://localhost:9090/api/v1/health

# Check frontend
curl http://localhost:3000

# Check via public URL
curl https://nomia.vkiel.com/api/v1/health
```

---

## Database Setup

### 1. Verify Database Initialization

The migrations are automatically applied on first startup due to the volume mount in docker-compose:

```yaml
volumes:
  - ./migrations:/docker-entrypoint-initdb.d:ro
```

### 2. Verify Schema

```bash
# Connect to database
docker compose -f docker-compose.prod.yml exec postgres psql -U postgres -d nomia

# Inside psql:
\dt              # List tables
\d countries     # Describe countries table
\d names         # Describe names table

# Verify extensions
SELECT extname, extversion FROM pg_extension WHERE extname = 'pg_trgm';

# Exit
\q
```

### 3. Manual Migration (if needed)

If migrations didn't run automatically:

```bash
# Apply migration manually
docker compose -f docker-compose.prod.yml exec -T postgres \
    psql -U postgres -d nomia < migrations/001_initial_schema.sql
```

### 4. Verify Country Data

```bash
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "SELECT * FROM countries ORDER BY code;"
```

---

## Data Import

### 1. Download US SSA Data

```bash
cd /opt/nomia

# Create data directory
mkdir -p names-data/us

# Download US data
cd names-data/us
curl -L "https://www.ssa.gov/oact/babynames/names.zip" -o names.zip

# Extract
unzip names.zip
rm names.zip

# Verify files
ls -l yob*.txt | head
```

### 2. Build Import Tool

Since we're using Docker, we'll run the import from within a backend container:

```bash
cd /opt/nomia

# Copy import command into backend container
docker compose -f docker-compose.prod.yml exec backend /bin/sh -c '
# This approach works if golang is in the image
# Otherwise, we need to pre-build the import binary
'
```

**Alternative: Build import tool locally**

If Go is installed on the server:

```bash
cd /opt/nomia

# Clone full repository if not already done
git clone https://github.com/YOUR_USERNAME/nomia.git nomia-source
cd nomia-source/backend

# Build import tool
go build -o /opt/nomia/nomia-import ./cmd/import/main.go

# Run import
cd /opt/nomia
DATABASE_URL="postgresql://postgres:$(grep DB_PASSWORD .env | cut -d= -f2)@localhost:5432/nomia?sslmode=disable" \
    ./nomia-import -country=US -dir=./names-data/us
```

**Recommended: Use Docker exec with database connection**

```bash
# Get database password
DB_PASS=$(grep DB_PASSWORD /opt/nomia/.env | cut -d= -f2)

# Run import via docker
docker run --rm \
    --network nomia-network \
    -v /opt/nomia/names-data:/data \
    -e DATABASE_URL="postgresql://postgres:${DB_PASS}@nomia-postgres:5432/nomia?sslmode=disable" \
    ghcr.io/supercakecrumb/nomia/backend:latest \
    /root/server import -country=US -dir=/data/us
```

**Note:** The import command needs to be added to the backend if not already present. You may need to build a separate import binary or add import functionality to the server binary.

### 3. Verify Data Import

```bash
# Check imported data
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "
SELECT 
    c.name as country,
    COUNT(DISTINCT n.year) as years,
    MIN(n.year) as first_year,
    MAX(n.year) as last_year,
    COUNT(DISTINCT n.name) as unique_names,
    COUNT(*) as total_records,
    SUM(n.count) as total_occurrences
FROM names n
JOIN countries c ON n.country_id = c.id
GROUP BY c.name
ORDER BY c.name;
"

# Check recent years
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "
SELECT 
    year,
    COUNT(*) as total_names,
    COUNT(DISTINCT name) as unique_names,
    SUM(count) as total_occurrences
FROM names n
JOIN countries c ON n.country_id = c.id
WHERE c.code = 'US'
GROUP BY year 
ORDER BY year DESC 
LIMIT 10;
"
```

---

## Post-Deployment Verification

### 1. Service Health Checks

```bash
# Check all containers are running
docker compose -f docker-compose.prod.yml ps

# Test backend API
curl https://nomia.vkiel.com/api/v1/health

# Test backend countries endpoint
curl https://nomia.vkiel.com/api/v1/countries

# Test frontend
curl -I https://nomia.vkiel.com
```

### 2. API Endpoint Tests

```bash
# Get countries list
curl https://nomia.vkiel.com/api/v1/countries | jq

# Get meta years
curl https://nomia.vkiel.com/api/v1/meta/years?country=US | jq

# Search names
curl "https://nomia.vkiel.com/api/v1/names?country=US&year_from=2020&year_to=2023&limit=10" | jq

# Get name detail
curl "https://nomia.vkiel.com/api/v1/names/Emma/trend?country=US&year_from=2020&year_to=2023" | jq
```

### 3. Browser Testing

Open in browser:
- https://nomia.vkiel.com
- Check all pages load correctly
- Test filtering functionality
- Verify charts render properly
- Test different name searches

### 4. SSL Certificate Verification

```bash
# Check certificate details
openssl s_client -connect nomia.vkiel.com:443 -servername nomia.vkiel.com < /dev/null

# Online tools:
# https://www.ssllabs.com/ssltest/analyze.html?d=nomia.vkiel.com
```

### 5. Performance Testing

```bash
# Simple load test
ab -n 1000 -c 10 https://nomia.vkiel.com/

# API endpoint test
ab -n 100 -c 5 "https://nomia.vkiel.com/api/v1/countries"
```

---

## Monitoring & Maintenance

### 1. Log Management

**View application logs:**

```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f postgres

# Last 100 lines
docker compose -f docker-compose.prod.yml logs --tail=100 backend

# Since specific time
docker compose -f docker-compose.prod.yml logs --since 1h backend
```

**View Nginx logs:**

```bash
# Access log
sudo tail -f /var/log/nginx/nomia.access.log

# Error log
sudo tail -f /var/log/nginx/nomia.error.log

# Search for errors
sudo grep "error" /var/log/nginx/nomia.error.log
```

**Log rotation** is handled automatically by Docker and Nginx's logrotate configuration.

### 2. Resource Monitoring

```bash
# Docker container stats
docker stats

# Disk usage
df -h
docker system df

# Database size
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "
SELECT 
    pg_size_pretty(pg_database_size('nomia')) as database_size,
    pg_size_pretty(pg_total_relation_size('names')) as names_table_size;
"
```

### 3. Health Monitoring Script

Create `/opt/nomia/health-check.sh`:

```bash
#!/bin/bash

# Health check script
echo "=== Nomia Health Check ==="
echo "Time: $(date)"
echo ""

# Check containers
echo "1. Container Status:"
docker compose -f /opt/nomia/docker-compose.prod.yml ps
echo ""

# Check backend
echo "2. Backend Health:"
curl -s https://nomia.vkiel.com/api/v1/health | jq
echo ""

# Check database
echo "3. Database:"
docker compose -f /opt/nomia/docker-compose.prod.yml exec -T postgres \
    psql -U postgres -d nomia -c "SELECT COUNT(*) as total_names FROM names;"
echo ""

# Check disk
echo "4. Disk Usage:"
df -h /
echo ""

# Check logs for errors
echo "5. Recent Errors:"
docker compose -f /opt/nomia/docker-compose.prod.yml logs --tail=50 | grep -i error | tail -10
echo ""

echo "=== Check Complete ==="
```

Make executable and run:

```bash
chmod +x /opt/nomia/health-check.sh
/opt/nomia/health-check.sh
```

### 4. Setup Cron for Regular Checks

```bash
# Edit crontab
crontab -e

# Add health check every hour
0 * * * * /opt/nomia/health-check.sh >> /var/log/nomia-health.log 2>&1

# Add disk space alert (90% threshold)
0 0 * * * df -h / | awk '$5 > 90 {print "Disk usage high: " $5}' | mail -s "Disk Alert" your-email@example.com
```

---

## Backup Strategy

### 1. Database Backup Script

Create `/opt/nomia/backup-db.sh`:

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/opt/nomia/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
DB_PASSWORD=$(grep DB_PASSWORD /opt/nomia/.env | cut -d= -f2)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
echo "Starting backup at $(date)"
docker compose -f /opt/nomia/docker-compose.prod.yml exec -T postgres \
    pg_dump -U postgres -d nomia | gzip > "$BACKUP_DIR/nomia_backup_$TIMESTAMP.sql.gz"

# Verify backup
if [ -f "$BACKUP_DIR/nomia_backup_$TIMESTAMP.sql.gz" ]; then
    SIZE=$(du -h "$BACKUP_DIR/nomia_backup_$TIMESTAMP.sql.gz" | cut -f1)
    echo "Backup completed successfully: nomia_backup_$TIMESTAMP.sql.gz ($SIZE)"
else
    echo "ERROR: Backup failed"
    exit 1
fi

# Remove backups older than 30 days
find $BACKUP_DIR -name "nomia_backup_*.sql.gz" -mtime +30 -delete

# List recent backups
echo "Recent backups:"
ls -lh $BACKUP_DIR/nomia_backup_*.sql.gz | tail -5

echo "Backup complete at $(date)"
```

Make executable:

```bash
chmod +x /opt/nomia/backup-db.sh
```

### 2. Schedule Automatic Backups

```bash
# Edit crontab
crontab -e

# Daily backup at 2 AM
0 2 * * * /opt/nomia/backup-db.sh >> /var/log/nomia-backup.log 2>&1

# Weekly backup to remote location (e.g., S3)
0 3 * * 0 /opt/nomia/backup-db.sh && aws s3 cp /opt/nomia/backups/ s3://your-bucket/nomia-backups/ --recursive
```

### 3. Restore from Backup

```bash
# Stop services
docker compose -f docker-compose.prod.yml stop backend

# Restore database
gunzip < /opt/nomia/backups/nomia_backup_TIMESTAMP.sql.gz | \
    docker compose -f docker-compose.prod.yml exec -T postgres \
    psql -U postgres -d nomia

# Restart services
docker compose -f docker-compose.prod.yml start backend

# Verify
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "SELECT COUNT(*) FROM names;"
```

### 4. Backup Docker Volumes

```bash
# Backup PostgreSQL volume
docker run --rm \
    -v nomia_postgres_data:/data \
    -v /opt/nomia/backups:/backup \
    alpine tar czf /backup/postgres_volume_$(date +%Y%m%d).tar.gz /data
```

---

## Update Procedures

### 1. Update Application (Rolling Update)

```bash
cd /opt/nomia

# Pull latest images
docker compose -f docker-compose.prod.yml pull

# Recreate containers with new images (zero-downtime)
docker compose -f docker-compose.prod.yml up -d

# Check logs
docker compose -f docker-compose.prod.yml logs -f
```

### 2. Update Specific Service

```bash
# Update backend only
docker compose -f docker-compose.prod.yml pull backend
docker compose -f docker-compose.prod.yml up -d backend

# Update frontend only
docker compose -f docker-compose.prod.yml pull frontend
docker compose -f docker-compose.prod.yml up -d frontend
```

### 3. Rollback to Previous Version

```bash
# Pull specific version
docker pull ghcr.io/supercakecrumb/nomia/backend:v1.0.0
docker pull ghcr.io/supercakecrumb/nomia/frontend:v1.0.0

# Update docker-compose to use specific tags
# Edit docker-compose.prod.yml to use version tags

# Restart with old version
docker compose -f docker-compose.prod.yml up -d
```

### 4. Database Migrations

```bash
# Backup before migration
/opt/nomia/backup-db.sh

# Apply new migration
docker compose -f docker-compose.prod.yml exec -T postgres \
    psql -U postgres -d nomia < migrations/002_new_migration.sql

# Verify
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "\dt"
```

### 5. System Updates

```bash
# Update system packages
sudo apt update
sudo apt upgrade -y

# Update Docker
sudo apt install --only-upgrade docker-ce docker-ce-cli containerd.io

# Reboot if kernel was updated
sudo reboot
```

---

## Troubleshooting

### Common Issues

#### Issue 1: Container Won't Start

```bash
# Check logs
docker compose -f docker-compose.prod.yml logs backend

# Check if port is already in use
sudo netstat -tlnp | grep 9090

# Restart specific service
docker compose -f docker-compose.prod.yml restart backend

# Force recreate
docker compose -f docker-compose.prod.yml up -d --force-recreate backend
```

#### Issue 2: Database Connection Error

```bash
# Check if PostgreSQL is running
docker compose -f docker-compose.prod.yml ps postgres

# Check database logs
docker compose -f docker-compose.prod.yml logs postgres

# Test connection manually
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "SELECT 1;"

# Verify DATABASE_URL environment variable
docker compose -f docker-compose.prod.yml exec backend env | grep DATABASE_URL

# Check if containers are on same network
docker network inspect nomia_nomia-network
```

#### Issue 3: 502 Bad Gateway from Nginx

```bash
# Check if backend is running
curl http://localhost:9090/api/v1/health

# Check Nginx error log
sudo tail -50 /var/log/nginx/nomia.error.log

# Test Nginx configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx

# Check if backend is accessible from host
telnet localhost 9090
```

#### Issue 4: SSL Certificate Issues

```bash
# Check certificate status
sudo certbot certificates

# Renew certificate manually
sudo certbot renew --force-renewal

# Check certificate files exist
ls -l /etc/letsencrypt/live/nomia.vkiel.com/

# Test SSL
curl -vI https://nomia.vkiel.com
```

#### Issue 5: Out of Disk Space

```bash
# Check disk usage
df -h

# Clean Docker system
docker system prune -a --volumes

# Remove old logs
sudo journalctl --vacuum-time=7d

# Remove old backups
find /opt/nomia/backups -mtime +30 -delete
```

#### Issue 6: High CPU/Memory Usage

```bash
# Check container stats
docker stats

# Check processes
top

# Restart resource-heavy container
docker compose -f docker-compose.prod.yml restart backend

# Check database queries
docker compose -f docker-compose.prod.yml exec postgres \
    psql -U postgres -d nomia -c "
SELECT pid, query, state 
FROM pg_stat_activity 
WHERE state != 'idle' 
ORDER BY query_start;
"
```

### Debug Commands

```bash
# Enter container shell
docker compose -f docker-compose.prod.yml exec backend sh
docker compose -f docker-compose.prod.yml exec postgres bash

# Check container environment variables
docker compose -f docker-compose.prod.yml exec backend env

# View container processes
docker compose -f docker-compose.prod.yml top

# Inspect container
docker inspect nomia-backend

# Check network connectivity between containers
docker compose -f docker-compose.prod.yml exec backend ping postgres

# View container resource usage
docker stats nomia-backend nomia-frontend nomia-postgres
```

### Log Analysis

```bash
# Search for specific error
docker compose -f docker-compose.prod.yml logs | grep "error"

# Count error occurrences
docker compose -f docker-compose.prod.yml logs | grep -c "ERROR"

# Export logs
docker compose -f docker-compose.prod.yml logs > nomia-logs-$(date +%Y%m%d).log

# Monitor logs in real-time with filtering
docker compose -f docker-compose.prod.yml logs -f | grep -E "ERROR|WARN"
```

### Recovery Procedures

**Complete system restart:**

```bash
cd /opt/nomia

# Stop all services
docker compose -f docker-compose.prod.yml down

# Start services with fresh state
docker compose -f docker-compose.prod.yml up -d

# Verify
docker compose -f docker-compose.prod.yml ps
```

**Emergency database restore:**

```bash
# Stop backend
docker compose -f docker-compose.prod.yml stop backend

# Restore latest backup
gunzip < /opt/nomia/backups/nomia_backup_LATEST.sql.gz | \
    docker compose -f docker-compose.prod.yml exec -T postgres \
    psql -U postgres -d nomia

# Restart services
docker compose -f docker-compose.prod.yml start backend
```

---

## Security Best Practices

### 1. Firewall Configuration

Ensure only necessary ports are open:

```bash
sudo ufw status
# Should show: 22 (SSH), 80 (HTTP), 443 (HTTPS)
```

### 2. Regular Updates

```bash
# Update system weekly
sudo apt update && sudo apt upgrade -y

# Update Docker images
docker compose -f docker-compose.prod.yml pull
```

### 3. Strong Passwords

- Use 20+ character random passwords
- Store in password manager
- Never commit to version control

### 4. SSL/TLS

- Always use HTTPS
- Keep certificates renewed
- Use strong SSL configuration

### 5. Database Security

- Don't expose PostgreSQL port to internet
- Use strong database password
- Regular backups
- Enable connection encryption if needed

### 6. Monitoring

- Set up alerts for downtime
- Monitor logs for suspicious activity
- Track resource usage

---

## Performance Optimization

### 1. Database Optimization

```sql
-- Analyze query performance
EXPLAIN ANALYZE SELECT ...;

-- Update statistics
ANALYZE;

-- Reindex if needed
REINDEX DATABASE nomia;

-- Check index usage
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan;
```

### 2. Docker Optimization

```yaml
# In docker-compose.prod.yml, add resource limits:
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### 3. Nginx Caching

Add to Nginx config:

```nginx
# Cache zone
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=api_cache:10m max_size=100m inactive=60m;

# In location block
proxy_cache api_cache;
proxy_cache_valid 200 5m;
```

---

## Support & Contacts

For issues or questions:

1. Check this documentation
2. Review logs: `/var/log/nginx/` and `docker compose logs`
3. Verify configuration files
4. Check GitHub repository issues
5. Contact system administrator

---

## Appendix

### A. Complete File Structure

```
/opt/nomia/
├── docker-compose.prod.yml
├── .env
├── migrations/
│   └── 001_initial_schema.sql
├── backups/
│   └── *.sql.gz
├── names-data/
│   └── us/
│       └── yob*.txt
├── backup-db.sh
└── health-check.sh
```

### B. Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| DB_PASSWORD | PostgreSQL password | `random_secure_password_123` |
| FIXTURE_MODE | Use fixture data (false in prod) | `false` |
| PORT | Backend internal port | `8080` |
| LOG_LEVEL | Logging level | `info` |
| FRONTEND_URL | Frontend URL for CORS | `https://nomia.vkiel.com` |

### C. Port Reference

| Service | Internal Port | External/Nginx |
|---------|--------------|----------------|
| Backend | 8080 | 9090 (localhost only) |
| Frontend | 80 | 3000 (localhost only) |
| PostgreSQL | 5432 | Not exposed |
| Nginx | 80, 443 | Public |

### D. Useful Links

- **Nomia Repository**: https://github.com/YOUR_USERNAME/nomia
- **Docker Documentation**: https://docs.docker.com/
- **Let's Encrypt**: https://letsencrypt.org/
- **Nginx Documentation**: https://nginx.org/en/docs/
- **PostgreSQL Documentation**: https://www.postgresql.org/docs/

---

**Document Version**: 1.0  
**Last Updated**: 2025-11-17  
**Maintained By**: DevOps Team
