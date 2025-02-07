package cli

import (
	"github.com/arelate/southern_light/vangogh_integration"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
	"net/url"
)

const limitVideoRequests = 1000

func GetVideoMetadataHandler(u *url.URL) error {
	ids, err := vangogh_integration.IdsFromUrl(u)
	if err != nil {
		return err
	}

	return GetVideoMetadata(
		ids,
		vangogh_integration.FlagFromUrl(u, "missing"),
		vangogh_integration.FlagFromUrl(u, "force"))
}

func GetVideoMetadata(ids []string, missing, force bool) error {

	gvma := nod.NewProgress("getting video metadata...")
	defer gvma.End()

	rdx, err := videoReduxAssets()
	if err != nil {
		return gvma.EndWithError(err)
	}

	videoIds := make([]string, 0, len(ids))
	for _, id := range ids {
		if vip, ok := rdx.GetAllValues(vangogh_integration.VideoIdProperty, id); ok {
			for _, vid := range vip {
				if rdx.HasKey(vangogh_integration.VideoTitleProperty, vid) && !force {
					continue
				}
				if rdx.HasKey(vangogh_integration.VideoErrorProperty, vid) && !force {
					continue
				}
				videoIds = append(videoIds, vid)
			}
		}
	}

	if missing {
		for id := range rdx.Keys(vangogh_integration.VideoIdProperty) {
			if vip, ok := rdx.GetAllValues(vangogh_integration.VideoIdProperty, id); ok {
				for _, vid := range vip {
					if rdx.HasKey(vangogh_integration.VideoErrorProperty, vid) && !force {
						continue
					}
					if !rdx.HasKey(vangogh_integration.VideoTitleProperty, vid) {
						videoIds = append(videoIds, vid)
					}
				}
			}
		}
	}

	if len(videoIds) > limitVideoRequests {
		gvma.EndWithResult("limiting number of videos to avoid IP blacklisting")
		gvma = nod.NewProgress("getting %d videos metadata...", limitVideoRequests)
		videoIds = videoIds[:limitVideoRequests]
	}

	gvma.TotalInt(len(videoIds))
	videoTitles := make(map[string][]string)
	videoDurations := make(map[string][]string)
	videoErrors := make(map[string][]string)

	for _, videoId := range videoIds {

		ipr, err := youtube_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			videoErrors[videoId] = append(videoErrors[videoId], err.Error())
			gvma.Error(err)
			gvma.Increment()
			continue
		}

		videoTitles[videoId] = []string{ipr.VideoDetails.Title}
		videoDurations[videoId] = []string{ipr.VideoDetails.LengthSeconds}

		gvma.Increment()
	}

	if err := rdx.BatchAddValues(vangogh_integration.VideoTitleProperty, videoTitles); err != nil {
		return gvma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(vangogh_integration.VideoDurationProperty, videoDurations); err != nil {
		return gvma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(vangogh_integration.VideoErrorProperty, videoErrors); err != nil {
		return gvma.EndWithError(err)
	}

	gvma.EndWithResult("done")

	return nil
}

func videoReduxAssets() (redux.Writeable, error) {

	propSet := make(map[string]bool)
	propSet[vangogh_integration.TitleProperty] = true

	for _, vp := range vangogh_integration.VideoProperties() {
		propSet[vp] = true
	}

	properties := make([]string, len(propSet))
	for p := range propSet {
		properties = append(properties, p)
	}

	return vangogh_integration.NewReduxWriter(properties...)
}
