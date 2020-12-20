package main

import (
	//"time"
	// "fmt"
	utils "c2server/utils"
	//types "c2server/types"
	spotify "c2server/states/spotify"
)


func main() {
	logger := utils.BuildLogger( "C2Server" )
	logger.Info( "main()" )
	// fmt.Println( result )
	// spotify.StartNextInCircularListOfMiscGenrePlaylists()
	spotify.Start()

	// test_data := types.TestStruct{
	// 	Wadu: "wadu wadu wer234" ,
	// 	Waduagain: []int{1,2,3,4} ,
	// }
	// json_string := utils.ConvertCommonTypeStructToJSONString( &test_data )
	// fmt.Println( json_string )
}