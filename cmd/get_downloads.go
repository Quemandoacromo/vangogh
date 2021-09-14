package cmd

import (
	"fmt"
	"github.com/arelate/gog_auth"
	"github.com/arelate/gog_media"
	"github.com/arelate/gog_urls"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/vangogh/cmd/http_client"
	"github.com/boggydigital/vangogh/cmd/itemize"
	"github.com/boggydigital/vangogh/cmd/url_helpers"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func GetDownloadsHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	operatingSystems := url_helpers.OperatingSystems(u)
	downloadTypes := url_helpers.DownloadTypes(u)
	langCodes := url_helpers.Values(u, "language-code")

	missing := url_helpers.Flag(u, "missing")

	forceUpdate := url_helpers.Flag(u, "force-update")

	return GetDownloads(
		idSet,
		mt,
		operatingSystems,
		downloadTypes,
		langCodes,
		missing,
		forceUpdate)
}

func GetDownloads(
	idSet gost.StrSet,
	mt gog_media.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	missing,
	forceUpdate bool) error {

	httpClient, err := http_client.Default()
	if err != nil {
		return err
	}

	li, err := gog_auth.LoggedIn(httpClient)
	if err != nil {
		return err
	}

	if !li {
		log.Fatalf("user is not logged in")
	}

	exl, err := vangogh_extracts.NewList(
		vangogh_properties.NativeLanguageNameProperty,
		vangogh_properties.SlugProperty,
		vangogh_properties.LocalManualUrl,
		vangogh_properties.DownloadStatusError)
	if err != nil {
		return err
	}

	if missing {
		missingIds, err := itemize.MissingLocalDownloads(mt, exl, operatingSystems, downloadTypes, langCodes)
		if err != nil {
			return err
		}

		if missingIds.Len() == 0 {
			fmt.Println("all downloads are available locally")
			return nil
		}

		idSet.AddSet(missingIds)
	}

	gdd := &getDownloadsDelegate{
		exl:         exl,
		forceUpdate: forceUpdate,
	}

	if err := vangogh_downloads.Map(
		idSet,
		mt,
		exl,
		operatingSystems,
		downloadTypes,
		langCodes,
		gdd.DownloadList); err != nil {
		return nil
	}

	return nil
}

func printCompletion(current, total uint64) {
	percent := (float32(current) / float32(total)) * 100
	if current == 0 {
		//we'll get the first notification before download starts and will output 4 spaces (XXX%)
		//that'll be deleted on updates
		fmt.Printf(strings.Repeat(" ", 4))
	}
	if current < total {
		//every update except the first (pre-download) and last (completion) are the same
		//move cursor 4 spaces back and print over current percent completion
		fmt.Printf("\x1b[4D%3.0f%%", percent)
	} else {
		//final update moves the cursor back 4 spaces to overwrite on the following update
		fmt.Printf("\x1b[4D")
	}
}

type getDownloadsDelegate struct {
	exl         *vangogh_extracts.ExtractsList
	forceUpdate bool
}

func (gdd *getDownloadsDelegate) DownloadList(_ string, slug string, list vangogh_downloads.DownloadsList) error {
	fmt.Println("downloading", slug)

	if len(list) == 0 {
		fmt.Println(" (no downloads for requested operating systems + download types + languages)")
		return nil
	}

	httpClient, err := http_client.Default()
	if err != nil {
		return err
	}

	//there is no need to use internal httpClient with cookie support for downloading
	//manual downloads, so we're going to rely on default http.Client
	defaultClient := http.DefaultClient
	dlClient := dolo.NewClient(defaultClient, printCompletion, dolo.Defaults())

	for _, dl := range list {
		if err := gdd.downloadManualUrl(slug, &dl, httpClient, dlClient); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (gdd *getDownloadsDelegate) downloadManualUrl(
	slug string,
	dl *vangogh_downloads.Download,
	httpClient *http.Client,
	dlClient *dolo.Client) error {
	//downloading a manual URL is the following set of steps:
	//1 - check if local file exists (based on manualUrl -> relative localFile association) before attempting to resolve manualUrl
	//2 - resolve the source URL to an actual session URL
	//3 - construct local relative dir and filename based on manualUrl type (installer, movie, dlc, extra)
	//4 - for a given set of extensions - download validation file
	//5 - download authorized session URL to a file
	//6 - set association from manualUrl to a resolved filename
	if err := gdd.exl.AssertSupport(
		vangogh_properties.LocalManualUrl,
		vangogh_properties.DownloadStatusError); err != nil {
		return err
	}

	//1
	if !gdd.forceUpdate {
		if localFilename, ok := gdd.exl.Get(vangogh_properties.LocalManualUrl, dl.ManualUrl); ok {
			pDir, err := vangogh_urls.ProductDownloadsAbsDir(slug)
			if err != nil {
				return err
			}
			if _, err := os.Stat(path.Join(pDir, localFilename)); !os.IsNotExist(err) {
				return nil
			}
		}
	}

	fmt.Printf(" %s...", dl.String())

	//2
	resp, err := httpClient.Head(gog_urls.ManualDownloadUrl(dl.ManualUrl).String())
	if err != nil {
		fmt.Println()
		return err
	}
	//check for error status codes and store them for the manualUrl to provide a hint that locally missing file
	//is not a problem that can be solved locally (it's a remote source error)
	if resp.StatusCode > 299 {
		if err := gdd.exl.Set(vangogh_properties.DownloadStatusError, dl.ManualUrl, strconv.Itoa(resp.StatusCode)); err != nil {
			return err
		}
		return fmt.Errorf(resp.Status)
	}

	resolvedUrl := resp.Request.URL

	if err := resp.Body.Close(); err != nil {
		fmt.Println()
		return err
	}

	//3
	_, filename := path.Split(resolvedUrl.Path)
	//ProductDownloadsAbsDir would return absolute dir path, e.g. downloads/s/slug
	pAbsDir, err := vangogh_urls.ProductDownloadsAbsDir(slug)
	if err != nil {
		return err
	}
	//we need to add suffix to a dir path, e.g. dlc, extras
	absDir := filepath.Join(pAbsDir, dl.DirSuffix())

	//4
	remoteChecksumPath := vangogh_urls.RemoteChecksumPath(resolvedUrl.Path)
	if remoteChecksumPath != "" {
		localChecksumPath := vangogh_urls.LocalChecksumPath(path.Join(absDir, filename))
		if _, err := os.Stat(localChecksumPath); os.IsNotExist(err) {
			fmt.Print("xml")
			originalPath := resolvedUrl.Path
			resolvedUrl.Path = remoteChecksumPath
			valDir, valFilename := path.Split(localChecksumPath)
			if _, err := dlClient.Download(
				resolvedUrl, valDir, valFilename); err != nil {
				return err
			}
			resolvedUrl.Path = originalPath
			fmt.Print("...")
		}
	}

	//5
	if _, err := dlClient.Download(resolvedUrl, absDir, filename); err != nil {
		return err
	}

	//6
	//ProductDownloadsRelDir would return relative (to downloads/ root) dir path, e.g. s/slug
	pRelDir, err := vangogh_urls.ProductDownloadsRelDir(slug)
	//we need to add suffix to a dir path, e.g. dlc, extras
	relDir := filepath.Join(pRelDir, dl.DirSuffix())
	if err != nil {
		return err
	}
	//store association for ManualUrl (/downloads/en0installer) to local file (s/slug/local_filename)
	if err := gdd.exl.Set(vangogh_properties.LocalManualUrl, dl.ManualUrl, path.Join(relDir, filename)); err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}
