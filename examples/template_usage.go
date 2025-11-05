package main

import (
	"fmt"
	"strings"

	"github.com/aitmiloud/ngxtui/internal/forms"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           NgxTUI Template System - Usage Examples             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Example 1: List all templates
	fmt.Println("📋 Available Templates:")
	fmt.Println(strings.Repeat("─", 70))
	templates := forms.GetSiteTemplates()
	for i, t := range templates {
		fmt.Printf("%2d. %-30s - %s\n", i+1, t.Name, t.Description)
	}
	fmt.Println()

	// Example 2: Get templates by category
	fmt.Println("📂 Templates by Category:")
	fmt.Println(strings.Repeat("─", 70))
	categories := forms.GetTemplatesByCategory()
	for category, temps := range categories {
		if len(temps) > 0 {
			fmt.Printf("\n%s:\n", category)
			for _, t := range temps {
				fmt.Printf("  • %s\n", t.Name)
			}
		}
	}
	fmt.Println()

	// Example 3: Get a specific template
	fmt.Println("🔍 Template Details - Node.js Application:")
	fmt.Println(strings.Repeat("─", 70))
	nodeTemplate := forms.GetTemplateByName("Node.js Application")
	if nodeTemplate != nil {
		fmt.Printf("Name: %s\n", nodeTemplate.Name)
		fmt.Printf("Description: %s\n", nodeTemplate.Description)
		fmt.Printf("Use Case: %s\n", nodeTemplate.UseCase)
		fmt.Printf("\nConfiguration:\n")
		fmt.Printf("  Server Name: %s\n", nodeTemplate.Config.ServerName)
		fmt.Printf("  Port: %s\n", nodeTemplate.Config.Port)
		fmt.Printf("  SSL Enabled: %v\n", nodeTemplate.Config.EnableSSL)
		fmt.Printf("  Is Proxy: %v\n", nodeTemplate.Config.IsProxy)
		fmt.Printf("  Proxy Pass: %s\n", nodeTemplate.Config.ProxyPass)
	}
	fmt.Println()

	// Example 4: Template recommendation
	fmt.Println("💡 Template Recommendations:")
	fmt.Println(strings.Repeat("─", 70))
	
	keywords := [][]string{
		{"react", "vue"},
		{"wordpress", "php"},
		{"docker", "container"},
		{"websocket", "realtime"},
	}
	
	for _, kw := range keywords {
		recommended := forms.GetTemplateRecommendation(kw...)
		if recommended != nil {
			fmt.Printf("Keywords: %v → Recommended: %s\n", kw, recommended.Name)
		}
	}
	fmt.Println()

	// Example 5: Template preview
	fmt.Println("👁️  Template Preview - WordPress Site:")
	fmt.Println(strings.Repeat("─", 70))
	preview := forms.GetTemplatePreview("WordPress Site")
	fmt.Println(preview)

	// Example 6: Generate NGINX config
	fmt.Println("⚙️  Generated NGINX Configuration - SPA:")
	fmt.Println(strings.Repeat("─", 70))
	spaTemplate := forms.GetTemplateByName("Single Page Application (SPA)")
	if spaTemplate != nil {
		// Customize the template
		spaTemplate.Config.ServerName = "myapp.example.com"
		spaTemplate.Config.RootPath = "/var/www/myapp"
		
		nginxConfig := spaTemplate.Config.GenerateNginxConfig()
		fmt.Println(nginxConfig)
	}

	// Example 7: Validate configuration
	fmt.Println("✅ Configuration Validation:")
	fmt.Println(strings.Repeat("─", 70))
	
	// Valid config
	validConfig := &forms.SiteConfig{
		ServerName:        "example.com",
		Port:              "443",
		EnableSSL:         true,
		SSLCertPath:       "/etc/ssl/certs/cert.pem",
		SSLKeyPath:        "/etc/ssl/private/key.pem",
		IsProxy:           true,
		ProxyPass:         "http://localhost:3000",
		ClientMaxBodySize: "10M",
	}
	
	errors := forms.ValidateTemplateConfig(validConfig)
	if len(errors) == 0 {
		fmt.Println("✓ Valid configuration - no errors")
	} else {
		fmt.Println("✗ Invalid configuration:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}
	
	// Invalid config
	invalidConfig := &forms.SiteConfig{
		ServerName: "",  // Missing required field
		Port:       "443",
		EnableSSL:  true,
		// Missing SSL paths
	}
	
	errors = forms.ValidateTemplateConfig(invalidConfig)
	if len(errors) > 0 {
		fmt.Println("\n✗ Example invalid configuration:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}
	fmt.Println()

	// Example 8: Quick start guide
	fmt.Println("📖 Quick Start Guide - Laravel Application:")
	fmt.Println(strings.Repeat("─", 70))
	guide := forms.GetQuickStartGuide("Laravel Application")
	fmt.Println(guide)

	// Example 9: Template names
	fmt.Println("📝 All Template Names:")
	fmt.Println(strings.Repeat("─", 70))
	names := forms.GetTemplateNames()
	for i, name := range names {
		fmt.Printf("%2d. %s\n", i+1, name)
	}
	fmt.Println()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    End of Examples                             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}
