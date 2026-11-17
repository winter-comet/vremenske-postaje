package vreme

import (
	"fmt"
)

type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura float64, minTemp float64, maxTemp float64) {
	_, obstaja := postaje[oznaka]
	if !obstaja {
		fmt.Printf("Zapis o postaji \"%s\" ne obstaja\n", oznaka)
		return
	}

	if temperatura < minTemp || temperatura > maxTemp {
		fmt.Printf("Temperatura %.1f ne leži na intervalu [%g, %g]\n", temperatura, minTemp, maxTemp)
		return
	}

	postaja := postaje[oznaka]
	postaja.meritve = append(postaja.meritve, temperatura)
	postaje[oznaka] = postaja
}

func DodajPostajo(postaje map[string]Postaja, oznaka string, ime string, kraj string) {
	if _, obstaja := postaje[oznaka]; obstaja {
		fmt.Printf("Zapis o postaji \"%s\" že obstaja\n", oznaka)
		return
	}

	postaje[oznaka] = Postaja{
		ime:     ime,
		kraj:    kraj,
		meritve: []float64{},
	}
}

func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for ključ, vrednost := range postaje {
		fmt.Printf("%s - %s (%s): %v\n", ključ, vrednost.ime, vrednost.kraj, vrednost.meritve)
	}
}

func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	fmt.Println("POROČILO:")
	for ključ, vrednost := range postaje {
		najnizja, najvisja, povprecna, ok := statistika(postaje, ključ, stMeritev)

		if !ok {
			fmt.Printf("%s - %s (%s): premalo meritev\n", ključ, vrednost.ime, vrednost.kraj)
		} else {
			var opis string

			switch {
			case povprecna >= 25:
				opis = "Vroče"
			case povprecna >= 15:
				opis = "Toplo"
			case povprecna >= 5:
				opis = "Hladno"
			default:
				opis = "Mrzlo"
			}

			fmt.Printf("%s - %s (%s): min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", ključ, vrednost.ime, vrednost.kraj, najnizja, najvisja, povprecna, opis)

		}
	}
}

func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	_, obstaja := postaje[oznaka]
	if !obstaja {
		fmt.Printf("Zapis o postaji %s ne obstaja\n", oznaka)
		ok = false
		return
	}

	if len(postaje[oznaka].meritve) < 1 {
		ok = false
		return
	}

	najnizja = postaje[oznaka].meritve[0]
	najvisja = postaje[oznaka].meritve[0]
	povprecna = 0
	for _, vrednost := range postaje[oznaka].meritve {
		if vrednost < najnizja {
			najnizja = vrednost
		}
		if vrednost > najvisja {
			najvisja = vrednost
		}
		povprecna += vrednost / float64(len(postaje[oznaka].meritve))
	}

	ok = len(postaje[oznaka].meritve) >= stMeritev
	return
}
