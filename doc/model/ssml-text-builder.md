

# model
`import "github.com/johnnyeven/bot-sdk-go/bot/model"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [type SSMLTextBuilder](#SSMLTextBuilder)
  * [func NewSSMLTextBuilder() *SSMLTextBuilder](#NewSSMLTextBuilder)
  * [func (this *SSMLTextBuilder) AppendAudio(src string) *SSMLTextBuilder](#SSMLTextBuilder.AppendAudio)
  * [func (this *SSMLTextBuilder) AppendBackground(text string, src string, repeat bool) *SSMLTextBuilder](#SSMLTextBuilder.AppendBackground)
  * [func (this *SSMLTextBuilder) AppendPlainSpeech(text string) *SSMLTextBuilder](#SSMLTextBuilder.AppendPlainSpeech)
  * [func (this *SSMLTextBuilder) AppendSilence(time int) *SSMLTextBuilder](#SSMLTextBuilder.AppendSilence)
  * [func (this *SSMLTextBuilder) AppendSubstitution(text, alias string) *SSMLTextBuilder](#SSMLTextBuilder.AppendSubstitution)
  * [func (this *SSMLTextBuilder) ApplyBackground(src string, repeat bool) *SSMLTextBuilder](#SSMLTextBuilder.ApplyBackground)
  * [func (this *SSMLTextBuilder) Build() string](#SSMLTextBuilder.Build)


#### <a name="pkg-files">Package files</a>
[dialog.go](/src/github.com/johnnyeven/bot-sdk-go/bot/model/dialog.go) [intent.go](/src/github.com/johnnyeven/bot-sdk-go/bot/model/intent.go) [request.go](/src/github.com/johnnyeven/bot-sdk-go/bot/model/request.go) [response.go](/src/github.com/johnnyeven/bot-sdk-go/bot/model/response.go) [session.go](/src/github.com/johnnyeven/bot-sdk-go/bot/model/session.go) [ssml_builder.go](/src/github.com/johnnyeven/bot-sdk-go/bot/model/ssml_builder.go) 






## <a name="SSMLTextBuilder">type</a> [SSMLTextBuilder](/src/target/ssml_builder.go?s=239:292#L15)
``` go
type SSMLTextBuilder struct {
    // contains filtered or unexported fields
}
```






### <a name="NewSSMLTextBuilder">func</a> [NewSSMLTextBuilder](/src/target/ssml_builder.go?s=294:336#L19)
``` go
func NewSSMLTextBuilder() *SSMLTextBuilder
```




### <a name="SSMLTextBuilder.AppendAudio">func</a> (\*SSMLTextBuilder) [AppendAudio](/src/target/ssml_builder.go?s=522:591#L30)
``` go
func (this *SSMLTextBuilder) AppendAudio(src string) *SSMLTextBuilder
```



### <a name="SSMLTextBuilder.AppendBackground">func</a> (\*SSMLTextBuilder) [AppendBackground](/src/target/ssml_builder.go?s=2062:2162#L93)
``` go
func (this *SSMLTextBuilder) AppendBackground(text string, src string, repeat bool) *SSMLTextBuilder
```



### <a name="SSMLTextBuilder.AppendPlainSpeech">func</a> (\*SSMLTextBuilder) [AppendPlainSpeech](/src/target/ssml_builder.go?s=394:470#L23)
``` go
func (this *SSMLTextBuilder) AppendPlainSpeech(text string) *SSMLTextBuilder
```



### <a name="SSMLTextBuilder.AppendSilence">func</a> (\*SSMLTextBuilder) [AppendSilence](/src/target/ssml_builder.go?s=1707:1776#L79)
``` go
func (this *SSMLTextBuilder) AppendSilence(time int) *SSMLTextBuilder
```



### <a name="SSMLTextBuilder.AppendSubstitution">func</a> (\*SSMLTextBuilder) [AppendSubstitution](/src/target/ssml_builder.go?s=1876:1960#L86)
``` go
func (this *SSMLTextBuilder) AppendSubstitution(text, alias string) *SSMLTextBuilder
```



### <a name="SSMLTextBuilder.ApplyBackground">func</a> (\*SSMLTextBuilder) [ApplyBackground](/src/target/ssml_builder.go?s=2355:2441#L103)
``` go
func (this *SSMLTextBuilder) ApplyBackground(src string, repeat bool) *SSMLTextBuilder
```



### <a name="SSMLTextBuilder.Build">func</a> (\*SSMLTextBuilder) [Build](/src/target/ssml_builder.go?s=2705:2748#L116)
``` go
func (this *SSMLTextBuilder) Build() string
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
