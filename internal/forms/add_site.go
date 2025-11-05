package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// SiteConfig holds all configuration for a new NGINX site
type SiteConfig struct {
	// Basic Configuration
	ServerName   string
	Port         string
	RootPath     string
	IndexFiles   string
	
	// SSL Configuration
	EnableSSL    bool
	SSLCertPath  string
	SSLKeyPath   string
	ForceHTTPS   bool
	
	// Proxy Configuration
	IsProxy      bool
	ProxyPass    string
	ProxyHeaders bool
	
	// Advanced Options
	EnableGzip   bool
	ClientMaxBodySize string
	AccessLog    string
	ErrorLog     string
	
	// PHP Configuration
	EnablePHP    bool
	PHPSocket    string
	
	// Custom Directives
	CustomConfig string
	
	// Confirmation
	Confirmed bool
}

// NewAddSiteForm creates a new form for adding an NGINX site
func NewAddSiteForm() (*huh.Form, *SiteConfig) {
	config := &SiteConfig{
		Port:              "80",
		RootPath:          "/var/www/html",
		IndexFiles:        "index.html index.htm",
		ClientMaxBodySize: "10M",
		AccessLog:         "/var/log/nginx/access.log",
		ErrorLog:          "/var/log/nginx/error.log",
		PHPSocket:         "/var/run/php/php-fpm.sock",
		ProxyHeaders:      true,
		EnableGzip:        true,
	}

	form := huh.NewForm(
		// Page 1: Basic Configuration
		huh.NewGroup(
			huh.NewInput().
				Title("Server Name").
				Description("Domain name or server name (e.g., example.com)").
				Placeholder("example.com").
				Value(&config.ServerName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("server name is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Listen Port").
				Description("Port number for the server").
				Placeholder("80").
				Value(&config.Port).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("port is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Root Path").
				Description("Document root directory").
				Placeholder("/var/www/html").
				Value(&config.RootPath),

			huh.NewInput().
				Title("Index Files").
				Description("Space-separated list of index files").
				Placeholder("index.html index.htm").
				Value(&config.IndexFiles),
		).Title("Basic Configuration").Description("Configure basic server settings"),

		// Page 2: SSL Configuration
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable SSL/TLS").
				Description("Enable HTTPS for this site").
				Value(&config.EnableSSL),

			huh.NewInput().
				Title("SSL Certificate Path").
				Description("Path to SSL certificate file").
				Placeholder("/etc/ssl/certs/cert.pem").
				Value(&config.SSLCertPath).
				Validate(func(s string) error {
					if config.EnableSSL && s == "" {
						return fmt.Errorf("SSL certificate path is required when SSL is enabled")
					}
					return nil
				}),

			huh.NewInput().
				Title("SSL Key Path").
				Description("Path to SSL private key file").
				Placeholder("/etc/ssl/private/key.pem").
				Value(&config.SSLKeyPath).
				Validate(func(s string) error {
					if config.EnableSSL && s == "" {
						return fmt.Errorf("SSL key path is required when SSL is enabled")
					}
					return nil
				}),

			huh.NewConfirm().
				Title("Force HTTPS Redirect").
				Description("Redirect HTTP to HTTPS automatically").
				Value(&config.ForceHTTPS),
		).Title("SSL/TLS Configuration").Description("Configure HTTPS settings"),

		// Page 3: Proxy Configuration
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Reverse Proxy").
				Description("Use this site as a reverse proxy").
				Value(&config.IsProxy),

			huh.NewInput().
				Title("Proxy Pass URL").
				Description("Backend server URL (e.g., http://localhost:3000)").
				Placeholder("http://localhost:3000").
				Value(&config.ProxyPass).
				Validate(func(s string) error {
					if config.IsProxy && s == "" {
						return fmt.Errorf("proxy pass URL is required when proxy is enabled")
					}
					return nil
				}),

			huh.NewConfirm().
				Title("Add Proxy Headers").
				Description("Add standard proxy headers (X-Real-IP, X-Forwarded-For, etc.)").
				Value(&config.ProxyHeaders),
		).Title("Reverse Proxy Configuration").Description("Configure reverse proxy settings"),

		// Page 4: PHP Configuration
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable PHP Support").
				Description("Enable PHP-FPM processing for .php files").
				Value(&config.EnablePHP),

			huh.NewInput().
				Title("PHP-FPM Socket").
				Description("Path to PHP-FPM socket or TCP address").
				Placeholder("/var/run/php/php-fpm.sock").
				Value(&config.PHPSocket).
				Validate(func(s string) error {
					if config.EnablePHP && s == "" {
						return fmt.Errorf("PHP socket is required when PHP is enabled")
					}
					return nil
				}),
		).Title("PHP Configuration").Description("Configure PHP-FPM settings"),

		// Page 5: Advanced Options
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Gzip Compression").
				Description("Enable gzip compression for responses").
				Value(&config.EnableGzip),

			huh.NewInput().
				Title("Client Max Body Size").
				Description("Maximum allowed size of client request body").
				Placeholder("10M").
				Value(&config.ClientMaxBodySize),

			huh.NewInput().
				Title("Access Log Path").
				Description("Path to access log file").
				Placeholder("/var/log/nginx/access.log").
				Value(&config.AccessLog),

			huh.NewInput().
				Title("Error Log Path").
				Description("Path to error log file").
				Placeholder("/var/log/nginx/error.log").
				Value(&config.ErrorLog),

			huh.NewText().
				Title("Custom Configuration").
				Description("Additional NGINX directives (optional)").
				Placeholder("# Add custom directives here").
				CharLimit(1000).
				Value(&config.CustomConfig),
		).Title("Advanced Options").Description("Configure advanced settings"),

		// Page 6: Confirmation
		huh.NewGroup(
			huh.NewNote().
				Title("Ready to Create Site").
				Description("Review your configuration and confirm to create the NGINX site.\n\n" +
					"The configuration will be tested before being applied.\n" +
					"Press Enter to continue."),

			huh.NewConfirm().
				Title("Create Site?").
				Description("This will create and activate the NGINX site configuration").
				Affirmative("Yes, Create!").
				Negative("No, Cancel").
				Value(&config.Confirmed),
		).Title("Confirm Creation").Description("Final step"),
	)

	return form, config
}

