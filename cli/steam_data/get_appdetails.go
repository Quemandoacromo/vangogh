package steam_data

import (
	"encoding/json"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/arelate/southern_light/vangogh_integration"
	"github.com/arelate/vangogh/cli/fetch"
	"github.com/arelate/vangogh/cli/reqs"
	"github.com/arelate/vangogh/cli/shared_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"slices"
	"strconv"
)

func GetAppDetails(since int64, force bool) error {

	gada := nod.NewProgress("getting %s...", vangogh_integration.SteamAppDetails)
	defer gada.Done()

	if force {
		since = -1
	}

	reduxDir, err := pathways.GetAbsRelDir(vangogh_integration.Redux)
	if err != nil {
		return err
	}

	rdx, err := redux.NewReader(reduxDir, vangogh_integration.SteamAppIdProperty)
	if err != nil {
		return err
	}

	steamAppDetailsDir, err := vangogh_integration.AbsProductTypeDir(vangogh_integration.SteamAppDetails)
	if err != nil {
		return err
	}

	kvSteamAppDetails, err := kevlar.New(steamAppDetailsDir, kevlar.JsonExt)
	if err != nil {
		return err
	}

	var newSteamAppIds []string

	for gogId := range rdx.Keys(vangogh_integration.SteamAppIdProperty) {
		if steamAppIds, ok := rdx.GetAllValues(vangogh_integration.SteamAppIdProperty, gogId); ok {
			for _, steamAppId := range steamAppIds {
				if kvSteamAppDetails.Has(steamAppId) && !force {
					continue
				}
				newSteamAppIds = append(newSteamAppIds, steamAppId)
			}
		}
	}

	gada.TotalInt(len(newSteamAppIds))

	if err = fetch.Items(slices.Values(newSteamAppIds), reqs.SteamAppDetails(), kvSteamAppDetails, gada); err != nil {
		return err
	}

	return reduceSteamAppDetails(kvSteamAppDetails, since)
}

func reduceSteamAppDetails(kvSteamAppDetails kevlar.KeyValues, since int64) error {

	rspada := nod.Begin(" reducing %s...", vangogh_integration.SteamAppDetails)
	defer rspada.Done()

	reduxDir, err := pathways.GetAbsRelDir(vangogh_integration.Redux)
	if err != nil {
		return err
	}

	rdx, err := redux.NewWriter(reduxDir, vangogh_integration.SteamAppDetailsProperties()...)
	if err != nil {
		return err
	}

	steamAppDetailsReductions := shared_data.InitReductions(vangogh_integration.SteamAppDetailsProperties()...)

	updatedSteamAppDetails := kvSteamAppDetails.Since(since, kevlar.Create, kevlar.Update)

	for steamAppId := range updatedSteamAppDetails {

		for gogId := range rdx.Match(map[string][]string{vangogh_integration.SteamAppIdProperty: {steamAppId}}, redux.FullMatch) {
			if err = reduceSteamAppDetailsProduct(gogId, steamAppId, kvSteamAppDetails, steamAppDetailsReductions); err != nil {
				return err
			}
		}

	}

	return shared_data.WriteReductions(rdx, steamAppDetailsReductions)
}

func reduceSteamAppDetailsProduct(gogId, steamAppId string, kvSteamAppDetails kevlar.KeyValues, piv shared_data.PropertyIdValues) error {

	rcSteamAppDetailsResponse, err := kvSteamAppDetails.Get(steamAppId)
	if err != nil {
		return err
	}
	defer rcSteamAppDetailsResponse.Close()

	var sadr steam_integration.AppDetailsResponse
	if err = json.NewDecoder(rcSteamAppDetailsResponse).Decode(&sadr); err != nil {
		return err
	}

	ad := sadr.GetAppDetails()

	for property := range piv {

		var values []string

		switch property {
		case vangogh_integration.RequiredAgeProperty:
			values = []string{strconv.FormatInt(int64(ad.GetRequiredAge()), 10)}
		case vangogh_integration.ControllerSupportProperty:
			values = []string{ad.GetControllerSupport()}
		case vangogh_integration.ShortDescriptionProperty:
			values = []string{ad.GetShortDescription()}
		case vangogh_integration.WebsiteProperty:
			values = []string{ad.GetWebsite()}
		case vangogh_integration.MetacriticScoreProperty:
			values = []string{strconv.FormatInt(int64(ad.GetMetacriticScore()), 10)}
		case vangogh_integration.MetacriticUrlProperty:
			values = []string{ad.GetMetacriticUrl()}
		case vangogh_integration.SteamCategoriesProperty:
			values = ad.GetCategories()
		case vangogh_integration.SteamGenresProperty:
			values = ad.GetGenres()
		case vangogh_integration.SteamSupportUrlProperty:
			values = []string{ad.GetSupportUrl()}
		case vangogh_integration.SteamSupportEmailProperty:
			values = []string{ad.GetSupportEmail()}
		}

		piv[property][gogId] = values

	}

	return nil
}
