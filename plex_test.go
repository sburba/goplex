package goplex

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestSignInSuccess(t *testing.T) {
	resp := `<?xml version="1.0" encoding="UTF-8"?>
	<user email="email@address.com" id="123456" thumb="http://thumb.com" username="username" title="title" cloudSyncDevice="" locale="locale" authenticationToken="authtoken" restricted="0" home="0" queueEmail="queue@email.com" queueUid="queueId" maxHomeSize="15">
	  <subscription active="1" status="Active" plan="lifetime">
	    <feature id="pass"/>
	    <feature id="sync"/>
	    <feature id="cloudsync"/>
	    <feature id="home"/>
	  </subscription>
	  <roles>
	    <role id="plexpass"/>
	  </roles>
	  <entitlements all="1">
	    <entitlement id="roku"/>
	    <entitlement id="android"/>
	    <entitlement id="xbox_one"/>
	    <entitlement id="xbox_360"/>
	    <entitlement id="windows"/>
	    <entitlement id="windows_phone"/>
	  </entitlements>
	  <username>username</username>
	  <email>email@address.com</email>
	  <joined-at type="datetime">2013-03-26 00:45:24 UTC</joined-at>
	  <authentication-token>authtoken</authentication-token>
	</user>`

	expectedReq, err := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", nil)
	if err != nil {
		t.Fatal(err)
	}
	expectedReq.SetBasicAuth("username", "password")
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")

	client = makeFakeClient(t, http.StatusCreated, resp, expectedReq)

	user, err := GetUser("username", "password")
	if err != nil {
		t.Fatal(err)
	}

	expected := User{
		Email:      "email@address.com",
		Id:         123456,
		Thumb:      HttpUrl{url.URL{Scheme: "http", Host: "thumb.com"}},
		Username:   "username",
		Title:      "title",
		Locale:     "locale",
		AuthToken:  "authtoken",
		QueueEmail: "queue@email.com",
		Subscription: Subscription{
			Active: true,
			Plan:   "lifetime",
		},
	}

	if user != expected {
		t.Fatalf("\nExpected: %+v\nGot: %+v", expected, user)
	}
}

func TestSignInFail(t *testing.T) {
	expectedReq, err := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", nil)
	if err != nil {
		t.Fatal(err)
	}
	expectedReq.SetBasicAuth("username", "password")
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")

	client = makeFakeClient(t, http.StatusUnauthorized, "", expectedReq)

	_, err = GetUser("username", "password")
	if err == nil {
		t.Fatal("Should return error when username and password are wrong")
	}
}

