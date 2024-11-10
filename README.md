## How to use

tlog is forked from repo https://github.com/trpc-group/trpc-go

```golang

import "github.com/phpgao/tlog"


func main(){
    var c = []tlog.OutputConfig{
        {
            Writer:    "console",
            Level:     "info",
            Formatter: "json",
        },
        {
            Writer: "file",
            WriteConfig: tlog.WriteConfig{
                LogPath:    "/tmp/",
                Filename:   "trpc.log",
                RollType:   "size",
                MaxAge:     1,
                MaxBackups: 5,
                Compress:   false,
                MaxSize:    1024,
            },
            Level:     "debug",
            Formatter: "json",
        },
    }
    
    
    tlog.Infof("123")
    tlog.Register("default", tlog.NewZapLog(c))
    tlog.Infof("123")
    ctx := tlog.WithContextFields(context.TODO(), "key1", "value1", "key2", "value2")
    tlog.InfoContextf(ctx, "123%s", "sss")
    tlog.ErrorContextf(ctx, "123%s", "zzz")


    tlog.SetLevel("0", tlog.LevelInfo)
	
	// use RegisterHandler to set log level
    mux := http.NewServeMux()
    RegisterHandlerWithPath(mux, "/")
    http.ListenAndServe(":8080", mux)
}
```
