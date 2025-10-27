//go:build webui

package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

//go:embed ../../webui/dist
var webuiFiles embed.FS

func registerWebUI() {
}
