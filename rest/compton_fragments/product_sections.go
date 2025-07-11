package compton_fragments

import (
	"github.com/arelate/southern_light/vangogh_integration"
	"github.com/arelate/vangogh/rest/compton_data"
	"github.com/boggydigital/redux"
)

func ProductSections(id string, rdx redux.Readable) []string {

	hasSections := make([]string, 0)

	hasSections = append(hasSections, compton_data.InfoSection)

	if sr, ok := rdx.GetLastVal(vangogh_integration.SummaryRatingProperty, id); ok && sr != "" {
		hasSections = append(hasSections, compton_data.ReceptionSection)
	}

	offeringsCount := 0
	for _, rpp := range compton_data.OfferingsProperties {
		if rps, ok := rdx.GetAllValues(rpp, id); ok {
			offeringsCount += len(rps)
		}
	}

	if offeringsCount > 0 {
		hasSections = append(hasSections, compton_data.OfferingsSection)
	}

	if rdx.HasKey(vangogh_integration.ScreenshotsProperty, id) ||
		rdx.HasKey(vangogh_integration.VideoIdProperty, id) {
		hasSections = append(hasSections, compton_data.MediaSection)
	}

	if rdx.HasValue(vangogh_integration.TypesProperty, id, vangogh_integration.SteamAppNews.String()) ||
		rdx.HasKey(vangogh_integration.ChangelogProperty, id) {
		hasSections = append(hasSections, compton_data.NewsSection)
	}

	if rdx.HasKey(vangogh_integration.SteamDeckAppCompatibilityCategoryProperty, id) ||
		rdx.HasKey(vangogh_integration.ProtonDBTierProperty, id) {
		hasSections = append(hasSections, compton_data.CompatibilitySection)
	}

	if val, ok := rdx.GetLastVal(vangogh_integration.OwnedProperty, id); ok && val == vangogh_integration.TrueValue {
		if productType, sure := rdx.GetLastVal(vangogh_integration.ProductTypeProperty, id); sure && productType == vangogh_integration.GameProductType {
			hasSections = append(hasSections, compton_data.InstallersSection)
		}
	}

	return hasSections
}
