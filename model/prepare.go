package model

import (
	"os"

	"github.com/oswaldoooo/octools/toolsbox"
)

var ROOTPATH = os.Getenv("HTTPSQL_HOME")
var Errorlog = toolsbox.LogInit("error", ROOTPATH+"/log/error.log")
