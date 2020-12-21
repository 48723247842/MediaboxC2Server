package substates

import (
	"fmt"
	utils "c2server/utils"
	"encoding/json"
	"strings"
	"context"
	"os"
	"time"
	types "c2server/types"
	logrus "github.com/sirupsen/logrus"
	redis "github.com/0187773933/RedisManagerUtils/manager"
	ffprobe "gopkg.in/vansante/go-ffprobe.v2"
	vlc "github.com/0187773933/VLCWrapper/wrapper"
)

var logger *logrus.Entry = utils.BuildLogger( "LocalTVShow" )

func get_title_from_episode_path( episode_path string ) ( title string ){
	parts := strings.Split( episode_path , "/" )
	name := parts[len(parts)-1]
	name_parts := strings.Split( name , ".mp4" )
	title = name_parts[0]
	return
}
func ffprobe_local_file_for_duration_seconds( episode_path string ) ( duration_seconds int64 ) {
	// https://pkg.go.dev/gopkg.in/vansante/go-ffprobe.v2#Format
	duration_seconds = 0
	ctx , cancelFn := context.WithTimeout( context.Background() , 30*time.Second )
	defer cancelFn()
	fileReader , err := os.Open( episode_path )
	if err != nil {
		logger.Info( "Couldn't Find Video Information via FFProbe" )
		return
	}
	data , err := ffprobe.ProbeReader( ctx , fileReader )
	if err != nil {
		logger.Info( "Couldn't Find Video Information via FFProbe" )
		return
	}
	duration_seconds = int64( data.Format.DurationSeconds )
	return
}

func build_default_episode_meta( current_tv_show_name_b64 string , current_tv_show_name string , current_tv_show_index string ,
	current_episode_name_b64 string , current_episode_name string , current_tv_show_episode_index string ) ( current_episode_meta_struct types.NowPlayingMeta ) {

		episode_duration_seconds := ffprobe_local_file_for_duration_seconds( current_episode_name )

		current_episode_meta_struct = types.NowPlayingMeta {
		Title: get_title_from_episode_path( current_episode_name ) ,
		Artist: current_tv_show_name ,
		LocalFilePath: current_episode_name ,
		LocalFilePathB64: current_episode_name_b64 ,
		ShowIndex: current_tv_show_index ,
		EpisodeIndex: current_tv_show_episode_index ,
		Times: types.TimesObject {
			Duration: types.TimeObject {
				Seconds: episode_duration_seconds ,
				TimeStamp: "00:00:00" ,
			} ,
			CurrentPosition: types.TimeObject {
				Seconds: 0 ,
				TimeStamp: "00:00:00" ,
			} ,
			Remaining: types.TimeObject {
				Seconds: 0 ,
				TimeStamp: "00:00:00" ,
			} ,
		} ,
		Stats: types.StatsObject {
			Skipped: false ,
			NumberOfTimesSkipped: 0	,
			Watched: false ,
			NumberOfTimesWatched: 0 ,
			Completed: false ,
			NumberOfTimesCompleted: 0 ,
		} ,
	}
	return
}

func get_current_episode( redis *redis.Manager ) ( current_episode_meta_struct types.NowPlayingMeta ) {
	current_tv_show_name_b64 , current_tv_show_index := redis.CircleCurrent( "MEDIA_MANAGER.TVShows.LIST" )
	current_tv_show_name := utils.Base64Decode( current_tv_show_name_b64 )
	current_episode_name_b64 , current_tv_show_episode_index := redis.CircleCurrent( fmt.Sprintf( "MEDIA_MANAGER.TVShows.%s" , current_tv_show_name_b64 ) )
	current_episode_name := utils.Base64Decode( current_episode_name_b64 )
	current_episode_meta_json_string := redis.Get( fmt.Sprintf( "MEDIA_MANAGER.TVShows.META.%s" , current_episode_name_b64 ) )
	if current_episode_meta_json_string == "failed" {
		current_episode_meta_struct = build_default_episode_meta( current_tv_show_name_b64 , current_tv_show_name , current_tv_show_index ,
	current_episode_name_b64 , current_episode_name , current_tv_show_episode_index )
		json_marshal_result , json_marshal_error := json.Marshal( &current_episode_meta_struct )
		if json_marshal_error != nil { panic( json_marshal_error ) }
		json_string := string( json_marshal_result )
		redis.Set( fmt.Sprintf( "MEDIA_MANAGER.TVShows.META.%s" , current_episode_name_b64 ) , json_string )
	} else {
		json_unmarshal_error := json.Unmarshal( []byte( current_episode_meta_json_string ) , &current_episode_meta_struct )
		if json_unmarshal_error != nil { panic( json_unmarshal_error ) }
	}
	return
}

