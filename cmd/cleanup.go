package cmd

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/vangogh/cmd/url_helpers"
	"net/url"
	"os"
	"path/filepath"
)

const (
	dirPerm os.FileMode = 0755
)

// TODO: implement this in a cleaner way
var testMode bool

func CleanupHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	operatingSystems := url_helpers.OperatingSystems(u)
	langCodes := url_helpers.Values(u, "language-code")
	downloadTypes := url_helpers.DownloadTypes(u)

	all := url_helpers.Flag(u, "all")

	testMode = url_helpers.Flag(u, "test")

	return Cleanup(idSet, mt, operatingSystems, langCodes, downloadTypes, all)
}

func Cleanup(
	idSet gost.StrSet,
	mt gog_media.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	langCodes []string,
	downloadTypes []vangogh_downloads.DownloadType,
	all bool) error {

	exl, err := vangogh_extracts.NewList(
		vangogh_properties.SlugProperty,
		vangogh_properties.NativeLanguageNameProperty,
		vangogh_properties.LocalManualUrl)
	if err != nil {
		return err
	}

	if all {
		vrDetails, err := vangogh_values.NewReader(vangogh_products.Details, mt)
		if err != nil {
			return err
		}
		idSet.Add(vrDetails.All()...)
	}

	cd := &cleanupDelegate{exl: exl}

	if err := vangogh_downloads.Map(
		idSet,
		mt,
		exl,
		operatingSystems,
		downloadTypes,
		langCodes,
		cd.CleanupList,
		0,
		false); err != nil {
		return err
	}

	return nil
}

func moveToRecycleBin(fp string) error {
	if testMode {
		return nil
	}
	rbFilepath := filepath.Join(vangogh_urls.RecycleBinDir(), fp)
	rbDir, _ := filepath.Split(rbFilepath)
	if _, err := os.Stat(rbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rbDir, dirPerm); err != nil {
			return err
		}
	}
	return os.Rename(fp, rbFilepath)
}

type cleanupDelegate struct {
	exl *vangogh_extracts.ExtractsList
}

func (cd *cleanupDelegate) CleanupList(_ string, slug string, list vangogh_downloads.DownloadsList) error {

	if err := cd.exl.AssertSupport(vangogh_properties.LocalManualUrl); err != nil {
		return err
	}

	fmt.Printf("cleaning up %s\n", slug)

	//cleanup process:
	//1. enumerate all expected files for a downloadList
	//2. enumerate all files present for a slug (files present in a `downloads/slug` folder)
	//3. delete (present files).Except(expected files) and corresponding xml files

	expectedSet := gost.NewStrSet()

	//pDir = s/slug
	pDir, err := vangogh_urls.ProductDownloadsRelDir(slug)
	if err != nil {
		return err
	}

	for _, dl := range list {
		if localFilename, ok := cd.exl.Get(vangogh_properties.LocalManualUrl, dl.ManualUrl); ok {
			//local filenames are saved as relative to root downloads folder (e.g. s/slug/local_filename)
			//so filepath.Rel would trim to local_filename (or dlc/local_filename, extra/local_filename)
			relFilename, err := filepath.Rel(pDir, localFilename)
			if err != nil {
				return err
			}
			expectedSet.Add(relFilename)
		}
	}

	//LocalSlugDownloads returns list of files relative to s/slug product directory
	presentSet, err := vangogh_urls.LocalSlugDownloads(slug)
	if err != nil {
		return err
	}

	unexpectedFiles := presentSet.Except(expectedSet)
	if len(unexpectedFiles) == 0 {
		fmt.Println(" already clean")
		return nil
	}

	for _, unexpectedFile := range unexpectedFiles {
		//restore absolute from local_filename to s/slug/local_filename
		downloadFilename := vangogh_urls.DownloadRelToAbs(filepath.Join(pDir, unexpectedFile))
		if _, err := os.Stat(downloadFilename); os.IsNotExist(err) {
			continue
		}
		prefix := ""
		if testMode {
			prefix = " TEST"
		}
		fmt.Println(prefix, downloadFilename)
		if err := moveToRecycleBin(downloadFilename); err != nil {
			return err
		}

		checksumFile := vangogh_urls.LocalChecksumPath(unexpectedFile)
		if _, err := os.Stat(checksumFile); os.IsNotExist(err) {
			continue
		}
		fmt.Println(prefix, checksumFile)
		if err := moveToRecycleBin(checksumFile); err != nil {
			return err
		}
	}

	return nil
}
