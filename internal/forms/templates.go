package forms

// SiteTemplate represents a predefined site configuration template
type SiteTemplate struct {
	Name        string
	Description string
	UseCase     string
	Config      SiteConfig
}

// GetSiteTemplates returns all available site templates
func GetSiteTemplates() []SiteTemplate {
	return []SiteTemplate{
		{
			Name:        "Static Website",
			Description: "HTML/CSS/JS static site with optional SSL",
			UseCase:     "Portfolio, landing page, documentation site",
			Config: SiteConfig{
				ServerName:        "example.com",
				Port:              "443",
				RootPath:          "/var/www/html",
				IndexFiles:        "index.html index.htm",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/example.com.key",
				ForceHTTPS:        true,
				IsProxy:           false,
				EnableGzip:        true,
				ClientMaxBodySize: "10M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig:      "",
			},
		},
		{
			Name:        "Single Page Application (SPA)",
			Description: "React, Vue, Angular with client-side routing",
			UseCase:     "Modern JavaScript frameworks",
			Config: SiteConfig{
				ServerName:        "app.example.com",
				Port:              "443",
				RootPath:          "/var/www/spa",
				IndexFiles:        "index.html",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/app.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/app.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           false,
				EnableGzip:        true,
				ClientMaxBodySize: "5M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig:      "# SPA routing support\ntry_files $uri $uri/ /index.html;",
			},
		},
		{
			Name:        "Node.js Application",
			Description: "Express, Nest.js, or any Node.js backend",
			UseCase:     "API server, Node.js web application",
			Config: SiteConfig{
				ServerName:        "api.example.com",
				Port:              "443",
				RootPath:          "",
				IndexFiles:        "",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/api.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/api.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           true,
				ProxyPass:         "http://localhost:3000",
				ProxyHeaders:      true,
				EnableGzip:        true,
				ClientMaxBodySize: "50M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig:      "",
			},
		},
		{
			Name:        "WordPress Site",
			Description: "WordPress blog or website with PHP-FPM",
			UseCase:     "WordPress CMS",
			Config: SiteConfig{
				ServerName:        "blog.example.com",
				Port:              "443",
				RootPath:          "/var/www/wordpress",
				IndexFiles:        "index.php index.html",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/blog.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/blog.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           false,
				EnableGzip:        true,
				ClientMaxBodySize: "100M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         true,
				PHPSocket:         "/var/run/php/php8.2-fpm.sock",
				CustomConfig: `# WordPress permalinks
location / {
    try_files $uri $uri/ /index.php?$args;
}

# Deny access to sensitive files
location ~ /\.ht {
    deny all;
}`,
			},
		},
		{
			Name:        "Laravel Application",
			Description: "Laravel PHP framework",
			UseCase:     "Laravel web application",
			Config: SiteConfig{
				ServerName:        "laravel.example.com",
				Port:              "443",
				RootPath:          "/var/www/laravel/public",
				IndexFiles:        "index.php",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/laravel.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/laravel.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           false,
				EnableGzip:        true,
				ClientMaxBodySize: "50M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         true,
				PHPSocket:         "/var/run/php/php8.2-fpm.sock",
				CustomConfig: `# Laravel routing
location / {
    try_files $uri $uri/ /index.php?$query_string;
}

# Security headers
add_header X-Frame-Options "SAMEORIGIN";
add_header X-Content-Type-Options "nosniff";`,
			},
		},
		{
			Name:        "Python/Django Application",
			Description: "Django or Flask with Gunicorn",
			UseCase:     "Python web application",
			Config: SiteConfig{
				ServerName:        "django.example.com",
				Port:              "443",
				RootPath:          "",
				IndexFiles:        "",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/django.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/django.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           true,
				ProxyPass:         "http://127.0.0.1:8000",
				ProxyHeaders:      true,
				EnableGzip:        true,
				ClientMaxBodySize: "25M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig: `# Static files
location /static/ {
    alias /var/www/django/static/;
}

# Media files
location /media/ {
    alias /var/www/django/media/;
}`,
			},
		},
		{
			Name:        "Docker Container Proxy",
			Description: "Proxy to a Docker container",
			UseCase:     "Containerized application",
			Config: SiteConfig{
				ServerName:        "container.example.com",
				Port:              "80",
				RootPath:          "",
				IndexFiles:        "",
				EnableSSL:         false,
				IsProxy:           true,
				ProxyPass:         "http://172.17.0.2:8080",
				ProxyHeaders:      true,
				EnableGzip:        true,
				ClientMaxBodySize: "10M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig: `# WebSocket support
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";`,
			},
		},
		{
			Name:        "WebSocket Application",
			Description: "Real-time app with WebSocket support",
			UseCase:     "Chat, real-time notifications, live updates",
			Config: SiteConfig{
				ServerName:        "ws.example.com",
				Port:              "443",
				RootPath:          "",
				IndexFiles:        "",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/ws.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/ws.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           true,
				ProxyPass:         "http://localhost:3001",
				ProxyHeaders:      true,
				EnableGzip:        false, // Gzip not compatible with WebSocket
				ClientMaxBodySize: "10M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig: `# WebSocket configuration
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
proxy_read_timeout 86400;`,
			},
		},
		{
			Name:        "Multiple Domain Redirect",
			Description: "Redirect www to non-www or vice versa",
			UseCase:     "Domain canonicalization",
			Config: SiteConfig{
				ServerName:        "www.example.com example.com",
				Port:              "443",
				RootPath:          "/var/www/example",
				IndexFiles:        "index.html",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/example.com.key",
				ForceHTTPS:        true,
				IsProxy:           false,
				EnableGzip:        true,
				ClientMaxBodySize: "10M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig: `# Redirect www to non-www
if ($host = www.example.com) {
    return 301 https://example.com$request_uri;
}`,
			},
		},
		{
			Name:        "API Gateway",
			Description: "Route multiple microservices through one domain",
			UseCase:     "Microservices architecture",
			Config: SiteConfig{
				ServerName:        "api.example.com",
				Port:              "443",
				RootPath:          "",
				IndexFiles:        "",
				EnableSSL:         true,
				SSLCertPath:       "/etc/ssl/certs/api.example.com.crt",
				SSLKeyPath:        "/etc/ssl/private/api.example.com.key",
				ForceHTTPS:        true,
				IsProxy:           true,
				ProxyPass:         "http://localhost:3000",
				ProxyHeaders:      true,
				EnableGzip:        true,
				ClientMaxBodySize: "10M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				CustomConfig: `# User service
location /api/users {
    proxy_pass http://localhost:3001;
    proxy_set_header Host $host;
}

# Auth service
location /api/auth {
    proxy_pass http://localhost:3002;
    proxy_set_header Host $host;
}

# Product service
location /api/products {
    proxy_pass http://localhost:3003;
    proxy_set_header Host $host;
}`,
			},
		},
		{
			Name:        "Blank Template",
			Description: "Start from scratch with minimal configuration",
			UseCase:     "Custom configuration",
			Config: SiteConfig{
				ServerName:        "",
				Port:              "80",
				RootPath:          "/var/www/html",
				IndexFiles:        "index.html",
				EnableSSL:         false,
				IsProxy:           false,
				EnableGzip:        true,
				ClientMaxBodySize: "10M",
				AccessLog:         "/var/log/nginx/access.log",
				ErrorLog:          "/var/log/nginx/error.log",
				EnablePHP:         false,
				PHPSocket:         "/var/run/php/php-fpm.sock",
				ProxyHeaders:      true,
				CustomConfig:      "",
			},
		},
	}
}

// GetTemplateByName returns a template by its name
func GetTemplateByName(name string) *SiteTemplate {
	templates := GetSiteTemplates()
	for _, t := range templates {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// GetTemplateNames returns a list of all template names
func GetTemplateNames() []string {
	templates := GetSiteTemplates()
	names := make([]string, len(templates))
	for i, t := range templates {
		names[i] = t.Name
	}
	return names
}
