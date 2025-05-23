package compton_data

import "github.com/arelate/southern_light/vangogh_integration"

var BinaryTitles = map[string]string{
	vangogh_integration.TrueValue:  "Yes",
	vangogh_integration.FalseValue: "No",
}

var TypesTitles = map[string]string{
	vangogh_integration.Licences.String():                     "Licences",
	vangogh_integration.UserWishlist.String():                 "User Wishlist",
	vangogh_integration.AccountPage.String():                  "Account Page",
	vangogh_integration.ApiProducts.String():                  "API Products",
	vangogh_integration.CatalogPage.String():                  "Catalog Page",
	vangogh_integration.Details.String():                      "Details",
	vangogh_integration.HltbData.String():                     "HowLongToBeat Data",
	vangogh_integration.HltbRootPage.String():                 "HowLongToBeat Root Page",
	vangogh_integration.OrderPage.String():                    "Order Page",
	vangogh_integration.PcgwSteamPageId.String():              "PCGamingWiki Steam PageId",
	vangogh_integration.PcgwGogPageId.String():                "PCGamingWiki GOG PageId",
	vangogh_integration.PcgwRaw.String():                      "PCGamingWiki Raw",
	vangogh_integration.WikipediaRaw.String():                 "Wikipedia Raw",
	vangogh_integration.SteamAppDetails.String():              "Steam App Details",
	vangogh_integration.SteamAppNews.String():                 "Steam App News",
	vangogh_integration.SteamAppReviews.String():              "Steam Reviews",
	vangogh_integration.ProtonDbSummary.String():              "ProtonDB Summary",
	vangogh_integration.GamesDbGogProducts.String():           "GamesDB GOG Products",
	vangogh_integration.SteamDeckCompatibilityReport.String(): "Steam Deck Compat Report",
	vangogh_integration.OpenCriticApiGame.String():            "OpenCritic API Game",
}

var OperatingSystemTitles = map[string]string{
	vangogh_integration.MacOS.String():   "macOS",
	vangogh_integration.Linux.String():   "Linux",
	vangogh_integration.Windows.String(): "Windows",
}
