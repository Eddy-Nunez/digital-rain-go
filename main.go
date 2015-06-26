package main

import (
	"github.com/nsf/termbox-go"
    "math/rand"
    "time"
    "unicode/utf8"
	"log"
	"os"
)

var (
	fd      *os.File
	logfile *log.Logger
	rmap    map[int]rune
)

const (
	jrunes = "EvXw&()*49-+=#あぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ゘ゝゞゟ"
)

type positionedCell struct {
	posX, posY int
	cell       termbox.Cell
}

func init() {
	i := 0
	rmap = make(map[int]rune, utf8.RuneCountInString(jrunes))
	for _, r := range jrunes {
		rmap[i] = r
		i++
	}
	rand.Seed(time.Now().UnixNano())
}

func randRune(runes string) rune {
	runeCount := utf8.RuneCountInString(runes)
	return rmap[rand.Intn(runeCount)]
}

// consumer of positionedCell types
func renderCell(renderChannel <-chan positionedCell) {
	for ch := range renderChannel {
		termbox.SetCell(ch.posX, ch.posY, ch.cell.Ch, ch.cell.Fg, ch.cell.Bg)
		termbox.Flush()
	}
}

// producer of positionedCells
func generateCells(sleepTime int, renderChannel chan<- positionedCell) {
	x, y := termbox.Size()
	midline := int(float32(y)*0.75 + 1)
	for {
		tx, ty := rand.Intn(x), rand.Intn(midline)
		fg, bg := termbox.ColorGreen, termbox.ColorDefault
		for i, dy := rand.Intn(y-ty), 0; i > 0; i-- {
			switch v := rand.Intn(100) % 15; v {
			case 0:
				fg = termbox.ColorGreen | termbox.AttrBold
				bg = termbox.ColorDefault
			case 1:
				fg = termbox.ColorDefault
				bg = termbox.ColorGreen
			}
			if dy++; ty+dy <= y {
				cell := termbox.Cell{randRune(jrunes), fg, bg}
				//logfile.Printf("tx=%d, ty=%d, dy=%d, chan len=%d\n", tx, ty, dy, len(renderChannel))
				renderChannel <- positionedCell{tx, ty + dy, cell}
				time.Sleep(time.Duration(sleepTime) * time.Millisecond)
			}
		}
	}
}

// thread that watches for Ctrl-Q quit event
func checkForQuit(ch chan<- bool) {
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				switch ev.Key {
					case termbox.KeyCtrlQ, termbox.KeyEsc:
						ch <- true
					default:
						// noop
				}
			}
		}
	}()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// fd, _   = os.Create("log.txt")
	// logfile = log.New(fd, "", log.Lshortfile)
	// defer fd.Close()

	termbox.SetInputMode(termbox.InputEsc)

	quitChan   := make(chan bool)
	renderChan := make(chan positionedCell, 100)
    // launch goroutines
    go checkForQuit(quitChan)
	go renderCell(renderChan)
	for threads := 20; threads > 0; threads-- {
		go generateCells(threads*10, renderChan)
	}

loop:
	for {
		select {
		case <-quitChan:
			break loop
		default:
			// noop
		}
	}
}
