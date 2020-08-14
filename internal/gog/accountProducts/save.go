package accountProducts

import (
	media "github.com/boggydigital/vangogh/internal/gog/media"
	"github.com/boggydigital/vangogh/internal/gog/paths"
	"github.com/boggydigital/vangogh/internal/storage"
)

func Save(ap *AccountProduct, mt media.Type) error {
	return storage.Save(ap, paths.AccountProduct(ap.ID, mt))
}
