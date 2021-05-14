package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_urls"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/kvas"
)

func itemizeAll(
	ids []string,
	missing, updated bool,
	modifiedAfter int64,
	pt vangogh_products.ProductType,
	mt gog_media.Media) ([]string, error) {

	for _, mainPt := range vangogh_products.MainTypes(pt) {
		if missing {
			missingIds, err := itemizeMissing(pt, mainPt, mt, modifiedAfter)
			if err != nil {
				return ids, err
			}
			if len(missingIds) == 0 {
				fmt.Printf("no missing %s data for %s (%s)\n", pt, mainPt, mt)
			}
			ids = append(ids, missingIds...)
		}
		if updated {
			updatedIds, err := itemizeUpdated(modifiedAfter, mainPt, mt)
			if err != nil {
				return ids, err
			}
			if len(updatedIds) == 0 {
				fmt.Printf("no updated %s data for %s (%s)\n", pt, mainPt, mt)
			}
			ids = append(ids, updatedIds...)
		}
	}

	return ids, nil
}

func itemizeMissing(
	detailPt, mainPt vangogh_products.ProductType,
	mt gog_media.Media,
	modifiedAfter int64) ([]string, error) {

	//api-products-v2 provides
	//includes-games, is-included-by-games,
	//requires-games, is-required-by-games
	if mainPt == vangogh_products.ApiProductsV2 &&
		detailPt == vangogh_products.ApiProductsV2 {
		return itemizeAPV2LinkedGames(modifiedAfter)
	}

	//licences give a signal when DLC has been purchased, this would add
	//required (base) game details to the updates
	if mainPt == vangogh_products.LicenceProducts &&
		detailPt == vangogh_products.Details {
		return itemizeRequiredGames(modifiedAfter, mt)
	}

	missingIds := make([]string, 0)

	mainDestUrl, err := vangogh_urls.LocalProductsDir(mainPt, mt)
	if err != nil {
		return missingIds, err
	}

	detailDestUrl, err := vangogh_urls.LocalProductsDir(detailPt, mt)
	if err != nil {
		return missingIds, err
	}

	kvMain, err := kvas.NewJsonLocal(mainDestUrl)
	if err != nil {
		return missingIds, err
	}

	kvDetail, err := kvas.NewJsonLocal(detailDestUrl)
	if err != nil {
		return missingIds, err
	}
	for _, id := range kvMain.All() {
		if !kvDetail.Contains(id) {
			missingIds = append(missingIds, id)
		}
	}

	return missingIds, nil
}

func itemizeUpdated(
	since int64,
	pt vangogh_products.ProductType,
	mt gog_media.Media) ([]string, error) {

	updatedIds := make([]string, 0)

	//licence products can only update through creation and we've already handled
	//newly created in itemizeMissing func
	if pt == vangogh_products.LicenceProducts {
		return updatedIds, nil
	}

	mainDestUrl, err := vangogh_urls.LocalProductsDir(pt, mt)
	if err != nil {
		return updatedIds, err
	}

	kvMain, err := kvas.NewJsonLocal(mainDestUrl)
	if err != nil {
		return updatedIds, err
	}

	updatedIds = kvMain.ModifiedAfter(since, false)

	return updatedIds, nil
}

func itemizeAPV2LinkedGames(modifiedAfter int64) ([]string, error) {

	missing := make(map[string]bool, 0)

	//currently api-products-v2 support only gog_media.Game, and since this method is exclusively
	//using api-products-v2 we're fine specifying media directly and not taking as a parameter
	vrApv2, err := vangogh_values.NewReader(vangogh_products.ApiProductsV2, gog_media.Game)

	if err != nil {
		return []string{}, err
	}

	for _, id := range vrApv2.ModifiedAfter(modifiedAfter, false) {

		// have to use product reader and not extracts here, since extracts wouldn't be ready
		// while we're still getting data. Attempting to minimize the impact by only querying
		// new or updated api-product-v2 items since start to the sync
		apv2, err := vrApv2.ApiProductV2(id)

		if err != nil {
			return []string{}, err
		}

		linkedGames := apv2.GetIncludesGames()
		linkedGames = append(linkedGames, apv2.GetIsIncludedInGames()...)
		linkedGames = append(linkedGames, apv2.GetRequiresGames()...)
		linkedGames = append(linkedGames, apv2.GetIsRequiredByGames()...)

		for _, lid := range linkedGames {
			if !vrApv2.Contains(lid) {
				missing[lid] = true
			}
		}
	}

	missingIds := make([]string, 0, len(missing))
	for id, _ := range missing {
		missingIds = append(missingIds, id)
	}

	return missingIds, nil
}

//itemizeRequiredGames enumerates all base products for a newly acquired DLCs
func itemizeRequiredGames(createdAfter int64, mt gog_media.Media) ([]string, error) {
	requiredGamesForNewLicences := make([]string, 0)

	vrLicences, err := vangogh_values.NewReader(vangogh_products.LicenceProducts, mt)
	if err != nil {
		return requiredGamesForNewLicences, err
	}

	vrApv2, err := vangogh_values.NewReader(vangogh_products.ApiProductsV2, gog_media.Game)
	if err != nil {
		return requiredGamesForNewLicences, err
	}

	for _, id := range vrLicences.CreatedAfter(createdAfter) {
		//like in itemizeMissingIncludesGames, we can't use extracts here,
		//because we're in process of getting data and would rather query api-products-v2 directly.
		//the performance impact is expected to be minimal since we're only loading newly acquired licences
		apv2, err := vrApv2.ApiProductV2(id)
		if err != nil {
			return requiredGamesForNewLicences, err
		}

		for _, rg := range apv2.GetRequiresGames() {
			if !stringsContain(requiredGamesForNewLicences, rg) {
				requiredGamesForNewLicences = append(requiredGamesForNewLicences, rg)
			}
		}
	}

	return requiredGamesForNewLicences, nil
}
