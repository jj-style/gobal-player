module github.com/jj-style/gobal-player/pkg/globalplayer

go 1.22.0

replace github.com/jj-style/gobal-player/pkg/resty => ../resty

require github.com/jj-style/gobal-player/pkg/resty v0.0.0-00010101000000-000000000000

require (
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.18.0 // indirect
)
