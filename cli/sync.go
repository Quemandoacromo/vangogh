package cli

import (
	"errors"
	"maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/arelate/vangogh/cli/shared_data"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"

	"github.com/arelate/southern_light/vangogh_integration"
	"github.com/boggydigital/nod"
)

const (
	SyncOptionPurchases         = "purchases"
	SyncOptionDescriptionImages = "description-images"
	SyncOptionImages            = "images"
	SyncOptionScreenshots       = "screenshots"
	SyncOptionVideosMetadata    = "videos-metadata"
	SyncOptionDownloadsUpdates  = "downloads-updates"
	SynOptionWineBinaries       = "wine-binaries"
	negativePrefix              = "no-"
)

type syncOptions struct {
	purchases         bool
	descriptionImages bool
	images            bool
	screenshots       bool
	videosMetadata    bool
	downloadsUpdates  bool
	wineBinaries      bool
}

func NegOpt(option string) string {
	if !strings.HasPrefix(option, negativePrefix) {
		return negativePrefix + option
	}
	return option
}

func initSyncOptions(u *url.URL) *syncOptions {

	so := &syncOptions{
		purchases:         vangogh_integration.FlagFromUrl(u, SyncOptionPurchases),
		descriptionImages: vangogh_integration.FlagFromUrl(u, SyncOptionDescriptionImages),
		images:            vangogh_integration.FlagFromUrl(u, SyncOptionImages),
		screenshots:       vangogh_integration.FlagFromUrl(u, SyncOptionScreenshots),
		videosMetadata:    vangogh_integration.FlagFromUrl(u, SyncOptionVideosMetadata),
		downloadsUpdates:  vangogh_integration.FlagFromUrl(u, SyncOptionDownloadsUpdates),
		wineBinaries:      vangogh_integration.FlagFromUrl(u, SynOptionWineBinaries),
	}

	if vangogh_integration.FlagFromUrl(u, "all") {
		// not handling purchases here - doesn't make sense to negate it
		so.descriptionImages = !vangogh_integration.FlagFromUrl(u, NegOpt(SyncOptionDescriptionImages))
		so.images = !vangogh_integration.FlagFromUrl(u, NegOpt(SyncOptionImages))
		so.screenshots = !vangogh_integration.FlagFromUrl(u, NegOpt(SyncOptionScreenshots))
		so.videosMetadata = !vangogh_integration.FlagFromUrl(u, NegOpt(SyncOptionVideosMetadata))
		so.downloadsUpdates = !vangogh_integration.FlagFromUrl(u, NegOpt(SyncOptionDownloadsUpdates))
		so.wineBinaries = !vangogh_integration.FlagFromUrl(u, NegOpt(SynOptionWineBinaries))
	}

	return so
}

func SyncHandler(u *url.URL) error {
	syncOpts := initSyncOptions(u)

	since, err := vangogh_integration.SinceFromUrl(u)
	if err != nil {
		return err
	}

	q := u.Query()

	return Sync(
		since,
		syncOpts,
		vangogh_integration.OperatingSystemsFromUrl(u),
		vangogh_integration.LanguageCodesFromUrl(u),
		vangogh_integration.DownloadTypesFromUrl(u),
		q.Has("no-patches"),
		vangogh_integration.DownloadsLayoutFromUrl(u),
		!q.Has("no-cleanup"),
		q.Has("force"))
}

