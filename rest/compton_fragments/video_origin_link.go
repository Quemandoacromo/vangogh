package compton_fragments

import (
	"strconv"
	"time"

	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

func formatSeconds(ts int64) string {
	if ts == 0 {
		return "unknown"
	}

	t := time.Unix(ts, 0).UTC()

	layout := "4:05"
	if t.Hour() > 0 {
		layout = "15:04:05"
	}

	return t.Format(layout)
}

func VideoOriginLink(r compton.Registrar, videoId, videoTitle, videoDuration string) compton.Element {

	originLink := compton.A(youtube_urls.VideoUrl(videoId).String())
	originLink.SetAttribute("target", "_top")

	linkColumn := compton.FlexItems(r, direction.Column).
		RowGap(size.Small)

	if videoTitle == "" {
		videoTitle = "Watch at origin"
	}

	// linkRow := compton.FlexItems(r, direction.Row).
	// 	AlignItems(align.Center).
	// 	JustifyContent(align.Center).
	// 	ColumnGap(size.Small)
	// linkColumn.Append(linkRow)

	// linkRow.Append(compton.SvgUse(r, compton.VideoThumbnail).ForegroundColor(color.Cyan))

	linkText := compton.Fspan(r, videoTitle).
		TextAlign(align.Center).
		FontWeight(font_weight.Bolder).
		ForegroundColor(color.Cyan)
	linkColumn.Append(linkText)

	if dur, err := strconv.ParseInt(videoDuration, 10, 64); err == nil {
		durFrow := compton.Frow(r)
		durFrow.IconColor(r, compton.VideoThumbnail, color.RepGray)
		durFrow.FontSize(size.XSmall)
		durFrow.PropVal("Duration", formatSeconds(dur))
		linkColumn.Append(compton.FICenter(r, durFrow))
	}

	originLink.Append(linkColumn)

	return originLink
}
