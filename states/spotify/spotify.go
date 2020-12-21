package spotify

import (
	"fmt"
	utils "c2server/utils"
	"encoding/json"
	//"reflect"
	types "c2server/types"
	redis "github.com/0187773933/RedisManagerUtils/manager"
	spotify_dbus "github.com/0187773933/SpotifyDBUSController/controller"
	logrus "github.com/sirupsen/logrus"
)

var logger *logrus.Entry = utils.BuildLogger( "Spotify" )

func Stop() ( result string ) {
	logger.Info( "State === Spotify === Stop()" )
	result = "failed"
	spotify := spotify_dbus.Controller{}
	spotify.Connect()
	spotify.Stop()
	fmt.Println( spotify.Status )
	logger.WithFields( logrus.Fields{
		"command": "spotify_status" ,
		"spotify_status": spotify.Status ,
	}).Info( "State === Spotify === Stop() === Spotify Status")
	result = spotify.PlaybackStatus()
	return
}

func Play() ( result string ) {
	logger.Info( "State === Spotify === Play()" )
	result = "failed"
	spotify := spotify_dbus.Controller{}
	spotify.Connect()
	spotify.Play()
	logger.WithFields( logrus.Fields{
		"command": "spotify_status" ,
		"spotify_status": spotify.Status ,
	}).Info( "State == Spotify === Play() === Spotify Status")
	result = spotify.PlaybackStatus()
	return
}

func Pause() ( result string ) {
	logger.Info( "State === Spotify === Pause()" )
	result = "failed"
	spotify := spotify_dbus.Controller{}
	spotify.Connect()
	spotify.Play()
	logger.WithFields( logrus.Fields{
		"command": "spotify_status" ,
		"spotify_status": spotify.Status ,
	}).Info( "State === Spotify === StartNextInCircularListOfMiscGenrePlaylists() === Spotify Status")
	result = spotify.PlaybackStatus()
	return
}

func Previous() ( result string ) {
	logger.Info( "State === Spotify === Previous()" )
	result = "failed"
	spotify := spotify_dbus.Controller{}
	spotify.Connect()
	spotify.Previous()
	logger.WithFields( logrus.Fields{
		"command": "spotify_status" ,
		"spotify_status": spotify.Status ,
	}).Info( "State === Spotify === Previous() === Spotify Status")
	result = spotify.PlaybackStatus()
	return
}

func Next() ( result string ) {
	logger.Info( "State === Spotify === Next()" )
	result = "failed"
	spotify := spotify_dbus.Controller{}
	spotify.Connect()
	spotify.Next()
	logger.WithFields( logrus.Fields{
		"command": "spotify_status" ,
		"spotify_status": spotify.Status ,
	}).Info( "State === Spotify === Next() === Spotify Status" )
	result = spotify.PlaybackStatus()
	return
}

func StartNextInCircularListOfMiscGenrePlaylists() ( result string ) {
	logger.Info( "State === Spotify === StartNextInCircularListOfMiscGenrePlaylists()" )
	result = "failed"
	redis := redis.Manager{}
	redis.Connect( "localhost:6379" , 3 , "" )
	next_playlist_uri := redis.CircleNext( "CONFIG.SPOTIFY.PLAYLISTS.GENERES.MISC" )
	logger.WithFields( logrus.Fields{
		"command": "next_playlist_uri" ,
		"next_playlist_uri": next_playlist_uri ,
	}).Info( "State === Spotify === StartNextInCircularListOfMiscGenrePlaylists() === Next Playlist URI" )
	spotify := spotify_dbus.Controller{}
	spotify.Connect()
	spotify.OpenURI( next_playlist_uri )
	logger.WithFields( logrus.Fields{
		"command": "spotify_status" ,
		"spotify_status": spotify.Status ,
	}).Info( "State === Spotify === StartNextInCircularListOfMiscGenrePlaylists() === Spotify Status" )
	return
}

func build_state_meta_data( state_name string ) ( json_string string ) {
	state_data := types.StateMetaData {
		Name: state_name ,
		GenericType: "Spotify" ,
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
	}).Info( "State === Spotify === swap_current_and_previous_state_info() === STATE CURRENT" )
	redis.Set( "STATE.PREVIOUS" , state_current )
	state_meta_data := build_state_meta_data( state_name )
	logger.WithFields( logrus.Fields{
		"command": "new_state" ,
		"new_state": state_meta_data ,
	}).Info( "State === Spotify === swap_current_and_previous_state_info() === NEW STATE" )
	redis.Set( "STATE.CURRENT" , state_meta_data )
}

func Start() ( result string ) {
	logger.Info( "State === Spotify === Start()" )
	swap_current_and_previous_state_info( "SpotifyStartNextInCircularListOfMiscGenrePlaylists" )
	result = StartNextInCircularListOfMiscGenrePlaylists()
	return
}