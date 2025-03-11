package gog_data

import (
	"github.com/arelate/southern_light/vangogh_integration"
	"github.com/arelate/vangogh/cli/fetch"
	"github.com/arelate/vangogh/cli/reqs"
	"github.com/arelate/vangogh/cli/shared_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"maps"
	"net/http"
)

func GetDetails(hc *http.Client, uat string, since int64) error {

	gda := nod.NewProgress("getting new or updated %s...", vangogh_integration.Details)
	defer gda.Done()

	detailsDir, err := vangogh_integration.AbsProductTypeDir(vangogh_integration.Details)
	if err != nil {
		return err
	}

	kvDetails, err := kevlar.New(detailsDir, kevlar.JsonExt)
	if err != nil {
		return err
	}

	newUpdatedDetails, err := shared_data.GetDetailsUpdates(since)

	gda.TotalInt(len(newUpdatedDetails))

	if err = fetch.Items(maps.Keys(newUpdatedDetails), reqs.Details(hc, uat), kvDetails, gda); err != nil {
		return err
	}

	return ReduceDetails(kvDetails, since)
}

func ReduceDetails(kvDetails kevlar.KeyValues, since int64) error {
	rda := nod.Begin(" reducing %s...", vangogh_integration.Details)
	defer rda.Done()

	reduxDir, err := pathways.GetAbsRelDir(vangogh_integration.Redux)
	if err != nil {
		return err
	}

	rdx, err := redux.NewWriter(reduxDir, vangogh_integration.GOGDetailsProperties()...)
	if err != nil {
		return err
	}

	detailReductions := shared_data.InitReductions(vangogh_integration.GOGDetailsProperties()...)

	updatedDetails := kvDetails.Since(since, kevlar.Create, kevlar.Update)

	for id := range updatedDetails {
		if err = reduceDetailsProduct(id, kvDetails, detailReductions); err != nil {
			return err
		}
	}

	return shared_data.WriteReductions(rdx, detailReductions)
}

func reduceDetailsProduct(id string, kvDetails kevlar.KeyValues, piv shared_data.PropertyIdValues) error {

	det, err := vangogh_integration.UnmarshalDetails(id, kvDetails)
	if err != nil {
		return err
	}

	if det == nil {
		return nil
	}

	for property := range piv {

		var values []string

		switch property {
		case vangogh_integration.TitleProperty:
			values = []string{det.GetTitle()}
		case vangogh_integration.FeaturesProperty:
			values = det.GetFeatures()
		case vangogh_integration.TagIdProperty:
			values = det.GetTagIds()
		case vangogh_integration.GOGReleaseDateProperty:
			values = []string{det.GetGOGRelease()}
		case vangogh_integration.ForumUrlProperty:
			values = []string{det.GetForumUrl()}
		case vangogh_integration.ChangelogProperty:
			values = []string{det.GetChangelog()}
		}

		piv[property][id] = values
	}

	return nil
}
