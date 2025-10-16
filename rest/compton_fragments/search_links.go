package compton_fragments

import (
	"slices"

	"github.com/arelate/vangogh/rest/compton_data"
	"github.com/boggydigital/author"
	"github.com/boggydigital/compton"
)

func SearchLinks(r compton.Registrar, current string, permissions ...author.Permission) compton.Element {

	searchNavLinks := compton.NavLinks(r)
	searchNavLinks.SetAttribute("style", "view-transition-name:secondary-nav")

	searchScopes := compton_data.SearchScopes()

	for _, scope := range compton_data.SearchOrder {

		if prm, ok := compton_data.SearchScopePermissions[scope]; ok && !slices.Contains(permissions, prm) {
			continue
		}

		searchLink := searchNavLinks.AppendLink(r, &compton.NavTarget{
			Href:     "/search?" + searchScopes[scope],
			Title:    scope,
			Selected: current == scope,
			Symbol:   compton_data.SearchScopesSymbols[scope],
		})
		if current == scope {
			searchLink.SetAttribute("style", "view-transition-name:current-search-link")
		}
	}

	return searchNavLinks
}
