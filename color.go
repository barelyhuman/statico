package main

type TermColor struct {
	AnsiStart string
	AnsiEnd   string
}

type TermColorBuilder struct {
	TermColor
}

func (tb *TermColorBuilder) Build(text string) string {
	return tb.AnsiStart + text + tb.AnsiEnd
}

func (tb *TermColorBuilder) Dim() *TermColorBuilder {
	start, end := colorBuilder("2", "22")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}
func (tb *TermColorBuilder) Bold() *TermColorBuilder {
	start, end := colorBuilder("1", "22")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}
func (tb *TermColorBuilder) Green() *TermColorBuilder {
	start, end := colorBuilder("32", "39")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}

func (tb *TermColorBuilder) Yellow() *TermColorBuilder {
	start, end := colorBuilder("33", "39")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}
func (tb *TermColorBuilder) Red() *TermColorBuilder {
	start, end := colorBuilder("31", "39")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}

func (tb *TermColorBuilder) Reset() *TermColorBuilder {
	start, end := colorBuilder("0", "0")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}

func (tb *TermColorBuilder) Cyan() *TermColorBuilder {
	start, end := colorBuilder("36", "39")
	tb.AnsiStart += start
	tb.AnsiEnd += end
	return tb
}

func colorBuilder(ansiiOpen string, ansiiClose string) (string, string) {
	return "\x1b[" + ansiiOpen + "m", "\x1b[" + ansiiClose + "m"
}