func TestGetDevicesSuccess(t *testing.T) {
	resp := `<?xml version="1.0" encoding="UTF-8"?>
	<MediaContainer publicAddress="serverPublicAddress.com">
	  <Device name="My Nexus 5" publicAddress="24.56.78.91" product="Plex for Android" productVersion="4.2.3.358" platform="Android" platformVersion="5.1" device="Nexus 5" model="hammerhead" vendor="LGE" provides="controller,sync-target" clientIdentifier="caac4066dbaa6a9c-com-plexapp-android" version="4.2.3.358" id="12345678" token="token" createdAt="1422553670" lastSeenAt="1430577652" screenResolution="1920x1080" screenDensity="480">
	    <SyncList itemsCompleteCount="0" totalSize="0" version="1"/>
	    <Connection uri="http://192.168.1.1:32400"/>
	  </Device>
	  <Device name="Server" publicAddress="serverPublicAddress.com" product="Plex Media Server" productVersion="0.9.11.17.986-269b82b" platform="Linux" platformVersion="3.13.0-45-generic (#74-Ubuntu SMP Tue Jan 13 19:36:28 UTC 2015)" device="PC" model="x86_64" vendor="ubuntu" provides="server" clientIdentifier="clientIdentifier" version="0.9.11.17.986-269b82b" id="12345679" token="token" createdAt="1394924489" lastSeenAt="1430601269" screenResolution="" screenDensity="">
	    <Connection uri="http://serverPublicAddress.com:12345"/>
	    <Connection uri="http://192.168.1.2:32400"/>
	    <Connection uri="http://192.168.1.2:32400"/>
	  </Device>
	</MediaContainer>`

	user := User{
		AuthToken: "authToken",
	}

	expectedReq, err := http.NewRequest("GET", "https://plex.tv/devices.xml", nil)
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")
	expectedReq.Header.Add("X-Plex-Token", "authToken")
	client = makeFakeClient(t, 200, resp, expectedReq)

	result, err := user.GetDevices()
	if err != nil {
		t.Fatal(err)
	}

	expected := []Device{
		Device{
			Name:          "My Nexus 5",
			PublicAddress: HttpUrl{url.URL{Scheme: "http", Host: "24.56.78.91"}},
			Product:       "Plex for Android",
			Provides:      []string{"controller", "sync-target"},
			Connections:   []Connection{Connection{HttpUrl{url.URL{Scheme: "http", Host: "192.168.1.1:32400"}}}},
			Owner:         user,
		},
		Device{
			Name:          "Server",
			PublicAddress: HttpUrl{url.URL{Scheme: "http", Host: "serverPublicAddress.com"}},
			Product:       "Plex Media Server",
			Provides:      []string{"server"},
			Connections: []Connection{
				Connection{HttpUrl{url.URL{Scheme: "http", Host: "serverPublicAddress.com:12345"}}},
				Connection{HttpUrl{url.URL{Scheme: "http", Host: "192.168.1.2:32400"}}},
				Connection{HttpUrl{url.URL{Scheme: "http", Host: "192.168.1.2:32400"}}},
			},
			Owner: user,
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nExpected: %+v\n\nGot: %+v", expected, result)
	}
}

func TestGetDevicesFail(t *testing.T) {
	user := User{
		AuthToken: "authToken",
	}

	expectedReq, err := http.NewRequest("GET", "https://plex.tv/devices.xml", nil)
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")
	expectedReq.Header.Add("X-Plex-Token", "authToken")

	client = makeFakeClient(t, http.StatusUnauthorized, "", expectedReq)

	_, err = user.GetDevices()
	if err == nil {
		t.Fatal("Should err when server returns 401")
	}
}

func TestGetServersSuccess(t *testing.T) {
	resp := `<?xml version="1.0" encoding="UTF-8"?>
	<MediaContainer publicAddress="serverPublicAddress.com">
	  <Device name="My Nexus 5" publicAddress="24.56.78.91" product="Plex for Android" productVersion="4.2.3.358" platform="Android" platformVersion="5.1" device="Nexus 5" model="hammerhead" vendor="LGE" provides="controller,sync-target" clientIdentifier="caac4066dbaa6a9c-com-plexapp-android" version="4.2.3.358" id="12345678" token="token" createdAt="1422553670" lastSeenAt="1430577652" screenResolution="1920x1080" screenDensity="480">
	    <SyncList itemsCompleteCount="0" totalSize="0" version="1"/>
	    <Connection uri="http://192.168.1.1:32400"/>
	  </Device>
	  <Device name="Server" publicAddress="serverPublicAddress.com" product="Plex Media Server" productVersion="0.9.11.17.986-269b82b" platform="Linux" platformVersion="3.13.0-45-generic (#74-Ubuntu SMP Tue Jan 13 19:36:28 UTC 2015)" device="PC" model="x86_64" vendor="ubuntu" provides="server" clientIdentifier="clientIdentifier" version="0.9.11.17.986-269b82b" id="12345679" token="token" createdAt="1394924489" lastSeenAt="1430601269" screenResolution="" screenDensity="">
	    <Connection uri="http://serverPublicAddress.com:12345"/>
	    <Connection uri="http://192.168.1.2:32400"/>
	    <Connection uri="http://192.168.1.2:32400"/>
	  </Device>
	</MediaContainer>`

	user := User{
		AuthToken: "authToken",
	}

	expectedReq, err := http.NewRequest("GET", "https://plex.tv/devices.xml", nil)
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")
	expectedReq.Header.Add("X-Plex-Token", "authToken")
	client = makeFakeClient(t, 200, resp, expectedReq)

	result, err := user.GetServers()
	if err != nil {
		t.Fatal(err)
	}

	expected := []Server{
		Server{
			Device{
				Name:          "Server",
				PublicAddress: HttpUrl{url.URL{Scheme: "http", Host: "serverPublicAddress.com:12345"}},
				Product:       "Plex Media Server",
				Provides:      []string{"server"},
				Connections: []Connection{
					Connection{HttpUrl{url.URL{Scheme: "http", Host: "serverPublicAddress.com:12345"}}},
					Connection{HttpUrl{url.URL{Scheme: "http", Host: "192.168.1.2:32400"}}},
					Connection{HttpUrl{url.URL{Scheme: "http", Host: "192.168.1.2:32400"}}},
				},
				Owner: user,
			},
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nExpected: %+v\n\nGot: %+v", expected, result)
	}
}

func TestGetServersFail(t *testing.T) {
	user := User{
		AuthToken: "authToken",
	}

	expectedReq, err := http.NewRequest("GET", "https://plex.tv/devices.xml", nil)
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")
	expectedReq.Header.Add("X-Plex-Token", "authToken")

	client = makeFakeClient(t, http.StatusUnauthorized, "", expectedReq)

	_, err = user.GetDevices()
	if err == nil {
		t.Fatal("Should err when server returns 401")
	}
}

func TestGetActivitySuccess(t *testing.T) {
	resp := `<?xml version="1.0" encoding="UTF-8"?>
		<MediaContainer size="1">
		<Video addedAt="1430373171" art="/library/metadata/181/art/1430373196" chapterSource="chapterSource" contentRating="TV-PG" duration="1297172" grandparentArt="/library/metadata/181/art/1430373196" grandparentKey="/library/metadata/181" grandparentRatingKey="181" grandparentTheme="/library/metadata/181/theme/1430373196" grandparentThumb="/library/metadata/181/thumb/1430373196" grandparentTitle="Modern Family" guid="com.plexapp.agents.thetvdb://95011/6/21?lang=en" index="21" key="/library/metadata/1751" librarySectionID="1" parentIndex="6" parentKey="/library/metadata/1117" parentRatingKey="1117" parentThumb="/library/metadata/1117/thumb/1430373196" ratingKey="1751" sessionKey="11" summary="" thumb="/library/metadata/1751/thumb/1430373196" title="Episode 21" type="episode" updatedAt="1430373196">
		<Media aspectRatio="1.78" audioChannels="6" audioCodec="ac3" bitrate="3874" container="mkv" duration="1297172" height="720" id="1950" videoCodec="h264" videoFrameRate="24p" videoResolution="720" width="1280">
		<Part container="mkv" duration="1297172" file="/media/Media/TV/Modern Family/Season 6/Modern Family - S06E21 - Integrity.mkv" id="2147" indexes="sd" key="/library/parts/2147/file.mkv" size="628172169">
		<Stream bitDepth="8" bitrate="3413" cabac="1" chromaSubsampling="4:2:0" codec="h264" codecID="V_MPEG4/ISO/AVC" colorSpace="yuv" duration="1297172" frameRate="23.976" frameRateMode="cfr" hasScalingMatrix="0" height="720" id="10812" index="0" language="English" languageCode="eng" level="41" profile="high" refFrames="8" scanType="progressive" streamType="1" width="1280" />
		<Stream audioChannelLayout="5.1(side)" bitDepth="16" bitrate="384" bitrateMode="cbr" channels="6" codec="ac3" codecID="A_AC3" dialogNorm="-31" duration="1297152" id="10813" index="1" samplingRate="48000" selected="1" streamType="2" />
		</Part>
		</Media>
		<User id="1" thumb="http://www.thumb.com" title="title" />
		<Player machineIdentifier="5418fbf4404066f0-com-plexapp-android" platform="Android" product="Plex for Android" state="playing" title="My Nexus 7" />
		<TranscodeSession key="5418fbf4404066f0-com-plexapp-android" throttled="1" progress="2.0999999046325684" speed="2.0999999046325684" duration="1297000" videoDecision="transcode" audioDecision="transcode" protocol="hls" container="mpegts" videoCodec="h264" audioCodec="aac" audioChannels="2" width="1280" height="720" />
		</Video>
		</MediaContainer>`

	expectedReq, err := http.NewRequest("GET", "http://server.com:4040/status/sessions", nil)
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")
	expectedReq.Header.Add("X-Plex-Token", "authToken")

	client = makeFakeClient(t, http.StatusOK, resp, expectedReq)

	server := Server{
		Device{
			PublicAddress: HttpUrl{url.URL{Scheme: "http", Host: "server.com:4040"}},
			Owner:         User{AuthToken: "authToken"},
		},
	}

	result, err := server.GetActivity()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expected := []Video{
		Video{
			AddedAt:          UnixTime{time.Unix(1430373171, 0)},
			Art:              UrlPath{url.URL{Path: "/library/metadata/181/art/1430373196"}},
			ContentRating:    "TV-PG",
			Duration:         MillisDuration{time.Duration(1297172) * time.Millisecond},
			GrandparentArt:   UrlPath{url.URL{Path: "/library/metadata/181/art/1430373196"}},
			GrandparentTheme: UrlPath{url.URL{Path: "/library/metadata/181/theme/1430373196"}},
			GrandparentThumb: UrlPath{url.URL{Path: "/library/metadata/181/thumb/1430373196"}},
			GrandparentTitle: "Modern Family",
			GUID:             "com.plexapp.agents.thetvdb://95011/6/21?lang=en",
			ParentThumb:      UrlPath{url.URL{Path: "/library/metadata/1117/thumb/1430373196"}},
			Thumb:            UrlPath{url.URL{Path: "/library/metadata/1751/thumb/1430373196"}},
			Title:            "Episode 21",
			UpdatedAt:        UnixTime{time.Unix(1430373196, 0)},
			Media: Media{
				AspectRatio:    1.78,
				AudioChannels:  6,
				AudioCodec:     "ac3",
				VideoCodec:     "h264",
				VideoFrameRate: "24p",
				HeightPx:       720,
				WidthPx:        1280,
			},
			User: User{
				Id:    1,
				Thumb: HttpUrl{url.URL{Scheme: "http", Host: "www.thumb.com"}},
				Title: "title",
			},
			Player: Player{
				MachineIdentifier: "5418fbf4404066f0-com-plexapp-android",
				Platform:          "Android",
				Product:           "Plex for Android",
				State:             "playing",
				Title:             "My Nexus 7",
			},
			TranscodeSession: TranscodeSession{
				Key:           "5418fbf4404066f0-com-plexapp-android",
				Throttled:     true,
				Progress:      2.0999999046325684,
				Speed:         2.0999999046325684,
				Duration:      MillisDuration{time.Duration(1297000) * time.Millisecond},
				VideoDecision: "transcode",
				AudioDecision: "transcode",
				Protocol:      "hls",
				Container:     "mpegts",
				VideoCodec:    "h264",
				AudioCodec:    "aac",
				AudioChannels: 2,
				Width:         1280,
				Height:        720,
			},
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nExpected: %+v\n\nGot: %+v", expected, result)
	}
}

func TestGetActivityFail(t *testing.T) {
	expectedReq, err := http.NewRequest("GET", "http://server.com:4040/status/sessions", nil)
	expectedReq.Header.Add("X-Plex-Client-Identifier", "plextrack")
	expectedReq.Header.Add("X-Plex-Token", "authToken")
	client = makeFakeClient(t, http.StatusUnauthorized, "", expectedReq)

	server := Server{
		Device{
			PublicAddress: HttpUrl{url.URL{Scheme: "http", Host: "server.com:4040"}},
			Owner:         User{AuthToken: "authToken"},
		},
	}

	if _, err = server.GetActivity(); err == nil {
		t.Fatal("GetActivity returned success when it received bad status code")
	}
}
