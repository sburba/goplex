package goplex

type Video struct {
	AddedAt          UnixTime       `xml:"addedAt,attr"`
	Art              UrlPath        `xml:"art,attr"`
	ContentRating    string         `xml:"contentRating,attr"`
	Duration         MillisDuration `xml:"duration,attr"`
	GrandparentArt   UrlPath        `xml:"grandparentArt,attr"`
	GrandparentTheme UrlPath        `xml:"grandparentTheme,attr"`
	GrandparentThumb UrlPath        `xml:"grandparentThumb,attr"`
	GrandparentTitle string         `xml:"grandparentTitle,attr"`
	GUID             string         `xml:"guid,attr"`
	ParentThumb      UrlPath        `xml:"parentThumb,attr"`
	Thumb            UrlPath        `xml:"thumb,attr"`
	Title            string         `xml:"title,attr"`
	UpdatedAt        UnixTime       `xml:"updatedAt,attr"`
	Media            Media
	User             User
	Player           Player
	TranscodeSession TranscodeSession
}

type Media struct {
	AspectRatio    float32 `xml:"aspectRatio,attr"`
	AudioChannels  int     `xml:"audioChannels,attr"`
	AudioCodec     string  `xml:"audioCodec,attr"`
	VideoCodec     string  `xml:"videoCodec,attr"`
	VideoFrameRate string  `xml:"videoFrameRate,attr"`
	HeightPx       int     `xml:"videoResolution,attr"`
	WidthPx        int     `xml:"width,attr"`
}

type Player struct {
	MachineIdentifier string `xml:"machineIdentifier,attr"`
	Platform          string `xml:"platform,attr"`
	Product           string `xml:"product,attr"`
	State             string `xml:"state,attr"`
	Title             string `xml:"title,attr"`
}

type TranscodeSession struct {
	Key           string         `xml:"key,attr"`
	Throttled     IntAsBool      `xml:"throttled,attr"`
	Progress      float64        `xml:"progress,attr"`
	Speed         float64        `xml:"speed,attr"`
	Duration      MillisDuration `xml:"duration,attr"`
	VideoDecision string         `xml:"videoDecision,attr"`
	AudioDecision string         `xml:"audioDecision,attr"`
	Protocol      string         `xml:"protocol,attr"`
	Container     string         `xml:"container,attr"`
	VideoCodec    string         `xml:"videoCodec,attr"`
	AudioCodec    string         `xml:"audioCodec,attr"`
	AudioChannels int            `xml:"audioChannels,attr"`
	Width         int            `xml:"width,attr"`
	Height        int            `xml:"height,attr"`
}