func update_time_info( episode *types.NowPlayingMeta ) {}

func episode_is_over( episode *types.NowPlayingMeta ) ( result bool ) {
	result = false
	if episode.Times.Duration.Seconds > 0 {
		if episode.Times.CurrentPosition.Seconds >= ( episode.Times.Duration.Seconds - 20 ) {
			result = true
		}
	}
	return
}

func build_state_meta_data( state_name string ) ( state_data types.StateMetaData ) {
	state_data = types.StateMetaData {
		Name: state_name ,
		GenericType: "LocalTVShow" ,
		RestartOnFail: true ,
		NowPlaying: types.NowPlayingMeta{} ,
	}
	return
}

func build_state_meta_data_json( state_name string ) ( json_string string ) {
	state_data := types.StateMetaData {
		Name: state_name ,
		GenericType: "LocalTVShow" ,
		RestartOnFail: true ,
		NowPlaying: types.NowPlayingMeta{} ,
	}
	json_marshal_result , json_marshal_error := json.Marshal( state_data )
	if json_marshal_error != nil { panic( json_marshal_error ) }
	json_string = string( json_marshal_result )
	return
}

func swap_current_and_previous_state_info( state_name string ) {
	redis := redis.Manager{}
	redis.Connect( "localhost:6379" , 3 , "" )
	state_current := redis.Get( "STATE.CURRENT" )
	logger.WithFields( logrus.Fields{
		"command": "state_current" ,
		"state_current": state_current ,
	}).Info( "State === LocalTVShow === STATE CURRENT" )
	redis.Set( "STATE.PREVIOUS" , state_current )
}



func StartNextShowInCircularListAndNextEpisodeInCircularList() ( result string ) {
	logger.Info( "State === LocalTVShow === StartNextShowInCircularListAndNextEpisodeInCircularList()" )
	swap_current_and_previous_state_info( "LocalTVShowNextShowInCircularListAndNextEpisodeInCircularList" )
	result = "failed"
	redis := redis.Manager{}
	redis.Connect( "localhost:6379" , 3 , "" )
	current_episode := get_current_episode( &redis )

	// 1.) Check If Current Episode is Over
	if episode_is_over( &current_episode ) {
		redis.CircleNext( "MEDIA_MANAGER.TVShows.LIST" )
		current_episode = get_current_episode( &redis )
	}

	// 2.) Check If Current Episode's CurrentPosition is > 0
	if current_episode.Times.CurrentPosition.Seconds > 0 {
		// We Need To Seek Into current_episode.CurrentPosition.Seconds
	}

	logger.WithFields( logrus.Fields{
		"command": "local_tv_show_current_episode" ,
		"local_tv_show_current_episode": current_episode ,
	}).Info( "State === LocalTVShow === Current Episode" )

	state_meta_data := build_state_meta_data( "LocalTVShowNextShowInCircularListAndNextEpisodeInCircularList" )
	state_meta_data.NowPlaying = current_episode
	logger.WithFields( logrus.Fields{
		"command": "new_state" ,
		"new_state": state_meta_data ,
	}).Info( "State === LocalTVShow === NEW STATE" )
	json_marshal_result , json_marshal_error := json.Marshal( state_meta_data )
	if json_marshal_error != nil { panic( json_marshal_error ) }
	json_string := string( json_marshal_result )
	redis.Set( "STATE.CURRENT" , json_string )

	p := vlc.NewPlayer( nil )
	p.Start()
	p.Play( current_episode.LocalFilePath )
	time.Sleep( 3 * time.Second )
	p.Fullscreen()

	return
}