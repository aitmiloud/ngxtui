package forms

import (
	"fmt"
	"strings"
)

// GetTemplatePreview returns a formatted preview of a template configuration
func GetTemplatePreview(templateName string) string {
	template := GetTemplateByName(templateName)
	if template == nil {
		return "Template not found"
	}

	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("╔═══════════════════════════════════════════════════════════╗\n"))
	sb.WriteString(fmt.Sprintf("║ %s\n", template.Name))
	sb.WriteString(fmt.Sprintf("╠═══════════════════════════════════════════════════════════╣\n"))
	sb.WriteString(fmt.Sprintf("║ Use Case: %s\n", template.UseCase))
	sb.WriteString(fmt.Sprintf("╠═══════════════════════════════════════════════════════════╣\n"))
	sb.WriteString(fmt.Sprintf("║\n"))
	sb.WriteString(fmt.Sprintf("║ 📄 Basic Configuration:\n"))
	sb.WriteString(fmt.Sprintf("║    Server Name: %s\n", template.Config.ServerName))
	sb.WriteString(fmt.Sprintf("║    Port: %s\n", template.Config.Port))
	
	if template.Config.RootPath != "" {
		sb.WriteString(fmt.Sprintf("║    Root Path: %s\n", template.Config.RootPath))
	}
	if template.Config.IndexFiles != "" {
		sb.WriteString(fmt.Sprintf("║    Index Files: %s\n", template.Config.IndexFiles))
	}
	
	sb.WriteString(fmt.Sprintf("║\n"))
	sb.WriteString(fmt.Sprintf("║ 🔒 SSL Configuration:\n"))
	sb.WriteString(fmt.Sprintf("║    Enable SSL: %v\n", template.Config.EnableSSL))
	
	if template.Config.EnableSSL {
		sb.WriteString(fmt.Sprintf("║    SSL Certificate: %s\n", template.Config.SSLCertPath))
		sb.WriteString(fmt.Sprintf("║    SSL Key: %s\n", template.Config.SSLKeyPath))
		sb.WriteString(fmt.Sprintf("║    Force HTTPS: %v\n", template.Config.ForceHTTPS))
	}
	
	sb.WriteString(fmt.Sprintf("║\n"))
	sb.WriteString(fmt.Sprintf("║ 🔄 Proxy Configuration:\n"))
	sb.WriteString(fmt.Sprintf("║    Enable Reverse Proxy: %v\n", template.Config.IsProxy))
	
	if template.Config.IsProxy {
		sb.WriteString(fmt.Sprintf("║    Proxy Pass URL: %s\n", template.Config.ProxyPass))
		sb.WriteString(fmt.Sprintf("║    Add Proxy Headers: %v\n", template.Config.ProxyHeaders))
	}
	
	sb.WriteString(fmt.Sprintf("║\n"))
	sb.WriteString(fmt.Sprintf("║ 🐘 PHP Configuration:\n"))
	sb.WriteString(fmt.Sprintf("║    Enable PHP: %v\n", template.Config.EnablePHP))
	
	if template.Config.EnablePHP {
		sb.WriteString(fmt.Sprintf("║    PHP-FPM Socket: %s\n", template.Config.PHPSocket))
	}
	
	sb.WriteString(fmt.Sprintf("║\n"))
	sb.WriteString(fmt.Sprintf("║ ⚙️ Advanced Options:\n"))
	sb.WriteString(fmt.Sprintf("║    Enable Gzip: %v\n", template.Config.EnableGzip))
	sb.WriteString(fmt.Sprintf("║    Client Max Body Size: %s\n", template.Config.ClientMaxBodySize))
	
	if template.Config.CustomConfig != "" {
		sb.WriteString(fmt.Sprintf("║    Custom Config:\n"))
		lines := strings.Split(template.Config.CustomConfig, "\n")
		for _, line := range lines {
			if line != "" {
				sb.WriteString(fmt.Sprintf("║    %s\n", line))
			}
		}
	}
	
	sb.WriteString(fmt.Sprintf("╚═══════════════════════════════════════════════════════════╝\n"))
	
	return sb.String()
}

// GetTemplatesByCategory returns templates grouped by category
func GetTemplatesByCategory() map[string][]SiteTemplate {
	templates := GetSiteTemplates()
	categories := make(map[string][]SiteTemplate)
	
	// Define categories
	categories["Static & Frontend"] = []SiteTemplate{}
	categories["Backend & APIs"] = []SiteTemplate{}
	categories["PHP Applications"] = []SiteTemplate{}
	categories["Proxy & Containers"] = []SiteTemplate{}
	categories["Other"] = []SiteTemplate{}
	
	// Categorize templates
	for _, t := range templates {
		switch t.Name {
		case "Static Website", "Single Page Application (SPA)":
			categories["Static & Frontend"] = append(categories["Static & Frontend"], t)
		case "Node.js Application", "Python/Django Application", "API Gateway":
			categories["Backend & APIs"] = append(categories["Backend & APIs"], t)
		case "WordPress Site", "Laravel Application":
			categories["PHP Applications"] = append(categories["PHP Applications"], t)
		case "Docker Container Proxy", "WebSocket Application", "Multiple Domain Redirect":
			categories["Proxy & Containers"] = append(categories["Proxy & Containers"], t)
		default:
			categories["Other"] = append(categories["Other"], t)
		}
	}
	
	return categories
}

