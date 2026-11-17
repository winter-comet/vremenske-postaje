package vreme

import (
	"fmt"
)

func dodajMeritev(postaje map[string]Postaja, oznaka string, temperatura float64) {
	_, obstaja := postaje[oznaka]
	if !obstaja {
		fmt.Printf("Zapis o postaji \"%s\" ne obstaja\n", oznaka)
		return
	}

	if temperatura < -50 || temperatura > 50 {
		fmt.Printf("Temperatura %.1f ne leži na intervalu [-50, 50]\n", temperatura)
		return
	}

	postaja := postaje[oznaka]
	postaja.meritve = append(postaja.meritve, temperatura)
	postaje[oznaka] = postaja
}

func statistika(postaje map[string]Postaja, oznaka string) (najnizja, najvisja, povprecna float64, ok bool) {
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

	ok = len(postaje[oznaka].meritve) >= 3
	return
}

func izpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for ključ, vrednost := range postaje {
		fmt.Printf("%s - %s (%s): %v\n", ključ, vrednost.ime, vrednost.kraj, vrednost.meritve)
	}
}

func izpisMeritve(postaje map[string]Postaja, oznaka string) {
	_, obstaja := postaje[oznaka]
	if !obstaja {
		fmt.Printf("Zapis o postaji \"%s\" ne obstaja\n", oznaka)
		return
	}

	fmt.Printf("%s - %s (%s): %v\n", oznaka, postaje[oznaka].ime, postaje[oznaka].kraj, postaje[oznaka].meritve)
}

func izpisPorocila(postaje map[string]Postaja) {
	fmt.Println("POROČILO:")
	for ključ, vrednost := range postaje {
		najnizja, najvisja, povprecna, ok := statistika(postaje, ključ)

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

type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

func main() {
	var postaje = make(map[string]Postaja)

	// Začetni podatki v slovarju
	postaje["LJ01"] = Postaja{
		"Ljubljana",
		"Ljubljana",
		[]float64{9.3, 9.6, 9.7, 9.4},
	}

	postaje["LJ02"] = Postaja{
		"Letališče Jožeta Pučnika Ljubljana",
		"Ljubljana",
		[]float64{8.5, 8.9},
	}

	postaje["KP01"] = Postaja{
		"Koper Kapitanija",
		"Koper",
		[]float64{13.6, 14.1, 14.0, 14.3, 15.1},
	}

	postaje["KP02"] = Postaja{
		"Oceanografska boja Piran",
		"Piran",
		[]float64{15.9, 15.8, 15.7, 13.7},
	}

	postaje["KR01"] = Postaja{
		"Kredarica",
		"Radovljica",
		[]float64{-6.9, -7.3, -7.4},
	}

	// Klici funkcij
	fmt.Println("-> ZAČETNO STANJE:")
	izpisMeritev(postaje)

	fmt.Println()
	fmt.Println("-> DODAJANJE NEVELJAVNE MERITVE:")
	dodajMeritev(postaje, "LJ01", 60.0)

	fmt.Println()
	fmt.Println("-> DODAJANJE MERITVE NEOBSTOJEČI POSTAJI:")
	dodajMeritev(postaje, "LJ05", 9.9)

	fmt.Println()
	fmt.Println("-> MERITVE \"LJ01\" PRED IN PO DODAJANJU VELJAVNE MERITVE:")
	izpisMeritve(postaje, "LJ01")
	dodajMeritev(postaje, "LJ01", 9.9)
	izpisMeritve(postaje, "LJ01")

	fmt.Println()
	fmt.Println("-> STATISTIKA POSTAJE \"LJ01\":")
	izpisMeritve(postaje, "LJ01")
	najnizja, najvisja, povprecna, ok := statistika(postaje, "LJ01")
	fmt.Printf("ok: %v\n", ok)
	fmt.Printf("min: %.1f\n", najnizja)
	fmt.Printf("max: %.1f\n", najvisja)
	fmt.Printf("povprecje: %.1f\n", povprecna)

	fmt.Println()
	fmt.Println("-> IZPIS POROČILA:")
	izpisPorocila(postaje)
}
