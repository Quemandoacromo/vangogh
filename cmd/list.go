package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/gog_types"
	"github.com/arelate/vangogh_types"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/kvas"
	"strings"
)

type productTitle struct {
	Title string `json:"title"`
}

func List(ids []string, title string, productType, media string) error {
	pt := vangogh_types.ParseProductType(productType)
	mt := gog_types.Parse(media)

	dstUrl, err := vangogh_urls.DestinationUrl(pt, mt)
	if err != nil {
		return err
	}

	kv, err := kvas.NewJsonLocal(dstUrl)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		ids = kv.All()
	}

	for _, id := range ids {
		rc, err := kv.Get(id)
		if err != nil {
			return err
		}
		var tt productTitle
		err = json.NewDecoder(rc).Decode(&tt)
		if err != nil {
			return err
		}

		if err := rc.Close(); err != nil {
			return err
		}

		if title != "" && !strings.Contains(
			strings.ToLower(tt.Title),
			strings.ToLower(title)) {
			continue
		}

		fmt.Println(id, tt.Title)
	}

	return nil
}
