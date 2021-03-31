package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_images"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/vangogh/internal"
	"log"
	"time"
)

func Sync(mt gog_media.Media, noData, images, screenshots, verbose bool) error {

	syncStart := time.Now().Unix()

	if !noData {
		// get paginated data
		for _, pt := range vangogh_products.AllPaged() {
			if err := GetData(nil, nil, pt, mt, syncStart, false, false, verbose); err != nil {
				return err
			}
		}

		// get main - detail data
		for _, pt := range vangogh_products.AllDetail() {
			denyIds := internal.ReadLines(vangogh_urls.Denylist(pt))
			if err := GetData(nil, denyIds, pt, mt, syncStart, true, true, verbose); err != nil {
				return err
			}
		}

		// extract data
		if err := Extract(mt, vangogh_properties.AllExtracted(), false); err != nil {
			return err
		}
	}

	localImageIds, err := vangogh_urls.LocalImageIds()
	if err != nil {
		return err
	}
	// get images
	for _, it := range vangogh_images.All() {
		if !images ||
			(!screenshots && it == vangogh_images.Screenshots) {
			continue
		}
		if err := GetImages(nil, it, localImageIds, true); err != nil {
			return err
		}
	}

	// TODO: get files

	// print created, modified
	return reportCreatedModifiedAfter(syncStart, mt)

}

func reportCreatedModifiedAfter(timestamp int64, mt gog_media.Media) error {
	for _, pt := range vangogh_products.AllLocal() {
		fmt.Printf("%s (%s) created this sync:\n", pt, mt)
		if err := List(nil, timestamp, 0, pt, mt); err != nil {
			return err
		}
		fmt.Printf("%s (%s) modified this sync:\n", pt, mt)
		if err := List(nil, 0, timestamp, pt, mt); err != nil {
			return err
		}
	}
	log.Println("sync took:", time.Since(time.Unix(timestamp, 0)).String())
	return nil
}