func Sync(
	since int64,
	syncOpts *syncOptions,
	operatingSystems []vangogh_integration.OperatingSystem,
	langCodes []string,
	downloadTypes []vangogh_integration.DownloadType,
	noPatches bool,
	downloadsLayout vangogh_integration.DownloadsLayout,
	cleanup bool,
	force bool) error {

	sa := nod.Begin("syncing data...")
	defer sa.Done()

	syncStart := since
	if syncStart == 0 {
		syncStart = time.Now().Unix()
	}

	reduxDir, err := pathways.GetAbsRelDir(vangogh_integration.Redux)
	if err != nil {
		return err
	}

	syncEventsRdx, err := redux.NewWriter(reduxDir, vangogh_integration.SyncEventsProperty)
	if err != nil {
		return err
	}

	if err = setSyncEvent(vangogh_integration.SyncStartKey, syncEventsRdx); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	if err = GetData(nil, nil, since, syncOpts.purchases, true, force); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	if err = setSyncEvent(vangogh_integration.SyncDataKey, syncEventsRdx); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	if err = Reduce([]vangogh_integration.ProductType{vangogh_integration.UnknownProductType}); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	// summarize sync updates now, since other updates are digital artifacts
	// and won't affect the summaries
	if err = Summarize(syncStart); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	// get description images
	if syncOpts.descriptionImages {
		if err = GetDescriptionImages(nil, since, force); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = setSyncEvent(vangogh_integration.SyncDescriptionImagesKey, syncEventsRdx); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}
	}

	// get images
	if syncOpts.images {
		imageTypes := make([]vangogh_integration.ImageType, 0, len(vangogh_integration.AllImageTypes()))
		for _, it := range vangogh_integration.AllImageTypes() {
			if !syncOpts.screenshots && it == vangogh_integration.Screenshots {
				continue
			}
			imageTypes = append(imageTypes, it)
		}
		if err = GetImages(nil, imageTypes, true, force); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = setSyncEvent(vangogh_integration.SyncImagesKey, syncEventsRdx); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		dehydratedImageTypes := []vangogh_integration.ImageType{vangogh_integration.Image, vangogh_integration.VerticalImage}
		if err = Dehydrate(nil, dehydratedImageTypes, false); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = setSyncEvent(vangogh_integration.SyncDehydrateKey, syncEventsRdx); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}
	}

	if syncOpts.videosMetadata {
		if err = GetVideoMetadata(nil, true, false); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = setSyncEvent(vangogh_integration.SyncVideoMetadataKey, syncEventsRdx); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}
	}

	// get downloads updates
	if syncOpts.downloadsUpdates {

		updatedDetails, err := shared_data.GetDetailsUpdates(since)
		if err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = queueDownloads(maps.Keys(updatedDetails)); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		// process remaining queued downloads, e.g. downloads that were queued before this sync
		// and were not completed successfully due to an interruption. Download updates
		// itemized earlier in the sync cycle (ids) are intentionally excluded to
		// focus on previously queued and avoid attempting to download problematic ids
		// right after they didn't download successfully, waiting until the next sync
		// is likely a better strategy in that case
		if err = ProcessQueue(
			operatingSystems,
			langCodes,
			downloadTypes,
			noPatches,
			downloadsLayout); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = setSyncEvent(vangogh_integration.SyncDownloadsKey, syncEventsRdx); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if cleanup {
			if err = Cleanup(nil,
				operatingSystems,
				langCodes,
				downloadTypes,
				noPatches,
				downloadsLayout,
				true,
				false); err != nil {
				return setSyncInterrupted(err, syncEventsRdx)
			}

			if err = setSyncEvent(vangogh_integration.SyncCleanupKey, syncEventsRdx); err != nil {
				return setSyncInterrupted(err, syncEventsRdx)
			}
		}

	}

	if syncOpts.wineBinaries {
		if err = GetWineBinaries(operatingSystems, force); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}

		if err = setSyncEvent(vangogh_integration.SyncWineBinaries, syncEventsRdx); err != nil {
			return setSyncInterrupted(err, syncEventsRdx)
		}
	}

	// backing up data
	if err = Backup(); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	if err = setSyncEvent(vangogh_integration.SyncBackup, syncEventsRdx); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	// print new, updated
	if err = GetSummary(); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	if err = setSyncEvent(vangogh_integration.SyncCompleteKey, syncEventsRdx); err != nil {
		return setSyncInterrupted(err, syncEventsRdx)
	}

	return nil
}

func setSyncEvent(eventKey string, rdx redux.Writeable) error {

	if err := rdx.MustHave(vangogh_integration.SyncEventsProperty); err != nil {
		return err
	}

	if !slices.Contains(vangogh_integration.SyncEventsKeys, eventKey) {
		return errors.New("unknown sync event key: " + eventKey)
	}

	syncEvents := make(map[string][]string)
	syncEvents[eventKey] = []string{strconv.Itoa(int(time.Now().Unix()))}

	return rdx.BatchReplaceValues(vangogh_integration.SyncEventsProperty, syncEvents)
}

func setSyncInterrupted(err error, rdx redux.Writeable) error {
	if sse := setSyncEvent(vangogh_integration.SyncInterruptedKey, rdx); sse != nil {
		return errors.Join(sse, err)
	}

	return err
}