// GenerateNginxConfig generates NGINX configuration from SiteConfig
func (c *SiteConfig) GenerateNginxConfig() string {
	var sb strings.Builder

	// HTTP to HTTPS redirect server block (if SSL and force HTTPS)
	if c.EnableSSL && c.ForceHTTPS {
		sb.WriteString(fmt.Sprintf(`# HTTP to HTTPS redirect
server {
    listen 80;
    server_name %s;
    return 301 https://$server_name$request_uri;
}

`, c.ServerName))
	}

	// Main server block
	sb.WriteString("server {\n")

	// Listen directive
	if c.EnableSSL {
		sb.WriteString(fmt.Sprintf("    listen 443 ssl http2;\n"))
		sb.WriteString(fmt.Sprintf("    listen [::]:443 ssl http2;\n"))
	} else {
		sb.WriteString(fmt.Sprintf("    listen %s;\n", c.Port))
		sb.WriteString(fmt.Sprintf("    listen [::]:%s;\n", c.Port))
	}

	// Server name
	sb.WriteString(fmt.Sprintf("    server_name %s;\n\n", c.ServerName))

	// SSL Configuration
	if c.EnableSSL {
		sb.WriteString(fmt.Sprintf("    # SSL Configuration\n"))
		sb.WriteString(fmt.Sprintf("    ssl_certificate %s;\n", c.SSLCertPath))
		sb.WriteString(fmt.Sprintf("    ssl_certificate_key %s;\n", c.SSLKeyPath))
		sb.WriteString("    ssl_protocols TLSv1.2 TLSv1.3;\n")
		sb.WriteString("    ssl_ciphers HIGH:!aNULL:!MD5;\n")
		sb.WriteString("    ssl_prefer_server_ciphers on;\n\n")
	}

	// Logging
	sb.WriteString(fmt.Sprintf("    # Logging\n"))
	sb.WriteString(fmt.Sprintf("    access_log %s;\n", c.AccessLog))
	sb.WriteString(fmt.Sprintf("    error_log %s;\n\n", c.ErrorLog))

	// Client max body size
	sb.WriteString(fmt.Sprintf("    # Upload size limit\n"))
	sb.WriteString(fmt.Sprintf("    client_max_body_size %s;\n\n", c.ClientMaxBodySize))

	// Gzip
	if c.EnableGzip {
		sb.WriteString("    # Gzip Compression\n")
		sb.WriteString("    gzip on;\n")
		sb.WriteString("    gzip_vary on;\n")
		sb.WriteString("    gzip_proxied any;\n")
		sb.WriteString("    gzip_comp_level 6;\n")
		sb.WriteString("    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss application/rss+xml font/truetype font/opentype application/vnd.ms-fontobject image/svg+xml;\n\n")
	}

	// Root and index (if not proxy)
	if !c.IsProxy {
		sb.WriteString(fmt.Sprintf("    # Document Root\n"))
		sb.WriteString(fmt.Sprintf("    root %s;\n", c.RootPath))
		sb.WriteString(fmt.Sprintf("    index %s;\n\n", c.IndexFiles))
	}

	// Location blocks
	sb.WriteString("    location / {\n")

	if c.IsProxy {
		// Proxy configuration
		sb.WriteString(fmt.Sprintf("        proxy_pass %s;\n", c.ProxyPass))
		if c.ProxyHeaders {
			sb.WriteString("        proxy_set_header Host $host;\n")
			sb.WriteString("        proxy_set_header X-Real-IP $remote_addr;\n")
			sb.WriteString("        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n")
			sb.WriteString("        proxy_set_header X-Forwarded-Proto $scheme;\n")
		}
	} else {
		// Static file serving
		sb.WriteString("        try_files $uri $uri/ =404;\n")
	}

	sb.WriteString("    }\n\n")

	// PHP configuration
	if c.EnablePHP && !c.IsProxy {
		sb.WriteString("    # PHP Configuration\n")
		sb.WriteString("    location ~ \\.php$ {\n")
		sb.WriteString("        include snippets/fastcgi-php.conf;\n")
		
		if strings.HasPrefix(c.PHPSocket, "unix:") || strings.HasPrefix(c.PHPSocket, "/") {
			// Unix socket
			socketPath := strings.TrimPrefix(c.PHPSocket, "unix:")
			sb.WriteString(fmt.Sprintf("        fastcgi_pass unix:%s;\n", socketPath))
		} else {
			// TCP socket
			sb.WriteString(fmt.Sprintf("        fastcgi_pass %s;\n", c.PHPSocket))
		}
		
		sb.WriteString("    }\n\n")
	}

	// Custom configuration
	if c.CustomConfig != "" {
		sb.WriteString("    # Custom Configuration\n")
		// Indent custom config
		lines := strings.Split(c.CustomConfig, "\n")
		for _, line := range lines {
			if line != "" {
				sb.WriteString("    " + line + "\n")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("}\n")

	return sb.String()
}

// GetFileName returns the suggested filename for the site configuration
func (c *SiteConfig) GetFileName() string {
	// Remove protocol and special characters from server name
	name := strings.TrimPrefix(c.ServerName, "http://")
	name = strings.TrimPrefix(name, "https://")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, ":", "_")
	return name
}
