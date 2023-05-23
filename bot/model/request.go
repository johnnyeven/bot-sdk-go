package model

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/johnnyeven/bot-sdk-go/bot/data"
)

const (
	INENT_REQUEST         = "IntentRequest"
	LAUNCH_REQUEST        = "LaunchRequest"
	SESSION_ENDED_REQUEST = "SessionEndedRequest"

	AUDIO_PLAYER_PLAYBACK_STARTED                 = "AudioPlayer.PlaybackStarted"
	AUDIO_PLAYER_PLAYBACK_STOPPED                 = "AudioPlayer.PlaybackStopped"
	AUDIO_PLAYER_PLAYBACK_FINISHED                = "AudioPlayer.PlaybackFinished"
	AUDIO_PLAYER_PLAYBACK_NEARLY_FINISHED         = "AudioPlayer.PlaybackNearlyFinished"
	AUDIO_PLAYER_PROGRESS_REPORT_INTERVAL_ELAPSED = "AudioPlayer.ProgressReportIntervalElapsed"

	VIDEO_PLAYER_PLAYBACK_STARTED                 = "VideoPlayer.PlaybackStarted"
	VIDEO_PLAYER_PLAYBACK_STOPPED                 = "VideoPlayer.PlaybackStopped"
	VIDEO_PLAYER_PLAYBACK_FINISHED                = "VideoPlayer.PlaybackFinished"
	VIDEO_PLAYER_PLAYBACK_NEARLY_FINISHED         = "VideoPlayer.PlaybackNearlyFinished"
	VIDEO_PLAYER_PLAYBACK_SCHEDULED_STOP_REACHED  = "VideoPlayer.PlaybackScheduledStopReached"
	VIDEO_PLAYER_PROGRESS_REPORT_INTERVAL_ELAPSED = "VideoPlayer.ProgressReportIntervalElapsed"
)

type Request struct {
	Type   string
	Common data.RequestPart
}

type IntentRequest struct {
	Data   data.IntentRequest
	Dialog *Dialog
	Request
}

type LaunchRequest struct {
	Data data.LaunchRequest
	Request
}

type SessionEndedRequest struct {
	Data data.SessionEndedRequest
	Request
}

type EventRequest struct {
	Data data.EventRequest
	Request
}

type AudioPlayerEventRequest struct {
	Data data.AudioPlayerEventRequest
	EventRequest
}

type VideoPlayerEventRequest struct {
	Data data.VideoPlayerEventRequest
	EventRequest
}

func (this *EventRequest) GetUrl() string {
	return this.Data.Request.Url
}

func (this *EventRequest) GetName() string {
	return this.Data.Request.Name
}

func (this *AudioPlayerEventRequest) GetOffsetInMilliseconds() int32 {
	return this.Data.Request.OffsetInMilliseconds
}

func (this *VideoPlayerEventRequest) GetOffsetInMilliseconds() int32 {
	return this.Data.Request.OffsetInMilliseconds
}

// 获取意图名
func (r *IntentRequest) GetIntentName() (string, bool) {
	return r.Dialog.GetIntentName()
}

// 槽位填充是否完成
func (r *IntentRequest) IsDialogStateCompleted() bool {
	return r.Dialog.DialogState == "COMPLETED"
}

// 获取用户请求query
func (r *IntentRequest) GetQuery() string {
	query, _ := r.Dialog.GetQuery()
	return query
}

// 获取用户id
func (r *Request) GetUserId() string {
	return r.Common.Context.System.User.UserId
}

// 获取设备id
func (r *Request) GetDeviceId() string {
	return r.Common.Context.System.Device.DeviceId
}

// 获取音频播放上下文
func (r *Request) GetAudioPlayerContext() data.AudioPlayerContext {
	return r.Common.Context.AudioPlayer
}

// 获取视频播放上下文
func (r *Request) GetVideoPlayerContext() data.VideoPlayerContext {
	return r.Common.Context.VideoPlayer
}

