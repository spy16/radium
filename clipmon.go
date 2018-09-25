package radium

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

// ClipboardMonitor monitors the system clipboard and tries to
// use clipboard content as queries to radium.
// If the number of words in clipboard content is more than the
// maxWords, then ClipboardMonitor will not perform a radium
// query with it.
type ClipboardMonitor struct {
	MaxWords int
	Interval time.Duration
	Instance *Instance

	oldContent string
}

// Run starts an infinite for loop which continuously monitors
// the system clipboard for changes. Run() can be invoked as
// a goroutine, and a context can be passed in for stopping
// the monitor.
func (cbm *ClipboardMonitor) Run(ctx context.Context) error {
	if err := cbm.ensureInitialized(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			cbm.Instance.Infof("context cancellation received")
			return nil
		default:
			currentContent, _ := clipboard.ReadAll()
			if currentContent == cbm.oldContent {
				continue
			}

			cbm.Instance.Infof("new clipboard content received")
			cbm.oldContent = currentContent

			words := strings.Split(currentContent, " ")
			if len(words) <= cbm.MaxWords {
				query := Query{}
				query.Text = strings.TrimSpace(currentContent)

				cbm.Instance.Infof("running query for '%s'..", query.Text)
				rs, err := cbm.Instance.Search(ctx, query, Strategy1st)
				if err == nil && rs != nil {
					if len(rs) > 0 {
						cbm.Instance.Infof("received result. pasting back..")
						clipboard.WriteAll(rs[0].Content)
						cbm.oldContent = rs[0].Content
					} else {
						cbm.Instance.Infof("no results found")
					}
				}
			}
		}
		time.Sleep(cbm.Interval)
	}
}

func (cbm *ClipboardMonitor) ensureInitialized() error {
	if cbm.MaxWords <= 0 {
		cbm.MaxWords = 5
	}

	if cbm.Interval == 0 {
		cbm.Interval = 2 * time.Second
	}

	if cbm.Instance == nil {
		return errors.New("radium instance is not set")
	}

	return nil
}
