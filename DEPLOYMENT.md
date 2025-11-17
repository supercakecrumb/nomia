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

## GitHub Actions CI/CD Pipelines

The project uses GitHub Actions for building, releasing, and deploying the application.

### Workflows

1. **Build and Push Workflow** (`.github/workflows/build-and-push.yml`):
   - Triggered on push to `main` branch or when a tag is pushed (tags matching `v*`)
   - Builds Docker images for backend and frontend
   - Pushes images to GitHub Container Registry (ghcr.io) with two tags: `latest` and the git tag (e.g., `v1.0.0`)

2. **Release Workflow** (`.github/workflows/release.yml`):
   - Triggered when a tag is pushed (tags matching `v*`)
   - Generates a changelog from the commits since the last tag
   - Creates a GitHub release with the changelog as release notes

3. **Deployment Workflow** (`.github/workflows/deploy.yml`):
   - Manual trigger (workflow_dispatch) with input for the tag to deploy
   - Connects to the production server via SSH
   - Pulls the specified version of the images
   - Restarts the services using `docker-compose.prod.yml`

### Required GitHub Secrets

Configure the following secrets in your GitHub repository (Settings → Secrets and variables → Actions):

- `SSH_HOST`: The hostname or IP address of your production server
- `SSH_USER`: The SSH username for the production server
- `SSH_PRIVATE_KEY`: The private key for SSH authentication (without passphrase)

### Production Server Setup for GitHub Actions Deployment

1. **Docker and Docker Compose**: Ensure Docker and Docker Compose are installed on the server.
2. **Directory Setup**: Create a directory for the project (e.g., `/opt/affirm-name`) and place the `docker-compose.prod.yml` file there.
3. **Environment Variables**: Create a `.env` file in the project directory with the required environment variables (see `backend/.env.example` for reference). At minimum, set:
   - `DB_PASSWORD`: A strong password for the PostgreSQL database
   - `FRONTEND_URL`: The public URL of the frontend application
4. **Permissions**: Ensure the SSH user has permission to run Docker commands without sudo.

### How to Use

1. **Build and Push**: Every push to `main` will build and push the `latest` images. When you push a tag (e.g., `git tag v1.0.0 && git push origin v1.0.0`), it will build and push the images with the version tag and also create a release.
2. **Deploy**: Go to the Actions tab in GitHub, select the "Deploy to Production" workflow, click "Run workflow", enter the tag to deploy (e.g., `v1.0.0`), and run.

### Notes

- The deployment workflow uses the `docker-compose.prod.yml` file from the repository. Make sure to update this file in the repository if you make changes to the production setup.
- The deployment script will:
   - Pull the new images
   - Stop the existing containers
   - Start the new containers
   - Prune unused Docker objects to free space

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