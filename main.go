package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	var timerSeconds, timerMinutes, timerHours int
	var alarmFile string
	var help bool
	flag.IntVar(&timerSeconds, "s", 0, "The number of seconds in the timers total time")
	flag.IntVar(&timerMinutes, "m", 0, "The number of minutes in the timers total time")
	flag.IntVar(&timerHours, "h", 0, "The number of hours in the timers total time")
	flag.StringVar(&alarmFile, "alarmFile", "./alarm.mp3", "The alarm mp3 file to be played at the end of the timer.")
	flag.BoolVar(&help, "help", false, "Print Useage info.")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if timerSeconds == 0 && timerMinutes == 0 && timerHours == 0 {
		fmt.Println("Please set a time for the timer.")
		os.Exit(0)
	}

	var timerLength int
	if timerSeconds != 0 {
		timerLength = timerLength + timerSeconds
	}
	if timerMinutes != 0 {
		timerLength = timerLength + (timerMinutes * 60)
	}
	if timerHours != 0 {
		timerLength = timerLength + ((timerHours * 60) * 60)
	}
	timer(timerLength)
}

func timer(length int) {

	f, err := os.Open("alarm.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	for i := 0; i != length; i++ {
		fmt.Print(counter((length - i)))
		time.Sleep(1 * time.Second)
	}

	for true {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speaker.Play(streamer)
		select {}
	}
}

func counter(in int) string {
	h, m, s := 0.0, 0.0, 0.0

	if in < 60 {
		s = float64(in)
	} else if in < 3600 {
		m = float64(in) / 60
		s = math.Mod(math.Mod(float64(in), 60), 60)
	} else {
		h = float64(in) / 3600
		m = (math.Mod(float64(in), 3600) / 60)
		s = math.Mod(math.Mod(math.Mod(float64(in), 60), 60), 60)
	}

	return fmt.Sprintf("%v:%v:%v\n", int(h), int(m), int(s))
}
