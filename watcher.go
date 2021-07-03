package main

import (
	"time"

	"github.com/radovskyb/watcher"
)

func WatchFiles() error {
	w := watcher.New()
	// experimenting with 2
	// TODO: add it to the config so people can change if needed
	w.SetMaxEvents(2)

	go handleFileChanges(w)

	if err := w.AddRecursive(ConfigRef.ContentPath); err != nil {
		return err
	}

	if err := w.AddRecursive(ConfigRef.PublicFolder); err != nil {
		return err
	}

	if err := w.Start(time.Millisecond * 500); err != nil {
		return err
	}

	return nil
}

func handleFileChanges(w *watcher.Watcher) error {
	for {
		select {
		case <-w.Event:
			Statico()
		case err := <-w.Error:
			return err
		case <-w.Closed:
			return nil
		}
	}
}
