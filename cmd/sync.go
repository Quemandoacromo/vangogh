package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_images"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/vangogh/internal"
	"time"
)

func Sync(mt gog_media.Media, sinceHoursAgo int, data, images, screenshots, videos, verbose bool) error {

	var syncStart int64
	if sinceHoursAgo > 0 {
		syncStart = time.Now().Unix() - int64(sinceHoursAgo*60*60)
	} else {
		syncStart = time.Now().Unix()
	}

	if data {
		//get array and paged data
		paData := vangogh_products.Array()
		paData = append(paData, vangogh_products.Paged()...)
		for _, pt := range paData {
			if err := GetData(nil, "", nil, pt, mt, syncStart, false, false, verbose); err != nil {
				return err
			}
			fmt.Println()
		}

		//get main - detail data
		for _, pt := range vangogh_products.Detail() {
			denyIds := internal.ReadLines(vangogh_urls.Denylist(pt))
			if err := GetData(nil, "", denyIds, pt, mt, syncStart, true, true, verbose); err != nil {
				return err
			}
			fmt.Println()
		}

		//extract data
		if err := Extract(syncStart, mt, vangogh_properties.Extracted()); err != nil {
			return err
		}
		fmt.Println()
	}

	localImageIds, err := vangogh_urls.LocalImageIds()
	if err != nil {
		return err
	}
	// get images
	if images {
		for _, it := range vangogh_images.All() {
			if !screenshots && it == vangogh_images.Screenshots {
				continue
			}
			if err := GetImages(nil, "", it, localImageIds, true); err != nil {
				return err
			}
		}
		fmt.Println()
	}

	// get videos
	if videos {
		if err := GetVideos(nil, "", true); err != nil {
			return err
		}
		fmt.Println()
	}

	// TODO: get files

	// print new or updated
	return Summary(syncStart, mt)
}
