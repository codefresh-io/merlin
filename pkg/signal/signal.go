package signal

import (
	"os"
	"os/signal"
	"syscall"
)

type (
	Handler interface {
		Process()
		Push(chan os.Signal)
	}

	handler struct {
		signals  []syscall.Signal
		channels []chan os.Signal
		logger   logger
	}

	logger interface {
		Debugf(format string, args ...interface{})
		Debug(args ...interface{})
	}
)

func NewSignalhandler(signals []syscall.Signal, logger logger) Handler {
	return &handler{
		signals:  signals,
		channels: nil,
		logger:   logger,
	}
}

func (h *handler) Push(c chan os.Signal) {
	h.logger.Debug("Adding another signal handler to slice")
	h.channels = append(h.channels, c)
}

func (h *handler) Process() {
	var channel = make(chan os.Signal)
	for _, s := range h.signals {
		signal.Notify(channel, s)
	}
	go func() {
		sig := <-channel
		h.logger.Debugf("Got signal: %s\n", sig)
		for _, ch := range h.channels {
			h.logger.Debug("Sending signal to channel")
			ch <- sig
			res := <-ch
			h.logger.Debugf("Return from channel %s\n", res)
		}
	}()

}

