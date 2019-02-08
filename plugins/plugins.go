package plugins

import (
	"sort"

	"gitlab.com/Cacophony/Processor/plugins/ping"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin interface {
	Name() string

	// TODO: add context for deadline
	Start() error

	// TODO: add context for deadline
	Stop() error

	Priority() int

	Passthrough() bool

	Action(event events.Event) bool
}

// nolint: gochecknoglobals
var (
	PluginList = []Plugin{
		&ping.Ping{},
	}
)

// nolint: gochecknoinits
func init() {
	// sort plugins by priority
	sort.Sort(ByPriority(PluginList))
}

func StartPlugins(logger *zap.Logger) {
	var err error
	for _, plugin := range PluginList {
		err = plugin.Start()
		if err != nil {
			logger.Error("failed to start plugin",
				zap.Error(err),
			)
		}
	}
}

func StopPlugins(logger *zap.Logger) {
	var err error
	for _, plugin := range PluginList {
		err = plugin.Stop()
		if err != nil {
			logger.Error("failed to stop plugin",
				zap.Error(err),
			)
		}
	}
}

type ByPriority []Plugin

func (p ByPriority) Len() int           { return len(p) }
func (p ByPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByPriority) Less(i, j int) bool { return p[i].Priority() > p[j].Priority() }
