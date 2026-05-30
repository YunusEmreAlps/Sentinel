package utils

var DomainList = []string{
	// domain list here - supports multiple formats:
	// 1. domain.com (automatically adds port 443)
	// 2. domain.com:443 (explicit port)
	// 3. https://domain.com (uses port 443)
	// 4. http://domain.com (uses port 80)
	// 5. https://domain.com:8443 (custom HTTPS port)

	"yunusemrealpu.netlify.app:443", // format with port
	// "example.com",                  // without port (will use 443)
	// "https://google.com",           // with protocol (will use 443)
	// "http://example.com:8080",      // HTTP with custom port
}

// To Users
var ToUsers = []string{
	// team members here
}

// CC Users
var CCUsers = []string{
	// team members here
}
