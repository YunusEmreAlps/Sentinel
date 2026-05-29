package constants

import "time"

const (
	// HTTP Server timeouts
	ServerReadTimeout    = 10 * time.Second
	ServerWriteTimeout   = 10 * time.Second
	ServerIdleTimeout    = 120 * time.Second
	ServerMaxHeaderBytes = 1 << 20 // 1 MB

	// Shutdown timeout
	ShutdownTimeout = 30 * time.Second

	// Certificate check defaults
	DefaultExpireDays = 30
	DefaultBufferSize = 1024

	// HTTP request headers
	HTTPVersion    = "HTTP/1.1"
	HTTPGetRequest = "GET / HTTP/1.1\r\n"
	HTTPHostHeader = "Host: "
	HTTPLineEnd    = "\r\n"

	// Certificate status
	CertStatusExpired    = 1
	CertStatusNotExpired = 0
	CertStatusTimeout    = 2

	// Excel colors
	ExcelColorRed    = "#FF0000"
	ExcelColorGreen  = "#00FF00"
	ExcelColorYellow = "#FFFF00"

	// Port parsing
	DefaultPortIndex = 1
)
