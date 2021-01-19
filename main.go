package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

type baitStruct struct {
	name       string
	totalBaits int
}

//ders id
// 76561198128945703
func isBaitingFile(fd string, id int64, topBaiters bool) {
	f, err := os.Open(fd)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()
	totalBaits := 0
	pname := ""

	var tb map[int64]*baitStruct
	tb = make(map[int64]*baitStruct)

	p.RegisterEventHandler(func(e events.MatchStart) {
		fmt.Println("======PLAYERS======")
		for _, ply := range p.GameState().Participants().All() {
			if topBaiters {
				tb[int64(ply.SteamID64)] = &baitStruct{name: ply.Name, totalBaits: 0}
				fmt.Printf("[%v] %v\n", ply.SteamID64, tb[int64(ply.SteamID64)].name)
			} else {
				fmt.Printf("[%v] %v\n", ply.SteamID64, ply.Name)
			}
		}
		fmt.Println("===================")
	})

	p.RegisterEventHandler(func(e events.Kill) {
		if p.GameState().IsMatchStarted() {
			// fmt.Printf("[round #%v] %v killed %v\n", p.GameState().TotalRoundsPlayed(), e.Killer.Name, e.Victim.Name)
			if e.Victim.SteamID64 == uint64(id) || topBaiters {
				pname = e.Victim.Name
				// Calculate people alive on team
				aliveteam := 0
				for _, tm := range e.Victim.TeamState.Members() {
					if tm.LastAlivePosition != tm.Position() {
						aliveteam = aliveteam + 1

					}
				}
				if aliveteam == 1 {
					// println("Potential Baiter: " + e.Victim.Name + " at " + fmt.Sprint(p.GameState().TotalRoundsPlayed()))
					totalBaits = totalBaits + 1
					if topBaiters {
						tb[int64(e.Victim.SteamID64)].totalBaits = (tb[int64(e.Victim.SteamID64)].totalBaits) + 1
					}

				}

			}
		}
	})

	p.RegisterEventHandler(func(e events.AnnouncementWinPanelMatch) {
		if !topBaiters {
			roundsPlayed := float32(p.GameState().TotalRoundsPlayed())
			fTotalBaits := float32(totalBaits)
			fmt.Println("======BAIT CALC======")
			fmt.Printf("Player: %v(%v) on %v\n", pname, id, p.Header().MapName)
			fmt.Printf("Baited %v/%v = %v percent of rounds\n", totalBaits, p.GameState().TotalRoundsPlayed(), (fTotalBaits/roundsPlayed)*100)
			fmt.Println("=====================")
		} else {
			fmt.Println("======BAIT CALC======")
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 30, 8, 0, '\t', tabwriter.AlignRight|tabwriter.Debug)
			fmt.Fprint(w, "id\tname\tpercent of rounds baited\n")
			for _, ply := range p.GameState().Participants().All() {
				roundsPlayed := float32(p.GameState().TotalRoundsPlayed())
				fTotalBaits := float32(tb[int64(ply.SteamID64)].totalBaits)
				fmt.Fprintf(w, "%v\t%v\t%v\n", ply.Name, ply.SteamID64, (fTotalBaits/roundsPlayed)*100)
				// fmt.Printf("Player: %v(%v) ", ply.Name, ply.SteamID64)
				// fmt.Printf("Baited %v/%v = %v percent of rounds\n", tb[int64(ply.SteamID64)].totalBaits, p.GameState().TotalRoundsPlayed(), (fTotalBaits/roundsPlayed)*100)
				w.Flush()
			}
			fmt.Println("=====================")

		}
	})
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}

}

func main() {
	println("Is He Baiting For Frags? - The Age Old Question")
	demoFile := flag.String("demofile", "test.dem", "The demo file to parse")
	steamID := flag.Int64("steamId", 76561198128945703, "SteamId64 of player to Watch")
	topBaiters := flag.Bool("topBaiters", false, "Analyze every player")
	flag.Parse()

	isBaitingFile(*demoFile, *steamID, *topBaiters)

}
