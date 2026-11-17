// Package vreme beleži temperaturne meritve vremenskih postaj in izpisuje
// njihovo statistiko.
//
// Postaje hranimo v slovarju, v katerem je ključ oznaka postaje (npr. "LJ01"),
// vrednost pa [Postaja]. Ker so polja strukture Postaja neizvožena, postajo
// ustvarimo s funkcijo [DodajPostajo], meritve pa ji dodajamo s funkcijo
// [DodajMeritev]:
//
//	postaje := make(map[string]vreme.Postaja)
//	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Bežigrad", "Ljubljana")
//	vreme.DodajMeritev(postaje, "LJ01", 18.2, -50, 50)
//	vreme.IzpisMeritev(postaje)
//	vreme.IzpisPorocila(postaje, 3)
package vreme

import (
	"fmt"
)

// Postaja je vremenska postaja z imenom, krajem in izmerjenimi temperaturami
// v °C. Novo postajo ustvarimo s funkcijo [DodajPostajo].
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajMeritev postaji z dano oznako doda meritev temperature v °C. Meritev
// je veljavna, če leži na intervalu [minTemp, maxTemp].
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

// DodajPostajo v slovar postaje doda novo postajo brez meritev z dano oznako,
// imenom in krajem. Obstoječe postaje z isto oznako ne spremeni.
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

// IzpisMeritev za vsako postajo izpiše oznako, ime, kraj in vse meritve.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for ključ, vrednost := range postaje {
		fmt.Printf("%s - %s (%s): %v\n", ključ, vrednost.ime, vrednost.kraj, vrednost.meritve)
	}
}

// IzpisPorocila za vsako postajo izpiše najnižjo, najvišjo in povprečno
// temperaturo ter opis glede na povprečje:
//
//   - Vroče, če je povprečje vsaj 25 °C,
//   - Toplo, če je med 15 °C in 25 °C,
//   - Hladno, če je med 5 °C in 15 °C,
//   - Mrzlo, če je pod 5 °C.
//
// Za postajo z manj kot stMeritev meritvami izpiše "premalo meritev".
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

// statistika vrne najnižjo, najvišjo in povprečno temperaturo postaje z dano
// oznako. Če postaja ne obstaja ali ima manj kot stMeritev meritev, vrne
// ok == false.
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
