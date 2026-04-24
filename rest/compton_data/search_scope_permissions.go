package compton_data

import (
	"github.com/arelate/vangogh/perm"
	"github.com/boggydigital/author"
)

var SearchScopePermissions = map[string]author.Permission{
	SearchGog:      perm.ReadOwned,
	SearchWishlist: perm.ReadWishlist,
}
