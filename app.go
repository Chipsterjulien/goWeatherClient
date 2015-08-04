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

type OnlyTemperature struct {
	Temp float64
}

type Spam struct {
	Str string
}

func processingLines() float64 {
	counter := 0
	linesSplitted := []string{}

	for {
		linesSplitted = readTemporyLines()
		if strings.Contains(linesSplitted[0], "YES") {
			break
		} else {
			counter++
			log := logging.MustGetLogger("log")
			log.Warning("Sensor \"%s\" was not recognize", viper.GetString("device.name"))
			if counter > 4 {
				os.Exit(1)
			}
			time.Sleep(5 * time.Second)
		}
	}

	temperatureString := strings.Split(linesSplitted[1], "t=")[1]
	temperature, err := strconv.ParseFloat(temperatureString, 64)
	if err != nil {
	}

	return float64(int(temperature/100)) / 10
}

func readTemporyLines() []string {
	var lines []byte
	var err error

	for {
		lines, err = ioutil.ReadFile(viper.GetString("device.name"))
		if err == nil {
			break
		} else {
			log := logging.MustGetLogger("log")
			log.Warning("Unable to read \"%s\": %s", viper.GetString("device.name"), err)
			time.Sleep(5 * time.Second)
		}
	}

	return strings.Split(string(lines), "\n")
}

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
	}
}

func app() {
	temp := processingLines()
	sendTemperature(temp)
}
