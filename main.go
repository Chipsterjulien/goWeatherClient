package main

func main() {
	logFilename := "/var/log/goWeatherClient/errors.log"
	confPath := "/etc/goWeatherClient"
	confFilename := "goWeatherClient"

	/*
		logFilename := "errors.log"
		confFilename := "goWeatherClient"
		confPath := "."
	*/

	fd := initLogging(&logFilename)
	defer fd.Close()

	loadConfig(&confPath, &confFilename)
	app()
}