// GetTemplateRecommendation returns a recommended template based on keywords
func GetTemplateRecommendation(keywords ...string) *SiteTemplate {
	templates := GetSiteTemplates()
	
	for _, keyword := range keywords {
		lowerKeyword := strings.ToLower(keyword)
		
		for _, t := range templates {
			// Check in name, description, and use case
			if strings.Contains(strings.ToLower(t.Name), lowerKeyword) ||
				strings.Contains(strings.ToLower(t.Description), lowerKeyword) ||
				strings.Contains(strings.ToLower(t.UseCase), lowerKeyword) {
				return &t
			}
		}
	}
	
	// Return blank template as fallback
	return GetTemplateByName("Blank Template")
}

// ValidateTemplateConfig validates a template configuration
func ValidateTemplateConfig(config *SiteConfig) []string {
	var errors []string
	
	// Basic validation
	if config.ServerName == "" {
		errors = append(errors, "Server name is required")
	}
	
	if config.Port == "" {
		errors = append(errors, "Port is required")
	}
	
	// SSL validation
	if config.EnableSSL {
		if config.SSLCertPath == "" {
			errors = append(errors, "SSL certificate path is required when SSL is enabled")
		}
		if config.SSLKeyPath == "" {
			errors = append(errors, "SSL key path is required when SSL is enabled")
		}
	}
	
	// Proxy validation
	if config.IsProxy {
		if config.ProxyPass == "" {
			errors = append(errors, "Proxy pass URL is required when reverse proxy is enabled")
		}
	} else {
		// Static file serving validation
		if config.RootPath == "" {
			errors = append(errors, "Root path is required for static file serving")
		}
	}
	
	// PHP validation
	if config.EnablePHP && config.PHPSocket == "" {
		errors = append(errors, "PHP-FPM socket is required when PHP is enabled")
	}
	
	return errors
}

// GetQuickStartGuide returns a quick start guide for a template
func GetQuickStartGuide(templateName string) string {
	template := GetTemplateByName(templateName)
	if template == nil {
		return "Template not found"
	}
	
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("# Quick Start Guide: %s\n\n", template.Name))
	sb.WriteString(fmt.Sprintf("## Overview\n%s\n\n", template.Description))
	sb.WriteString(fmt.Sprintf("## Use Case\n%s\n\n", template.UseCase))
	
	sb.WriteString("## Configuration Steps\n\n")
	
	// Step 1: Basic setup
	sb.WriteString("### 1. Basic Configuration\n")
	sb.WriteString(fmt.Sprintf("- Set your server name (currently: `%s`)\n", template.Config.ServerName))
	sb.WriteString(fmt.Sprintf("- Configure port (currently: `%s`)\n", template.Config.Port))
	
	if template.Config.RootPath != "" {
		sb.WriteString(fmt.Sprintf("- Set document root (currently: `%s`)\n", template.Config.RootPath))
	}
	sb.WriteString("\n")
	
	// Step 2: SSL
	if template.Config.EnableSSL {
		sb.WriteString("### 2. SSL/TLS Configuration\n")
		sb.WriteString("- Ensure SSL certificates are in place\n")
		sb.WriteString(fmt.Sprintf("- Certificate: `%s`\n", template.Config.SSLCertPath))
		sb.WriteString(fmt.Sprintf("- Key: `%s`\n", template.Config.SSLKeyPath))
		sb.WriteString("\n")
	}
	
	// Step 3: Application-specific
	if template.Config.IsProxy {
		sb.WriteString("### 3. Backend Application\n")
		sb.WriteString(fmt.Sprintf("- Ensure your application is running on `%s`\n", template.Config.ProxyPass))
		sb.WriteString("- Verify the application is accessible\n")
		sb.WriteString("\n")
	}
	
	if template.Config.EnablePHP {
		sb.WriteString("### 3. PHP Configuration\n")
		sb.WriteString(fmt.Sprintf("- Ensure PHP-FPM is running on `%s`\n", template.Config.PHPSocket))
		sb.WriteString("- Verify PHP version compatibility\n")
		sb.WriteString("\n")
	}
	
	// Step 4: Testing
	sb.WriteString("### 4. Testing\n")
	sb.WriteString("- Test NGINX configuration: `nginx -t`\n")
	sb.WriteString("- Reload NGINX: `nginx -s reload`\n")
	sb.WriteString(fmt.Sprintf("- Access your site: `http%s://%s`\n", 
		map[bool]string{true: "s", false: ""}[template.Config.EnableSSL],
		template.Config.ServerName))
	
	return sb.String()
}
