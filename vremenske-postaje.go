package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/winter-comet/vremenske-postaje/vreme"
)

func main() {
	cmd := &cli.Command{
		Name:  "vremenskaepostaje",
		Usage: "beleži meritve vremenskih postaj",
		Flags: []cli.Flag{
			&cli.FloatFlag{
				Name:  "minTemp",
				Usage: "najnižja dovoljena temperatura meritve v °C",
				Value: -50,
			},
			&cli.FloatFlag{
				Name:  "maxTemp",
				Usage: "Najvišja dovoljena temperatura meritve v °C",
				Value: 50,
			},
			&cli.IntFlag{
				Name:  "stMeritev",
				Usage: "najmanjše število meritev, potrebnih za izračun statistike posamezne postaje",
				Value: 3,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return zazeni(cmd.Float("minTemp"), cmd.Float("maxTemp"), cmd.Int("stMeritev"))
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func zazeni(minTemp, maxTemp float64, stMeritev int) error {
	if minTemp > maxTemp {
		return fmt.Errorf("minTemp (%g) ne sme biti večja od maxTemp (%g)", minTemp, maxTemp)
	}
	if stMeritev < 1 {
		return fmt.Errorf("stMeritev mora biti vsaj 1, podano je %d", stMeritev)
	}

	postaje := make(map[string]vreme.Postaja)

	fmt.Println("-> DODAJANJE ZAČETNIH PODATKOV:")
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana", "Ljubljana")
	vreme.DodajPostajo(postaje, "LJ02", "Letališče Jožeta Pučnika Ljubljana", "Ljubljana")
	vreme.DodajPostajo(postaje, "KP01", "Koper Kapitanija", "Koper")
	vreme.DodajPostajo(postaje, "KP02", "Oceanografska boja Piran", "Piran")
	vreme.DodajPostajo(postaje, "KR01", "Kredarica", "Radovljica")

	zacetneMeritve := map[string][]float64{
		"LJ01": {9.3, 9.6, 9.7, 9.4},
		"LJ02": {8.5, 8.9},
		"KP01": {13.6, 14.1, 14.0, 14.3, 15.1},
		"KP02": {15.9, 15.8, 15.7, 13.7},
		"KR01": {-6.9, -7.3, -7.4},
	}
	for oznaka, meritve := range zacetneMeritve {
		for _, temperatura := range meritve {
			vreme.DodajMeritev(postaje, oznaka, temperatura, minTemp, maxTemp)
		}
	}

	fmt.Println()
	fmt.Println("-> ZAČETNO STANJE:")
	vreme.IzpisMeritev(postaje)

	fmt.Println()
	fmt.Println("-> DODAJANJE NEVELJAVNE MERITVE:")
	vreme.DodajMeritev(postaje, "LJ01", maxTemp+10, minTemp, maxTemp) // vedno izven območja

	fmt.Println()
	fmt.Println("-> DODAJANJE MERITVE NEOBSTOJEČI POSTAJI:")
	vreme.DodajMeritev(postaje, "LJ05", 9.9, minTemp, maxTemp)

	fmt.Println()
	fmt.Println("-> DODAJANJE OBSTOJEČE POSTAJE:")
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Bežigrad", "Ljubljana")

	fmt.Println()
	fmt.Println("-> DODAJANJE VELJAVNE MERITVE POSTAJI \"LJ01\":")
	vreme.DodajMeritev(postaje, "LJ01", 9.9, minTemp, maxTemp)
	vreme.IzpisMeritev(postaje)

	fmt.Println()
	fmt.Println("-> IZPIS POROČILA:")
	vreme.IzpisPorocila(postaje, stMeritev)

	return nil
}
