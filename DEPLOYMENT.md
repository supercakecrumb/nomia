# Deployment Guide

This guide covers deploying the Affirm Name backend to production.

## Prerequisites

- Docker and Docker Compose installed
- GitHub account with access to GitHub Container Registry
- Production server with Docker installed

## Local Production Testing

### 1. Build and Run with Production Docker Compose

```bash
# Set environment variables
export DB_PASSWORD=secure_password
export FRONTEND_URL=https://your-frontend-url.com

# Build and start
docker-compose -f docker-compose.prod.yml up --build
```

### 2. Test Health Endpoint

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "database": "connected"
}
```

### 3. Test API Endpoints

```bash
# Test meta endpoints
curl http://localhost:8080/api/meta/years
curl http://localhost:8080/api/meta/countries

# Test names endpoints
curl "http://localhost:8080/api/names?page=1&page_size=10"
curl "http://localhost:8080/api/names/trend?name=Oliver"
```

## GitHub Actions Deployment

### 1. Configure GitHub Secrets

Go to your repository → Settings → Secrets and add:

- `DB_PASSWORD`: Production database password
- `FRONTEND_URL`: Production frontend URL
- Additional secrets for your deployment method (SSH keys, API tokens, etc.)

### 2. Deployment Workflow

The deployment workflow (`.github/workflows/deploy.yml`) automatically:

**On push to main:**
- Builds Docker image
- Pushes to GitHub Container Registry
- Deploys to staging environment

**On release:**
- Builds Docker image
- Pushes to GitHub Container Registry with version tag
- Deploys to production environment

### 3. Manual Deployment

Trigger manual deployment:
```bash
gh workflow run deploy.yml
```

## Production Server Setup

### Option 1: Docker Compose on VPS

1. **Install Docker on server:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

2. **Clone repository:**
```bash
git clone https://github.com/your-org/affirm-name-backend.git
cd affirm-name-backend
```

3. **Configure environment:**
```bash
cp backend/.env.example backend/.env
# Edit backend/.env with production values
```

4. **Start services:**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

5. **Verify deployment:**
```bash
docker-compose -f docker-compose.prod.yml ps
curl http://localhost:8080/health
```

### Option 2: Kubernetes

See [kubernetes/README.md](kubernetes/README.md) for Kubernetes deployment guide.

### Option 3: Cloud Platforms

#### Railway
```bash
railway login
railway link
railway up
```

#### Heroku
```bash
heroku create affirm-name-backend
heroku addons:create heroku-postgresql:hobby-dev
git push heroku main
```

#### Fly.io
```bash
fly launch
fly deploy
```

## Database Migrations

### Run Migrations on Production

```bash
# Using migrate tool
migrate -path migrations -database "postgresql://user:pass@host:5432/affirm_name?sslmode=disable" up

# Or connect to database container
docker-compose -f docker-compose.prod.yml exec postgres psql -U postgres -d affirm_name
```

### Import Initial Data

```bash
# Import US name data
bash backend/scripts/import-us-data.sh
```

## Monitoring

### Health Check
```bash
curl https://your-domain.com/health
```

### Database Status
```bash
docker-compose -f docker-compose.prod.yml exec postgres psql -U postgres -d affirm_name -c "SELECT COUNT(*) FROM names;"
```

### View Logs
```bash
# Backend logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Database logs
docker-compose -f docker-compose.prod.yml logs -f postgres
```

## Backup and Recovery

### Backup Database
```bash
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres affirm_name > backup.sql
```

### Restore Database
```bash
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U postgres affirm_name < backup.sql
```

## Scaling

### Horizontal Scaling
Add multiple backend instances behind a load balancer:

```yaml
services:
  backend:
    deploy:
      replicas: 3
```

### Database Connection Pooling
Already configured in the application using pgxpool.

## Troubleshooting

### Check Container Status
```bash
docker-compose -f docker-compose.prod.yml ps
```

### View Container Logs
```bash
docker-compose -f docker-compose.prod.yml logs backend
```

### Restart Services
```bash
docker-compose -f docker-compose.prod.yml restart
```

### Reset Database
```bash
docker-compose -f docker-compose.prod.yml down -v
docker-compose -f docker-compose.prod.yml up -d
```

## Security Checklist

- [ ] Use strong database passwords
- [ ] Enable SSL/TLS for database connections
- [ ] Set up firewall rules
- [ ] Use environment variables for secrets
- [ ] Enable CORS only for trusted origins
- [ ] Keep Docker images updated
- [ ] Regular security audits
- [ ] Monitor logs for suspicious activity

## Performance Optimization

- Use Redis for caching (future enhancement)
- Enable gzip compression
- Set up CDN for static content
- Database query optimization
- Connection pooling (already configured)

## Support

For deployment issues, check:
- GitHub Issues: https://github.com/your-org/affirm-name-backend/issues
- Documentation: https://github.com/your-org/affirm-name-backend/wiki