package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/redux"
	"strconv"
)

const dehydratedCount = 3

func ProductsList(r compton.Registrar, ids []string, from, to int, rdx redux.Readable, topTarget bool) compton.Element {
	productCards := compton.GridItems(r).JustifyContent(align.Center)

	if (to - from) < 10 {
		productCards.AddClass("items-" + strconv.Itoa(to-from))
	}

	for ii := from; ii < to; ii++ {
		id := ids[ii]
		productLink := compton.A("/product?id=" + id)
		if topTarget {
			productLink.SetAttribute("target", "_top")
		}

		if productCard := ProductCard(r, id, ii-from < dehydratedCount, rdx); productCard != nil {
			productLink.Append(productCard)
			productCards.Append(productLink)
		}
	}

	return productCards
}
