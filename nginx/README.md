# Nginx Configuration for Nomia

This directory contains the production nginx configuration for nomia.vkiel.com.

## Installation

1. Copy the configuration to nginx sites-available:
   ```bash
   sudo cp nginx/nomia.vkiel.com.conf /etc/nginx/sites-available/nomia.vkiel.com
   ```

2. Enable the site:
   ```bash
   sudo ln -s /etc/nginx/sites-available/nomia.vkiel.com /etc/nginx/sites-enabled/
   ```

3. Test the configuration:
   ```bash
   sudo nginx -t
   ```

4. Reload nginx:
   ```bash
   sudo systemctl reload nginx
   ```

## SSL/TLS Certificates

The configuration expects Let's Encrypt certificates. Set them up with:

```bash
sudo certbot certonly --webroot -w /var/www/certbot -d nomia.vkiel.com
```

## Configuration Details

- **Backend API**: Proxied from `/api/` to `127.0.0.1:9090`
- **Frontend**: Proxied from `/` to `127.0.0.1:3000`
- **Health Check**: Available at `/health` (bypasses rate limiting)
- **Rate Limiting**: 
  - API: 10 requests/second (burst 20)
  - General: 100 requests/second (burst 50)
- **SSL/TLS**: TLSv1.2 and TLSv1.3 with modern ciphers
- **Security Headers**: HSTS, X-Frame-Options, X-Content-Type-Options, X-XSS-Protection

## Upstream Configuration

The configuration defines two upstreams:
- `nomia_backend`: Backend API at 127.0.0.1:9090
- `nomia_frontend`: Frontend SPA at 127.0.0.1:3000

Both use keepalive connections for better performance.