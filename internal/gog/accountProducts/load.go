package accountProducts

import (
	"encoding/json"
	media "github.com/boggydigital/vangogh/internal/gog/media"
	"github.com/boggydigital/vangogh/internal/gog/paths"
	"github.com/boggydigital/vangogh/internal/storage"
)

func Load(id int, mt media.Type) (ap *AccountProduct, err error) {
	apBytes, err := storage.Load(paths.AccountProduct(id, mt))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(apBytes, &ap)

	return ap, err
}
