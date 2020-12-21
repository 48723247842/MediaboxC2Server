package spotify

import (
	"fmt"
	utils "c2server/utils"
	"encoding/json"
	types "c2server/types"
	logrus "github.com/sirupsen/logrus"
	redis "github.com/0187773933/RedisManagerUtils/manager"
	//vlc "github.com/0187773933/VLCWrapper/wrapper"
)

var logger *logrus.Entry = utils.BuildLogger( "LocalTVShow" )

func Stop() ( result string ) {
	logger.Info( "State === LocalTVShow === Stop()" )
	result = "failed"
	// logger.WithFields( logrus.Fields{
	// 	"command": "spotify_status" ,
	// 	"vlc_status": spotify.Status ,
	// }).Info( "State === LocalTVShow === VLC Status" )
	return
}

func Play() ( result string ) {
	logger.Info( "State === LocalTVShow === Play()" )
	result = "failed"
	// logger.WithFields( logrus.Fields{
	// 	"command": "spotify_status" ,
	// 	"spotify_status": spotify.Status ,
	// }).Info( "State === LocalTVShow === VLC Status" )
	return
}

func Pause() ( result string ) {
	logger.Info( "State === LocalTVShow === Pause()" )
	result = "failed"
	// logger.WithFields( logrus.Fields{
	// 	"command": "spotify_status" ,
	// 	"spotify_status": spotify.Status ,
	// }).Info( "State === LocalTVShow === VLC Status" )
	return
}

func Previous() ( result string ) {
	logger.Info( "State === LocalTVShow === Previous()" )
	result = "failed"
	// logger.WithFields( logrus.Fields{
	// 	"command": "spotify_status" ,
	// 	"spotify_status": spotify.Status ,
	// }).Info( "State === LocalTVShow === VLC Status" )
	return
}

func Next() ( result string ) {
	logger.Info( "State === LocalTVShow === Next()" )
	result = "failed"
	// logger.WithFields( logrus.Fields{
	// 	"command": "spotify_status" ,
	// 	"spotify_status": spotify.Status ,
	// }).Info( "State === LocalTVShow === VLC Status" )
	return
}

func get_current_episode( redis *redis.Manager ) {
	current_tv_show_name_b64 , _ := redis.CircleCurrent( "MEDIA_MANAGER.TVShows.LIST" )
	current_tv_show_name := utils.Base64Decode( current_tv_show_name_b64 )
	fmt.Println( current_tv_show_name )
}

func StartNextShowInCircularListAndNextEpisodeInCircularList() ( result string ) {
	logger.Info( "State === LocalTVShow === StartNextShowInCircularListAndNextEpisodeInCircularList()" )
	result = "failed"
	redis := redis.Manager{}
	redis.Connect( "localhost:6379" , 3 , "" )
	//next_playlist_uri := redis.CircleNext( "CONFIG.SPOTIFY.PLAYLISTS.GENERES.MISC" )
	get_current_episode( &redis )

	// logger.WithFields( logrus.Fields{
	// 	"command": "next_playlist_uri" ,
	// 	"next_playlist_uri": next_playlist_uri ,
	// }).Info("State === LocalTVShow === StartNextShowInCircularListAndNextEpisodeInCircularList() === Next TV Show Episode")
	// spotify := spotify_dbus.Controller{}
	// spotify.Connect()
	// spotify.OpenURI( next_playlist_uri )
	// logger.WithFields( logrus.Fields{
	// 	"command": "spotify_status" ,
	// 	"spotify_status": spotify.Status ,
	// }).Info("State === LocalTVShow === StartNextShowInCircularListAndNextEpisodeInCircularList() === VLC Status")
	return
}

func build_state_meta_data( state_name string ) ( json_string string ) {
	state_data := types.StateMetaData {
		Name: state_name ,
		GenericType: "LocalTVShow" ,
		RestartOnFail: true ,
		NowPlaying: types.NowPlaying{} ,
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
	state_meta_data := build_state_meta_data( state_name )
	logger.WithFields( logrus.Fields{
		"command": "new_state" ,
		"new_state": state_meta_data ,
	}).Info( "State === LocalTVShow === NEW STATE" )
	redis.Set( "STATE.CURRENT" , state_meta_data )
}

func Start() ( result string ) {
	logger.Info( "State === LocalTVShow === Start()" )
	swap_current_and_previous_state_info( "LocalTVShowNextShowInCircularListAndNextEpisodeInCircularList" )
	result = StartNextShowInCircularListAndNextEpisodeInCircularList()
	return
}