// 获取access token
func (r *Request) GetAccessToken() string {
	return r.Common.Context.System.User.AccessToken
}

// 获取请求的时间戳
func (r *Request) GetTimestamp() int {
	i, err := strconv.Atoi(r.Common.Request.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// 获取请求id
func (r *Request) GetRequestId() string {
	return r.Common.Request.RequestId
}

// 获取技能id
func (r *Request) GetBotId() string {
	return r.Common.Context.System.Application.ApplicationId
}

// 验证请求时间戳合法性
func (r *Request) VerifyTimestamp() bool {

	if r.GetTimestamp()+180 > int(time.Now().Unix()) {
		return true
	}

	return false
}

// 获取设备支持的接口类型
func (r *Request) GetSupportedInterfaces() map[string]interface{} {
	return r.Common.Context.System.Device.SupportedInterfaces
}

func (r *Request) isSupportInterface(support string) bool {
	supportedInterfaces := r.GetSupportedInterfaces()
	_, ok := supportedInterfaces[support]

	if ok {
		return true
	}
	return false
}

// 检查是否支持展现
func (r *Request) IsSupportDisplay() bool {
	return r.isSupportInterface("Display")
}

// 检查是否支持音频播放
func (r *Request) IsSupportAudio() bool {
	return r.isSupportInterface("AudioPlayer")
}

// 检查是否支持视频播放
func (r *Request) IsSupportVideo() bool {
	return r.isSupportInterface("VideoPlayer")
}

// 验证技能id合法性
func (r *Request) VerifyBotID(myBotID string) bool {
	if r.GetBotId() == myBotID {
		return true
	}
	return false
}

func getType(rawData string) string {
	jsonBlob := []byte(rawData)
	d := data.LaunchRequest{}

	if err := json.Unmarshal(jsonBlob, &d); err != nil {
		log.Println(err)
	}

	return d.Request.Type
}

func GetSessionData(rawData string) data.Session {
	jsonBlob := []byte(rawData)
	common := data.RequestPart{}
	if err := json.Unmarshal(jsonBlob, &common); err != nil {
		log.Println(err)
	}

	return common.Session
}

func NewRequest(rawData string) interface{} {
	requestType := getType(rawData)

	jsonBlob := []byte(rawData)

	common := data.RequestPart{}
	if err := json.Unmarshal(jsonBlob, &common); err != nil {
		log.Println(err)
		return false
	}

	if requestType == INENT_REQUEST {
		request := IntentRequest{}
		request.Type = requestType
		request.Common = common
		if err := json.Unmarshal(jsonBlob, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		request.Dialog = NewDialog(request.Data.Request)

		return request
	} else if requestType == LAUNCH_REQUEST {
		request := LaunchRequest{}
		request.Type = requestType
		request.Common = common
		if err := json.Unmarshal(jsonBlob, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		return request
	} else if requestType == SESSION_ENDED_REQUEST {
		request := SessionEndedRequest{}
		request.Type = requestType
		request.Common = common
		if err := json.Unmarshal(jsonBlob, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		return request
	} else {
		if match, _ := regexp.MatchString("^AudioPlayer", requestType); match {
			request := AudioPlayerEventRequest{}
			request.Type = requestType
			request.Common = common
			if err := json.Unmarshal(jsonBlob, &request.Data); err != nil {
				log.Println(err)
				return false
			}
			return request
		} else if match, _ := regexp.MatchString("^VideoPlayer", requestType); match {
			request := VideoPlayerEventRequest{}
			request.Type = requestType
			request.Common = common
			if err := json.Unmarshal(jsonBlob, &request.Data); err != nil {
				log.Println(err)
				return false
			}
			return request
		}

		request := EventRequest{}
		request.Type = requestType
		request.Common = common
		if err := json.Unmarshal(jsonBlob, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		return request
	}
	return false
}
