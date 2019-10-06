package help

import (
	"encoding/json"
	"net/http"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/localization"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) endpointCommands() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseHelpList := make([]*common.PluginHelp, 0, len(p.pluginHelpList))

	NextPluginHelp:
		for _, pluginHelp := range p.pluginHelpList {
			if pluginHelp.Hide {
				continue
			}

			for _, permissionRequired := range pluginHelp.PermissionsRequired {
				if permissionRequired == permissions.BotAdmin {
					continue NextPluginHelp
				}
			}

			pluginHelp.Description = localization.Translate(
				p.localizations,
				pluginHelp.Description,
			)

			for i := range pluginHelp.Commands {
				pluginHelp.Commands[i].Name = localization.Translate(
					p.localizations,
					pluginHelp.Commands[i].Name,
				)
				pluginHelp.Commands[i].Description = localization.Translate(
					p.localizations,
					pluginHelp.Commands[i].Description,
				)
			}

			for i := range pluginHelp.Reactions {
				pluginHelp.Reactions[i].Description = localization.Translate(
					p.localizations,
					pluginHelp.Reactions[i].Description,
				)
			}

			responseHelpList = append(responseHelpList, pluginHelp)
		}

		resp, err := json.Marshal(responseHelpList)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}

}
