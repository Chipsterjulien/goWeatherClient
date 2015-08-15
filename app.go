package main

import (
	"github.com/jmcvetta/napping"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// temperature will be sent in json format
type OnlyTemperature struct {
	Temp float64
}

// Response of server after sending temperature
type Spam struct {
	Str string
}

// get and convert temperature string format to a float64
func processingLines() float64 {
	log := logging.MustGetLogger("log")
	counter := 0
	linesSplitted := []string{}
	log.Debug("Counter val before beginning: %d", counter)

	for {
		linesSplitted = readTemporyLines()
		log.Debug("Line was reading on device: %s", strings.Join(linesSplitted, "\n"))
		if strings.Contains(linesSplitted[0], "YES") {
			log.Debug("  It contain Yes")
			break
		} else {
			log.Debug("Counter val: %d", counter)
			counter++
			log.Warning("Sensor \"%s\" was not recognize", viper.GetString("device.name"))
			if counter > 4 {
				os.Exit(1)
			}
			time.Sleep(5 * time.Second)
		}
	}

	temperatureString := strings.Split(linesSplitted[1], "t=")[1]

	log.Debug("Temperature string is: %s", temperatureString)

	temperature, err := strconv.ParseFloat(temperatureString, 64)
	if err != nil {
		log.Critical("Unable to convert \"%s\" to a float !", temperatureString)
		os.Exit(1)
	}

	log.Debug("Temperature float is: %f", temperature)

	return float64(int(temperature/100)) / 10
}

// Read temperature on device
func readTemporyLines() []string {
	log := logging.MustGetLogger("log")
	counter := 0
	var lines []byte
	var err error

	for {
		lines, err = ioutil.ReadFile(viper.GetString("device.name"))
		if err == nil {
			log.Debug("Lines who was reading are: %s", lines)
			break
		} else {
			log.Debug("Counter val: %d", counter)
			counter++
			log.Warning("Unable to read \"%s\": %s", viper.GetString("device.name"), err)
			if counter > 4 {
				os.Exit(1)
			}
			time.Sleep(5 * time.Second)
		}
	}

	log.Debug("Lines will be return: %s", string(lines))

	return strings.Split(string(lines), "\n")
}

// Send temperature (float64) via tcp on json format
func sendTemperature(temp float64) {
	log := logging.MustGetLogger("log")
	data := OnlyTemperature{
		Temp: temp,
	}
	str := Spam{}

	resp, err := napping.Post(viper.GetString("server.url"), &data, &str, nil)
	if err != nil {
		log.Critical("Unable to send data:", err)
	}
	if resp.Status() != 200 {
		log.Critical("Status is not good: \"%d\". Error is: %s", resp.Status(), str.Str)
		os.Exit(1)
	}
}

func app() {
	temp := processingLines()
	sendTemperature(temp)
}
