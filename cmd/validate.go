package cmd

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/vangogh/cmd/iterate"
	"github.com/boggydigital/vangogh/cmd/validation"
	"io"
	"os"
	"path"
)

var (
	ErrMissingDownload        = errors.New("download file missing")
	ErrMissingValidationFile  = errors.New("validation file missing")
	ErrValidationNotSupported = errors.New("validation not supported")
	ErrValidationFailed       = errors.New("validation failed")
)

func Validate(
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

	validated := make(map[string]bool)
	missingDownload := make(map[string]bool)
	missingChecksum := make(map[string]bool)
	failed := make(map[string]bool)
	slugLastError := make(map[string]string)

	if err := iterate.DownloadsList(
		idSet,
		mt,
		exl,
		operatingSystems,
		downloadTypes,
		langCodes,
		func(slug string,
			list vangogh_downloads.DownloadsList,
			exl *vangogh_extracts.ExtractsList,
			_ bool) error {

			hasValidationTargets := false

			for _, dl := range list {
				if err := validateManualUrl(slug, &dl, exl); errors.Is(err, ErrValidationNotSupported) {
					continue
				} else if errors.Is(err, ErrMissingValidationFile) {
					missingChecksum[slug] = true
				} else if errors.Is(err, ErrMissingDownload) {
					missingDownload[slug] = true
				} else if errors.Is(err, ErrValidationFailed) {
					failed[slug] = true
				} else if err != nil {
					fmt.Println(err)
					slugLastError[slug] = err.Error()
					continue
				}
				// don't attempt to assess success for files that don't support validation
				hasValidationTargets = true
			}

			if hasValidationTargets &&
				!missingDownload[slug] &&
				!missingChecksum[slug] &&
				!failed[slug] &&
				slugLastError[slug] == "" {
				validated[slug] = true
			}

			return nil
		},
		0,
		false); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("validation summary:")
	fmt.Printf("%d product(s) successfully validated\n", len(validated))
	if len(missingDownload) > 0 {
		fmt.Printf("%d product(s) not downloaded\n", len(missingDownload))
		for slug := range missingDownload {
			fmt.Println("", slug)
		}
	}
	if len(missingChecksum) > 0 {
		fmt.Printf("%d product(s) without checksum:\n", len(missingChecksum))
		for slug := range missingChecksum {
			fmt.Println("", slug)
		}
	}
	if len(failed) > 0 {
		fmt.Printf("%d product(s) failed validation:\n", len(failed))
		for slug := range failed {
			fmt.Println("", slug)
		}
	}
	if len(slugLastError) > 0 {
		fmt.Printf("%d product(s) validation caused an error:\n", len(slugLastError))
		for slug, err := range slugLastError {
			fmt.Printf(" %s: %s\n", slug, err)
		}
	}

	return nil
}

type validationError struct {
}

func (ve validationError) Error() string {
	return "validation failed"
}

func validateManualUrl(
	slug string,
	dl *vangogh_downloads.Download,
	exl *vangogh_extracts.ExtractsList) error {

	if err := exl.AssertSupport(vangogh_properties.LocalManualUrl); err != nil {
		return err
	}

	relLocalFile, ok := exl.Get(vangogh_properties.LocalManualUrl, dl.ManualUrl)
	if !ok {
		return ErrMissingDownload
	}

	absLocalFile := path.Join(vangogh_urls.DownloadsDir(), relLocalFile)
	if !vangogh_urls.CanValidate(absLocalFile) {
		return ErrValidationNotSupported
	}

	if _, err := os.Stat(absLocalFile); os.IsNotExist(err) {
		return ErrMissingDownload
	}

	absValidationFile := vangogh_urls.LocalValidationPath(absLocalFile)

	if _, err := os.Stat(absValidationFile); os.IsNotExist(err) {
		return ErrMissingValidationFile
	}

	fmt.Printf("validing %s...", dl)

	valFile, err := os.Open(absValidationFile)
	if err != nil {
		return err
	}
	defer valFile.Close()

	var valData validation.File
	if err := xml.NewDecoder(valFile).Decode(&valData); err != nil {
		return err
	}

	sourceFile, err := os.Open(absLocalFile)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	h := md5.New()
	if _, err := io.Copy(h, sourceFile); err != nil {
		return err
	}
	sourceFileMD5 := fmt.Sprintf("%x", h.Sum(nil))

	if valData.MD5 == sourceFileMD5 {
		fmt.Println("ok")
	} else {
		fmt.Println("FAIL")
		return ErrValidationFailed
	}

	return nil
}
