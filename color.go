package main

type TermColorFunc func(text string) string

type TermColors struct {
	Dim   TermColorFunc
	Bold  TermColorFunc
	Green TermColorFunc
	Red   TermColorFunc
	Reset TermColorFunc
}

func colorBuilder(ansiiOpen string, ansiiClose string) func(text string) string {
	return func(text string) string {
		return "\x1b[" + ansiiOpen + "m" + text + "\x1b[" + ansiiClose + "m"
	}
}

func (t *TermColors) Init() {
	t.Dim = colorBuilder("2", "22")
	t.Bold = colorBuilder("1", "22")
	t.Green = colorBuilder("32", "39")
	t.Red = colorBuilder("31", "39")
	t.Reset = colorBuilder("0", "0")
}
