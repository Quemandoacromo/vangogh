package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/boggydigital/vangogh/internal/gog/const/aliases"
	"github.com/boggydigital/vangogh/internal/gog/const/names"
	"go.mongodb.org/mongo-driver/mongo"
)

type Setter interface {
	Set(data interface{}) error
}

func GetSetterByName(name string, mongoClient *mongo.Client, ctx context.Context) (Setter, error) {
	switch name {
	case aliases.Products:
		fallthrough
	case names.Products:
		return NewProducts(mongoClient, ctx), nil
	case aliases.AccountProducts:
		fallthrough
	case names.AccountProducts:
		return NewAccountProducts(mongoClient, ctx), nil
	case aliases.Wishlist:
		fallthrough
	case names.Wishlist:
		return NewWishlist(mongoClient, ctx), nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown source: %s", name))
	}
}
