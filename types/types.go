package types

type TestStruct struct {
	Wadu string `json:wadu`
	Waduagain []int `json:waduagain`
}

type LoggerMain struct {
	TimeStamp string `json:time_stamp`
	NanosecondsSinceEpoch int64 `json:nanoseconds_since_epoch`
	Msg string `json:msg`
	Author string `json:author`
	Fields map[string]interface{} `json:fields`
	File string `json:file`
	Function string `json:function`
	Line int `json:line`
	Level string `json:level`
}

type NowPlaying struct {
	Title string `json:title`
	Artist string `json:artist`
	LocalFilePath string `json:local_file_path`
	LocalFilePathB64 string `json:local_file_path_b64`
	URL string `json:url`
	Times struct {
		Duration struct {
			Seconds int64 `json:seconds`
			Time int64 `json:seconds`
		} `json:duration`
		CurrentPosition struct {
			Seconds int64 `json:seconds`
			Time int64 `json:seconds`
		} `json:current_position`
		Remaining struct {
			Seconds int64 `json:seconds`
			Time int64 `json:seconds`
		} `json:remaining`
	} `json:times`
	Stats struct {
		Skipped bool `json:skipped`
		NumberOfTimesSkipped bool `json:number_of_times_skipped`
		Watched bool `json:watched`
		NumberOfTimesWatched int `json:number_of_times_watched`
		Completed bool `json:completed`
		NumberOfTimesCompleted int `json:number_of_times_completed`
	} `json:stats`
}

type StateMetaData struct {
	Name string `json:name`
	GenericType string `json:generic_type`
	RestartOnFail bool `json:restart_on_fail`
	NowPlaying NowPlaying `json:now_playing`
}