package utils

import (
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"fmt"
	"strings"
	types "c2server/types"
	redis "github.com/0187773933/RedisManagerUtils/manager"
)

// Using Webdis Redis HTTP Server
// https://github.com/nicolasff/webdis
// sudo docker run -dit --restart='always' \
// --name webdis --network host \
// nicolas/webdis

// And then this script as a makeshift "log viewer"
// We can then watch all publishes to LOG.ALL in realtime on webbrowser
// <script>
// var previous_response_length = 0
// xhr = new XMLHttpRequest()
// xhr.open("GET", "http://100.112.11.51:7379/3/SUBSCRIBE/LOG.ALL", true);
// xhr.onreadystatechange = checkData;
// xhr.send(null);

// function checkData() {
// 	if(xhr.readyState == 3)  {
// 		response = xhr.responseText;
// 		chunk = response.slice(previous_response_length);
// 		previous_response_length = response.length;
// 		try{
// 			const message = chunk.split( '{"SUBSCRIBE":["message","LOG.ALL",' )[1].split( "]}" )[0]
// 			const parsed = JSON.parse( JSON.parse( message ) );
// 			console.log( `${parsed.TimeStamp} === ${parsed.Msg}` );
// 			console.log( parsed );
// 		}catch(e){}
// 	}
// };
// </script>
func AddLogToRedis( input_struct *types.LoggerMain ) {
	redis := redis.Manager{}
	redis.Connect( "localhost:6379" , 3 , "" )
	json_marshal_result , json_marshal_error := json.Marshal( input_struct )
	if json_marshal_error != nil { panic( json_marshal_error ) }
	json_string := string( json_marshal_result )
	fmt.Println( json_string )
	redis.ListPushRight( "LOG.ALL" , json_string )
	//var ctx = context.Background()
	//redis.Redis.Do( ctx , "PUBLISH" , "LOG.ALL" , json_string )
	redis.Publish( "LOG.ALL" , json_string )
	return
}

// func GetJSONStringFromRedis( redis_key string ) {
// 	json_get_test := redis.Get( "testmeta" )
// 	var json_get_test_struct TestStruct
// 	json_unmarshal_error := json.Unmarshal( []byte( json_get_test ) , &json_get_test_struct )
// 	if json_unmarshal_error != nil { panic( json_unmarshal_error ) }
// 	fmt.Println( json_get_test_struct )
// }

// https://godoc.org/github.com/sirupsen/logrus#Entry
// https://stackoverflow.com/a/54314594
type LoggerMain struct {}
type LoggerMainHook struct {}
func ( hook *LoggerMainHook ) Fire( entry *log.Entry ) error {
	time_stamp := fmt.Sprintf( "%d%s%d===%02d:%02d:%02d" ,
		entry.Time.Day() , strings.ToUpper( entry.Time.Month().String()[:3] ) , entry.Time.Year() ,
		entry.Time.Hour() , entry.Time.Minute() , entry.Time.Second() ,
	)
	new_log_line := types.LoggerMain{
		TimeStamp:  time_stamp ,
		NanosecondsSinceEpoch: entry.Time.UnixNano() ,
		Msg: entry.Message ,
		Author: entry.Data["author"].(string) ,
		Fields: entry.Data ,
		File: entry.Caller.File ,
		Function: entry.Caller.Function ,
		Line: entry.Caller.Line ,
		Level: entry.Level.String() ,
	}
	AddLogToRedis( &new_log_line )
	return nil
}
func ( hook *LoggerMainHook ) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
		log.DebugLevel,
	}
}
// func ( hook *CallerHook ) caller( skip int ) string {
// 	if _, file, line, ok := runtime.Caller(skip); ok {
// 		return strings.Join([]string{filepath.Base(file), strconv.Itoa(line)}, ":")
// 	}
// 	// not sure what the convention should be here
// 	return "???"
// }

func BuildLogger( author_name string ) ( logger *log.Entry ) {

	log.SetFormatter( &log.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true ,
	})
	log.SetFormatter( &log.JSONFormatter{ DisableHTMLEscape: true } )
	log.SetReportCaller( true )

	logger_main_hook := LoggerMainHook{}
	log.AddHook( &logger_main_hook )

	logger = log.WithFields( log.Fields{
		"author": author_name ,
	})

	return
}