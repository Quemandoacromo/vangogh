package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_properties"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/vangogh/cmd/iterate"
)

func Size(
	idSet gost.StrSet,
	mt gog_media.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	missing bool) error {

	dlList := vangogh_downloads.DownloadsList{}

	exl, err := vangogh_extracts.NewList(
		vangogh_properties.NativeLanguageNameProperty,
		vangogh_properties.SlugProperty)
	if err != nil {
		return err
	}

	if missing {
		missingIds, err := idMissingLocalDownloads(mt, exl, operatingSystems, downloadTypes, langCodes)
		if err != nil {
			return err
		}
		idSet.AddSet(missingIds)
	}

	if idSet.Len() == 0 {
		fmt.Println("no ids to estimate size")
		return nil
	}

	if err := iterate.DownloadsList(
		idSet,
		mt,
		exl,
		operatingSystems,
		downloadTypes,
		langCodes,
		func(
			_ string,
			list vangogh_downloads.DownloadsList,
			_ *vangogh_extracts.ExtractsList,
			_ bool) error {
			dlList = append(dlList, list...)
			return nil
		},
		0,
		false); err != nil {
		return err
	}

	fmt.Printf("estimated total download size: %.2fGB\n", dlList.TotalGBsEstimate())

	return nil
}
