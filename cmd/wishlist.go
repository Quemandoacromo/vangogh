package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/gog_urls"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/froth"
	"github.com/boggydigital/vangogh/internal"
	"log"
	"net/http"
	"net/url"
)

func Wishlist(mt gog_media.Media, addProductIds, removeProductIds []string) error {
	httpClient, err := internal.HttpClient()
	if err != nil {
		return err
	}

	vrStoreProducts, err := vangogh_values.NewReader(vangogh_products.StoreProducts, mt)
	if err != nil {
		return err
	}

	titleExtracts, err := froth.NewStash(
		vangogh_urls.ExtractsDir(),
		vangogh_properties.TitleProperty)
	if err != nil {
		return err
	}

	if err := wishlistAdd(addProductIds, httpClient, vrStoreProducts, mt, titleExtracts); err != nil {
		return err
	}

	if err := wishlistRemove(removeProductIds, httpClient, vrStoreProducts, mt, titleExtracts); err != nil {
		return err
	}

	return nil
}

func addError(id string) string {
	return fmt.Sprintf("can't add invalid id %s to the wishlist", id)
}

func removeError(id string) string {
	return fmt.Sprintf("can't remove invalid id %s from the wishlist", id)
}

func addProgress(title, id string) string {
	return fmt.Sprintf("add %s (%s) to the wishlist", title, id)
}

func removeProgress(title, id string) string {
	return fmt.Sprintf("remove %s (%s) from the wishlist", title, id)
}

func wishlistAdd(
	ids []string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_values.ValueReader,
	mt gog_media.Media,
	titleExtracts *froth.Stash) error {

	for _, id := range ids {
		if err := vrStoreProducts.CopyToType(id, vangogh_products.WishlistProducts, mt); err != nil {
			return err
		}
	}

	return remoteWishlistCommand(
		ids,
		gog_urls.AddToWishlist,
		addError,
		addProgress,
		httpClient,
		vrStoreProducts,
		titleExtracts)
}

func wishlistRemove(
	ids []string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_values.ValueReader,
	mt gog_media.Media,
	titleExtracts *froth.Stash) error {

	if err := removeData(ids, vangogh_products.WishlistProducts, mt); err != nil {
		return err
	}

	return remoteWishlistCommand(
		ids,
		gog_urls.RemoveFromWishlist,
		removeError,
		removeProgress,
		httpClient,
		vrStoreProducts,
		titleExtracts)
}

func remoteWishlistCommand(
	ids []string,
	wishlistUrl func(string) *url.URL,
	fmtError func(string) string,
	fmtProgress func(string, string) string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_values.ValueReader,
	titleExtracts *froth.Stash) error {
	for _, id := range ids {
		if !vrStoreProducts.Contains(id) {
			log.Printf("vangogh: %s", fmtError(id))
			continue
		}
		title, ok := titleExtracts.Get(id)
		if !ok {
			title = id
		}
		fmt.Println(fmtProgress(title, id))
		wUrl := wishlistUrl(id)
		resp, err := httpClient.Get(wUrl.String())
		if err != nil {
			return err
		}

		if err := resp.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}
