package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var timeFile bool

func main() {
	var timerSeconds, timerMinutes, timerHours int
	var alarmFile string
	var help, stopwatch bool
	flag.IntVar(&timerSeconds, "s", 0, "The number of seconds in the timers total time")
	flag.IntVar(&timerMinutes, "m", 0, "The number of minutes in the timers total time")
	flag.IntVar(&timerHours, "h", 0, "The number of hours in the timers total time")
	flag.StringVar(&alarmFile, "alarmFile", "./alarm.mp3", "The alarm mp3 file to be played at the end of the timer.")
	flag.BoolVar(&timeFile, "timeFile", false, "Write time to /tmp/go-timer.")
	flag.BoolVar(&stopwatch, "stopwatch", false, "Makes go timer act as a stopwatch.")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if timerSeconds == 0 && timerMinutes == 0 && timerHours == 0 && !stopwatch {
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
	os.Remove("/tmp/go-timer")
	if !stopwatch {
		timer(timerLength, alarmFile)
	} else {
		stopw()
	}
}

func stopw() {
	var i int
	for true {
		fmt.Print(counter(i))
		time.Sleep(1 * time.Second)
		i++
	}
}

func timer(length int, alarmFile string) {

	f, err := os.Open(alarmFile)
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

	hr := fmt.Sprint(int(h))
	mi := fmt.Sprint(int(m))
	se := fmt.Sprint(int(s))
	if int(h) == 0 {
		hr = "00"
	}
	if int(m) == 0 {
		mi = "00"
	}
	if int(s) == 0 {
		se = "00"
	}
	if int(h) < 10 {
		hr = fmt.Sprint("0", int(h))
	}
	if int(m) < 10 {
		mi = fmt.Sprint("0", int(m))
	}
	if int(s) < 10 {
		se = fmt.Sprint("0", int(s))
	}

	if timeFile {
		file, err := os.Create("/tmp/go-timer")
		if err != nil {
			log.Fatal("Error createing file timeFile", ":", err)
		}
		w := bufio.NewWriter(file)
		dump, err := w.WriteString(fmt.Sprintf("%s:%s:%s\n", hr, mi, se))
		if err != nil {
			log.Fatal(err)
		}
		dump = dump + 1
		w.Flush()
	}

	return fmt.Sprintf("%s:%s:%s\n", hr, mi, se)
}
