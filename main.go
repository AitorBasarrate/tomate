package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/schollz/progressbar/v3"
)

func alertMessage(msg string) {
	err := beeep.Alert("TOMATE", msg, "./img/fresh-tomato-vegetable.png")
	if err != nil {
		fmt.Printf("Error sending alert: %s", err)
	}
}

func formatTime(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	seconds -= 60 * minutes
	return fmt.Sprintf("%dmin %ds", minutes, seconds)
}

func countdown(seconds int) {
	bar := progressbar.NewOptions(
		seconds,
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)
	bar.RenderBlank()

	var halfway int
	if seconds >= 120 {
		halfway = seconds / 2
	}

	for i := 0; i < seconds; i++ {
		time.Sleep(1 * time.Second)
		if i == halfway && i > 0 {
			timeLeft := formatTime(i)
			msg := fmt.Sprintf("Keep up! You're halfway - %s", timeLeft)
			alertMessage(msg)
		}
		bar.Add(1)
	}
}

func main() {
	beeep.AppName = "tomate"

	focus_minutes := flag.Int("fm", 0, "Minutos que dura la concentracion")
	focus_seconds := flag.Int("fs", 0, "Segundos que dura la concentracion")
	rest_minutes := flag.Int("dm", 0, "Minutos que dura el descanso")
	rest_seconds := flag.Int("ds", 0, "Segundos que dura el descanso")
	repetitions := flag.Int("r", 1, "Veces que se va a repetir")
	flag.Parse()

	total_focus := *focus_minutes*60 + *focus_seconds
	total_rest := *rest_minutes*60 + *rest_seconds

	if total_focus <= 0 {
		fmt.Println("Introduce una duracion valida")
		return
	}

	i := 1
	for i <= *repetitions {
		alertMessage("Start to focus!")
		countdown(total_focus)
		if total_rest > 0 {
			alertMessage("You can rest now")
			countdown(total_rest)
		}
		i++
	}
	alertMessage("Finished!")
	fmt.Println("")

}
