package main

func main() {
	logFilename := "/var/log/goweatherclient/errors.log"
	confPath := "/etc/goweatherclient"
	confFilename := "goweatherclient"

	fd := initLogging(&logFilename)
	defer fd.Close()

	loadConfig(&confPath, &confFilename)
	app()
}
