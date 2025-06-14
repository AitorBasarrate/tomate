package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

func countdown(seconds int) {
	bar := progressbar.NewOptions(
		seconds,
		progressbar.OptionShowBytes(false),
	)
	for i := 0; i <= seconds; i++ {
		bar.Set(i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
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
		fmt.Println("")
		fmt.Println("Start to focus!")
		countdown(total_focus)
		if total_rest > 0 {
			fmt.Println("")
			fmt.Println("You can rest now")
			countdown(total_rest)
		}
		i++
	}
	fmt.Println("")
	fmt.Println("Finished")
}
