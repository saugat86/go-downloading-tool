package main

import (
	"os"

	"github.com/saugat86/go-downloading-tool/cmd"
	"github.com/saugat86/go-downloading-tool/util"
)

func main() {
	util.Log.WithField("args", os.Args).Info("App Started")
	cmd.Execute()
}
