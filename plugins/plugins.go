package plugins

import (
	"path/filepath"
	"sort"

	"github.com/go-chi/chi"
	"gitlab.com/Cacophony/Processor/plugins/eventlog"
	"gitlab.com/Cacophony/Processor/plugins/greeter"
	"gitlab.com/Cacophony/Processor/plugins/rpg"
	"gitlab.com/Cacophony/Processor/plugins/stats"
	"gitlab.com/Cacophony/Processor/plugins/tiktok"
	"gitlab.com/Cacophony/Processor/plugins/tools"
	"gitlab.com/Cacophony/Processor/plugins/uploads"
	"gitlab.com/Cacophony/Processor/plugins/weather"
	"gitlab.com/Cacophony/Processor/plugins/weverse"
	"gitlab.com/Cacophony/go-kit/localization"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/admin"
	"gitlab.com/Cacophony/Processor/plugins/automod"
	"gitlab.com/Cacophony/Processor/plugins/color"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/Processor/plugins/customcommands"
	"gitlab.com/Cacophony/Processor/plugins/dev"
	"gitlab.com/Cacophony/Processor/plugins/gall"
	"gitlab.com/Cacophony/Processor/plugins/help"
	"gitlab.com/Cacophony/Processor/plugins/instagram"
	"gitlab.com/Cacophony/Processor/plugins/lastfm"
	"gitlab.com/Cacophony/Processor/plugins/patreon"
	"gitlab.com/Cacophony/Processor/plugins/ping"
	"gitlab.com/Cacophony/Processor/plugins/prefix"
	"gitlab.com/Cacophony/Processor/plugins/quickactions"
	"gitlab.com/Cacophony/Processor/plugins/roles"
	"gitlab.com/Cacophony/Processor/plugins/rss"
	"gitlab.com/Cacophony/Processor/plugins/serverlist"
	"gitlab.com/Cacophony/Processor/plugins/stocks"
	"gitlab.com/Cacophony/Processor/plugins/trello"
	"gitlab.com/Cacophony/Processor/plugins/whitelist"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/featureflag"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin interface {
	Names() []string

	// TODO: add context for deadline
	Start(common.StartParameters) error

	// TODO: add context for deadline
	Stop(common.StopParameters) error

	Priority() int

	Passthrough() bool

	Action(event *events.Event) bool

	Help() *common.PluginHelp
}

var (
	PluginList = []Plugin{
		&weather.Plugin{},
		&trello.Plugin{},
		&admin.Plugin{},
		&patreon.Plugin{},
		&ping.Plugin{},
		&color.Plugin{},
		&dev.Plugin{},
		&lastfm.Plugin{},
		&automod.Plugin{},
		&whitelist.Plugin{},
		&gall.Plugin{},
		&rss.Plugin{},
		&roles.Plugin{},
		&eventlog.Plugin{},
		&instagram.Plugin{},
		&serverlist.Plugin{},
		&prefix.Plugin{},
		&help.Plugin{},
		&stocks.Plugin{},
		&quickactions.Plugin{},
		&customcommands.Plugin{},
		&stats.Plugin{},
		&weverse.Plugin{},
		&tools.Plugin{},
		&uploads.Plugin{},
		&greeter.Plugin{},
		&rpg.Plugin{},
		&tiktok.Plugin{},
	}

	LocalizationsList []interfaces.Localization
)

func init() {
	// sort plugins by priority
	sort.Sort(ByPriority(PluginList))
}

func StartPlugins(
	logger *zap.Logger,
	db *gorm.DB,
	redis *redis.Client,
	tokens map[string]string,
	state *state.State,
	featureFlagger *featureflag.FeatureFlagger,
	publisher *events.Publisher,
	httpMux *chi.Mux,
) {

	// load localizations
	loadLocalizations(logger)

	// get help list from all pluguins for help module
	pluginHelpList := make([]*common.PluginHelp, len(PluginList))
	for i, plugin := range PluginList {
		pluginHelpList[i] = plugin.Help()
	}

	var err error
	for _, plugin := range PluginList {
		if len(plugin.Names()) <= 0 {
			logger.Error("plugin has no names, failed to launch",
				zap.Error(err),
			)
			continue
		}

		err = plugin.Start(common.StartParameters{
			Logger:         logger,
			DB:             db,
			Redis:          redis,
			Tokens:         tokens,
			State:          state,
			FeatureFlagger: featureFlagger,
			PluginHelpList: pluginHelpList,
			Localizations:  LocalizationsList,
			Publisher:      publisher,
			HTTPMux:        httpMux,
		})
		if err != nil {
			logger.Error("failed to start plugin",
				zap.String("name", plugin.Names()[0]),
				zap.Error(err),
			)
		}
		// TODO: do not send events to plugins that failed to start
	}
}

func StopPlugins(
	logger *zap.Logger,
	db *gorm.DB,
	redis *redis.Client,
	tokens map[string]string,
	state *state.State,
	featureFlagger *featureflag.FeatureFlagger,
	publisher *events.Publisher,
) {
	var err error
	for _, plugin := range PluginList {
		err = plugin.Stop(common.StopParameters{
			Logger:         logger,
			DB:             db,
			Redis:          redis,
			Tokens:         tokens,
			State:          state,
			FeatureFlagger: featureFlagger,
			Publisher:      publisher,
		})
		if err != nil {
			logger.Error("failed to stop plugin",
				zap.String("name", plugin.Names()[0]),
				zap.Error(err),
			)
		}
	}
}

func loadLocalizations(logger *zap.Logger) {
	files, err := filepath.Glob("assets/translations/*.en.toml")
	if err != nil {
		logger.Error("could not load any localization files", zap.Error(err))
		return
	}

	for _, path := range files {
		local, err := localization.NewFileSource(path, "en")
		if err != nil {
			logger.Error("failed to load localization file", zap.String("file", path), zap.Error(err))
		}
		LocalizationsList = append(LocalizationsList, local)
	}
}

type ByPriority []Plugin

func (p ByPriority) Len() int           { return len(p) }
func (p ByPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByPriority) Less(i, j int) bool { return p[i].Priority() > p[j].Priority() }
