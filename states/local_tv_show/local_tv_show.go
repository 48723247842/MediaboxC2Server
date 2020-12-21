package localtvshow

import (
	logrus "github.com/sirupsen/logrus"
	//redis "github.com/0187773933/RedisManagerUtils/manager"
	utils "c2server/utils"
	substates "c2server/states/local_tv_show/substates"
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


func Start() ( result string ) {
	logger.Info( "State === LocalTVShow === Start()" )
	substates.StartNextShowInCircularListAndNextEpisodeInCircularList()
	return
}