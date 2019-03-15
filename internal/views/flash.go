package views

import (
	"context"
	"strings"
	"time"

	"github.com/derailed/k9s/internal/resource"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

const (
	flashInfo flashLevel = iota
	flashWarn
	flashErr
	flashFatal
	flashDelay = 2

	emoDoh   = "😗"
	emoRed   = "😡"
	emoDead  = "💀"
	emoHappy = "😎"
)

type (
	flashLevel int

	flashView struct {
		*tview.TextView

		cancel context.CancelFunc
		app    *tview.Application
	}
)

func newFlashView(app *tview.Application, m string) *flashView {
	f := flashView{app: app, TextView: tview.NewTextView()}
	{
		f.SetTextColor(tcell.ColorAqua)
		f.SetTextAlign(tview.AlignLeft)
		f.SetBorderPadding(0, 0, 1, 1)
		f.SetText(m)
	}
	return &f
}

func (f *flashView) setMessage(level flashLevel, msg ...string) {
	if f.cancel != nil {
		f.cancel()
	}
	var ctx context.Context
	{
		ctx, f.cancel = context.WithTimeout(context.TODO(), flashDelay*time.Second)
		go func(ctx context.Context) {
			_, _, width, _ := f.GetRect()
			if width <= 15 {
				width = 100
			}
			m := strings.Join(msg, " ")
			f.SetTextColor(flashColor(level))
			f.SetText(resource.Truncate(flashEmoji(level)+" "+m, width-3))
			f.app.Draw()
			for {
				select {
				case <-ctx.Done():
					f.Clear()
					f.app.Draw()
					return
				}
			}
		}(ctx)
	}
}

func flashEmoji(l flashLevel) string {
	switch l {
	case flashWarn:
		return emoDoh
	case flashErr:
		return emoRed
	case flashFatal:
		return emoDead
	default:
		return emoHappy
	}
}

func flashColor(l flashLevel) tcell.Color {
	switch l {
	case flashWarn:
		return tcell.ColorOrange
	case flashErr:
		return tcell.ColorOrangeRed
	case flashFatal:
		return tcell.ColorFuchsia
	default:
		return tcell.ColorNavajoWhite
	}
}
