package utils

import (
	"strings"
)

func UrlStringToOptions(url string) (string, string, string, string, string, string) {
	// parse redis url to username password and host and port and db
	// app_name://username:password@host:port/db

	// split all options (app_name, username, password, host, port, db)
	options := strings.Split(url, "://")
	// split protocol and info
	protocol := options[0]
	info := options[1]
	// split info to username, password, host, port, db
	infoOptions := strings.Split(info, "@")
	// split username and password
	usernamePassword := infoOptions[0]
	hostPortDb := infoOptions[1]
	// split username and password
	usernamePasswordOptions := strings.Split(usernamePassword, ":")
	username := usernamePasswordOptions[0]
	password := usernamePasswordOptions[1]
	// split host, port and db
	hostPortDbOptions := strings.Split(hostPortDb, "/")
	hostPort := hostPortDbOptions[0]
	db := hostPortDbOptions[1]
	// split host and port
	hostPortOptions := strings.Split(hostPort, ":")
	host := hostPortOptions[0]
	port := hostPortOptions[1]

	return protocol, username, password, host, port, db
}
