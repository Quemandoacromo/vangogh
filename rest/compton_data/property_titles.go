package compton_data

import (
	"github.com/arelate/southern_light/vangogh_integration"
)

var PropertyTitles = map[string]string{
	vangogh_integration.TitleProperty:                             "Title",
	vangogh_integration.DescriptionOverviewProperty:               "Description",
	vangogh_integration.TagIdProperty:                             "Account Tags",
	vangogh_integration.LocalTagsProperty:                         "Local Tags",
	vangogh_integration.SteamTagsProperty:                         "Steam Tags",
	vangogh_integration.OperatingSystemsProperty:                  "OS",
	vangogh_integration.DevelopersProperty:                        "Developers",
	vangogh_integration.PublishersProperty:                        "Publishers",
	vangogh_integration.EnginesProperty:                           "Engine",
	vangogh_integration.EnginesBuildsProperty:                     "Engine Build",
	vangogh_integration.SeriesProperty:                            "Series",
	vangogh_integration.GenresProperty:                            "Genres",
	vangogh_integration.ThemesProperty:                            "Themes",
	vangogh_integration.StoreTagsProperty:                         "Store Tags",
	vangogh_integration.FeaturesProperty:                          "Features",
	vangogh_integration.GameModesProperty:                         "Game Modes",
	vangogh_integration.LanguageCodeProperty:                      "Language",
	vangogh_integration.IncludesGamesProperty:                     "Includes",
	vangogh_integration.IsIncludedByGamesProperty:                 "Included By",
	vangogh_integration.RequiresGamesProperty:                     "Requires",
	vangogh_integration.IsRequiredByGamesProperty:                 "Required By",
	vangogh_integration.ProductTypeProperty:                       "Product Type",
	vangogh_integration.WishlistedProperty:                        "Wishlisted",
	vangogh_integration.OwnedProperty:                             "Owned",
	vangogh_integration.IsFreeProperty:                            "Free",
	vangogh_integration.IsDiscountedProperty:                      "On Sale",
	vangogh_integration.PreOrderProperty:                          "Pre-order",
	vangogh_integration.ComingSoonProperty:                        "Coming Soon",
	vangogh_integration.InDevelopmentProperty:                     "In Development",
	vangogh_integration.TypesProperty:                             "Data Type",
	vangogh_integration.SteamReviewScoreDescProperty:              "Steam Reviews",
	vangogh_integration.SteamDeckAppCompatibilityCategoryProperty: "Steam Deck",
	vangogh_integration.ProtonDBTierProperty:                      "ProtonDB Tier",
	vangogh_integration.ProtonDBConfidenceProperty:                "ProtonDB Confidence",
	vangogh_integration.SortProperty:                              "Sort",
	vangogh_integration.DescendingProperty:                        "Descending",
	vangogh_integration.GlobalReleaseDateProperty:                 "Global Release",
	vangogh_integration.GOGReleaseDateProperty:                    "GOG.com Release",
	vangogh_integration.GOGOrderDateProperty:                      "GOG.com Order",
	vangogh_integration.ProductValidationResultProperty:           "Validation Result",
	vangogh_integration.RatingProperty:                            "Rating",
	vangogh_integration.AggregatedRatingProperty:                  "Aggregated Rating",
	vangogh_integration.PriceProperty:                             "Price",
	vangogh_integration.BasePriceProperty:                         "Base Price",
	vangogh_integration.DiscountPercentageProperty:                "Discount",

	vangogh_integration.HLTBHoursToCompleteMainProperty: "HLTB Main Story",
	vangogh_integration.HLTBHoursToCompletePlusProperty: "HLTB Story + Extras",
	vangogh_integration.HLTBHoursToComplete100Property:  "HLTB Completionist",
	vangogh_integration.HLTBGenresProperty:              "HLTB Genres",
	vangogh_integration.HLTBPlatformsProperty:           "HLTB Platforms",
	vangogh_integration.HLTBReviewScoreProperty:         "HLTB Review Score",

	GauginGOGLinksProperty:   "GOG.com Links",
	GauginOtherLinksProperty: "Other Links",
	GauginSteamLinksProperty: "Steam Links",

	vangogh_integration.ForumUrlProperty:   "Forum",
	vangogh_integration.StoreUrlProperty:   "Store",
	vangogh_integration.SupportUrlProperty: "Support",

	GauginSteamCommunityUrlProperty: "Community",

	GauginGOGDBUrlProperty:        "GOGDB",
	GauginIGDBUrlProperty:         "IGDB",
	GauginHLTBUrlProperty:         "HLTB",
	GauginMobyGamesUrlProperty:    "MobyGames",
	GauginPCGamingWikiUrlProperty: "PCGamingWiki",
	GauginProtonDBUrlProperty:     "ProtonDB",
	GauginStrategyWikiUrlProperty: "StrategyWiki",
	GauginWikipediaUrlProperty:    "Wikipedia",
	GauginWineHQUrlProperty:       "WineHQ",
	GauginVNDBUrlProperty:         "VNDB",
	GauginIGNWikiUrlProperty:      "IGN Wiki",
}
