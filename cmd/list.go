package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"time"
)

//List prints products of a certain type and media.
//Can be filtered to products that were created or modified since a certain time.
//Provided properties will be printed for each product (if supported) in addition to default ID, Title.
func List(
	ids map[string]bool,
	modifiedAfter int64,
	pt vangogh_products.ProductType,
	mt gog_media.Media,
	properties map[string]bool) error {

	if !vangogh_products.Valid(pt) {
		return fmt.Errorf("can't list invalid product type %s", pt)
	}
	if !gog_media.Valid(mt) {
		return fmt.Errorf("can't list invalid media %s", mt)
	}

	if properties == nil {
		properties = make(map[string]bool, 0)
	}

	//if no properties have been provided - print ID, Title
	if len(properties) == 0 {
		properties[vangogh_properties.IdProperty] = true
		properties[vangogh_properties.TitleProperty] = true
	}

	//if Title property has not been provided - add it first.
	//we'll always print the title
	if !properties[vangogh_properties.TitleProperty] {
		properties[vangogh_properties.TitleProperty] = true
	}

	//rules for collecting IDs to print:
	//1. start with user provided IDs
	//2. if createdAfter has been provided - add products created since that time
	//3. if modifiedAfter has been provided - add products modified (not by creation!) since that time
	//4. if no IDs have been collected and the request have not provided createdAfter or modifiedAfter:
	// add all product IDs

	if ids == nil {
		ids = make(map[string]bool, 0)
	}

	vr, err := vangogh_values.NewReader(pt, mt)
	if err != nil {
		return err
	}

	if modifiedAfter > 0 {
		for _, mId := range vr.ModifiedAfter(modifiedAfter, false) {
			ids[mId] = true
		}
		if len(ids) == 0 {
			fmt.Printf("no new or updated %s (%s) since %v\n", pt, mt, time.Unix(modifiedAfter, 0).Format(time.Kitchen))
		}
	}

	if len(ids) == 0 &&
		modifiedAfter == 0 {
		for _, id := range vr.All() {
			ids[id] = true
		}
	}

	//load properties extract that will be used for printing
	exl, err := vangogh_extracts.NewListFromMap(properties)

	//use common printInfo func to display product information by ID
	for id, ok := range ids {
		if !ok {
			continue
		}
		if err := printInfo(
			id,
			nil,
			vangogh_properties.Supported(pt, properties),
			exl); err != nil {
			return err
		}
	}

	return nil
}
