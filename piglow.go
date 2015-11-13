package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/schoentoon/piglow"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GlowOff() {
	piglow.ShutDown()
}

func SetFan(res http.ResponseWriter, req *http.Request, parms martini.Params) (int, string) {
	numstr := parms["num"]
	num, _ := strconv.ParseInt(numstr, 10, 64)

	fmt.Println("Going to set fan!")
	
	intervalstr := req.FormValue(FORM_INTERVAL)
	if intervalstr == "" {
		intervalstr = "500"
	}
	interval, _ := strconv.ParseInt(intervalstr, 10, 64)
	
	GlowOff();
	
	// piglow.ShutDown()
	for i := 0; i < int(num); i++ {
		// turn off all of the lights
		GlowOff()
		setLegOn(i%3, 0.3)
		// Sleep for a while
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	
	return 200, "set fan " + numstr + "\n"
}

func SetFade(res http.ResponseWriter, req *http.Request, parms martini.Params) (int, string) {
	numstr := parms["num"]
	num, _ := strconv.ParseInt(numstr, 10, 64)

	fmt.Println("Going to set fan!")
	
	intensity := req.FormValue(FORM_INTENSITY)
	if intensity == "" {
		intensity = "0.5"
	}	
	brightness, _ := strconv.ParseFloat(intensity, 64)
	
	intervalstr := req.FormValue(FORM_INTERVAL)
	if intervalstr == "" {
		intervalstr = "500"
	}
	interval, _ := strconv.ParseInt(intervalstr, 10, 64)
		
	colors := []string{"white", "blue", "green", "yellow", "orange", "red"}
	
	GlowOff();
	
	// piglow.ShutDown()
	for i := 0; i < int(num); i++ {
		// turn off all of the lights
		
		for j := 0; j < 6; j ++ {
			err := glowToColor(getColorFromString(colors[j]), brightness)
			
			if err != nil {
				return 400, err.Error()
			}
			
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
		
		// Sleep for a while
		time.Sleep(time.Duration(interval) * time.Millisecond)
		GlowOff();
	}
	
	return 200, "set fade " + numstr + " to intensity " + intensity + "\n"
}

func checkPiGlow() bool {
	if !piglow.HasPiGlow() {
		fmt.Println("piglow is not connected")
		return false
	}
	return true
}

func glowToColor(color byte, brightness float64) error {
	brightnessByte := byte(math.Floor(brightness*255.0 + 0.5))
	return piglow.Ring(color, brightnessByte)

}

func setLegOn(leg int, brightness float64) error {
	return piglow.Leg(byte(leg), byte(brightness*255+0.5))

}

func setLedOn(leg int, color string, brightness float64) error {
	return piglow.Led(byte(leg), getColorFromString(color), byte(brightness*255+0.5))
}

func TurnAllOn() {
	for i := 0; i < 3; i++ {
		setLegOn(i, 0.3)
	}
}

func Flare() {
	for i := 0; i < 3; i++ {
		setLegOn(i, 1)
	}
}

func getColorFromString(color string) byte {
	if strings.Contains(color, "red") {
		return piglow.Red
	}
	if strings.Contains(color, "orange") {
		return piglow.Orange
	}
	if strings.Contains(color, "yellow") {
		return piglow.Yellow
	}
	if strings.Contains(color, "green") {
		return piglow.Green
	}
	if strings.Contains(color, "blue") {
		return piglow.Blue
	}
	if strings.Contains(color, "white") {
		return piglow.White
	}
	return 0x00
}

func SetLed(res http.ResponseWriter, req *http.Request, parms martini.Params) (int, string) {
	legId := parms["id"]
	colId := parms["colid"]
	intensity := req.FormValue(FORM_INTENSITY)
	if intensity == "" {
		intensity = "0.2"
	}
	legIdNum, _ := strconv.ParseInt(legId, 10, 64)
	brightness, _ := strconv.ParseFloat(intensity, 64)
	err := setLedOn(int(legIdNum), colId, brightness)
	if err != nil {
		return 400, err.Error()
	}
	return 200, "set leg " + legId + " to intensity " + intensity + "\n"
}

func SetLeg(res http.ResponseWriter, req *http.Request, parms martini.Params) (int, string) {
	legId := parms["id"]
	intensity := req.FormValue(FORM_INTENSITY)
	if intensity == "" {
		intensity = "0.2"
	}
	legIdNum, _ := strconv.ParseInt(legId, 10, 64)
	brightness, _ := strconv.ParseFloat(intensity, 64)
	err := setLegOn(int(legIdNum), brightness)
	if err != nil {
		return 400, err.Error()
	}
	return 200, "set leg " + legId + " to intensity " + intensity + "\n"
}

func SetGlowColor(res http.ResponseWriter, req *http.Request, parms martini.Params) (int, string) {
	colorString := parms["id"]
	intensity := req.FormValue(FORM_INTENSITY)
	if intensity == "" {
		intensity = "0.5"
	}
	brightness, _ := strconv.ParseFloat(intensity, 64)
	err := glowToColor(getColorFromString(colorString), brightness)
	if err != nil {
		return 400, err.Error()
	}
	return 200, "set color " + colorString + " to On\n"
}
