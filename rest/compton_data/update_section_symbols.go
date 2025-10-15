package compton_data

import (
	"github.com/arelate/southern_light/vangogh_integration"
	"github.com/boggydigital/compton"
)

var UpdateSectionSymbols = map[string]compton.Symbol{
	vangogh_integration.UpdatesInstallers:    compton.CompactDisk,
	vangogh_integration.UpdatesReleasedToday: compton.RisingSun,
	vangogh_integration.UpdatesNewProducts:   compton.ShoppingLabel,
	vangogh_integration.UpdatesSteamNews:     compton.NewsBroadcast,
}
