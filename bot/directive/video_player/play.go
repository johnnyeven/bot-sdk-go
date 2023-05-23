package video_player

import (
	"github.com/johnnyeven/bot-sdk-go/bot/directive"
)

var behaviorMap = map[string]bool{
	ENQUEUE:          true,
	REPLACE_ALL:      true,
	REPLACE_ENQUEUED: true,
}

type PlayDirective struct {
	directive.BaseDirective
	PlayBehavior string `json:"playBehavior"`
	VideoItem    struct {
		Stream struct {
			Url                  string `json:"url"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
			ExpiryTime           string `json:"expiryTime,omitempty"`
			ProgressReport       struct {
				ProgressReportDelayInMilliseconds    int `json:"progressReportDelayInMilliseconds,omitempty"`
				ProgressReportIntervalInMilliseconds int `json:"progressReportIntervalInMilliseconds,omitempty"`
			} `json:"progressReport,omitempty"`
			Token                 string `json:"token"`
			ExpectedPreviousToken string `json:"expectedPreviousToken,omitempty"`
		} `json:"stream"`
	} `json:"VideoItem"`
}

func NewPlayDirective(url string) *PlayDirective {
	play := &PlayDirective{}
	play.Type = "VideoPlayer.Play"
	play.PlayBehavior = REPLACE_ALL
	play.VideoItem.Stream.Url = url
	play.VideoItem.Stream.OffsetInMilliseconds = 0
	play.VideoItem.Stream.Token = play.GenToken()
	return play
}

func (d *PlayDirective) SetBehavior(behavior string) *PlayDirective {
	_, ok := behaviorMap[behavior]
	if ok {
		d.PlayBehavior = behavior
	}

	return d
}

func (d *PlayDirective) SetToken(token string) *PlayDirective {
	d.VideoItem.Stream.Token = token
	return d
}

func (d *PlayDirective) GetToken(token string) string {
	return d.VideoItem.Stream.Token
}

func (d *PlayDirective) SetUrl(url string) *PlayDirective {
	d.VideoItem.Stream.Url = url
	return d
}

func (d *PlayDirective) SetOffsetInMilliseconds(milliseconds int) *PlayDirective {
	d.VideoItem.Stream.OffsetInMilliseconds = milliseconds
	return d
}

func (d *PlayDirective) SetExpiryTime(expiryTime string) *PlayDirective {
	d.VideoItem.Stream.ExpiryTime = expiryTime
	return d
}

func (d *PlayDirective) SetReportDelayInMs(reportDelayInMs int) *PlayDirective {
	d.VideoItem.Stream.ProgressReport.ProgressReportDelayInMilliseconds = reportDelayInMs
	return d
}

func (d *PlayDirective) SetReportIntervalInMs(reportIntervalInMs int) *PlayDirective {
	d.VideoItem.Stream.ProgressReport.ProgressReportIntervalInMilliseconds = reportIntervalInMs
	return d
}

func (d *PlayDirective) SetExpectedPreviousToken(expectedPreviousToken string) *PlayDirective {
	d.VideoItem.Stream.ExpectedPreviousToken = expectedPreviousToken
	return d
}
