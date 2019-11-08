## MJPEG Streaming Library - And Proxy

I borrowed this code from `github.com/saljam/mjpeg` and made a few tweaks and fixed a few issues.
It needs some optimzing still

## Proxy - Config
The proxy is in the mjproxy directory. 

The config is very simple, by default the docker image looks for a config file `/opt/rts2p/rts2p.yaml` but the
config library supports yaml, toml or json. You can specify `-c` to the config file you want. 

Every config option is shown:
```
server:
  port: 8080
  username: myusername
  password: mypassword

streams:
  - url: "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mov"
    name: "mytestfeed"
```

This would create the mjpeg feed `http://server:8080/mytestfeed` and require a login of `myusername` with password `mypassword`
Omit the username and password fields if you do not want to require a login.

## Docker

Example docker command: `docker run -it -p 5554:5554 -v /my/config.yaml:/opt/mjproxy/mjproxy.yaml snowzach/mjproxy:latest`
