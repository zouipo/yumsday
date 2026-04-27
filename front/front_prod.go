//go:build !dev

package front

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var frontFS embed.FS

func Handler() http.HandlerFunc {
	sub, err := fs.Sub(frontFS, "dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServerFS(sub)

	return func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	}
}
