package cmd

import (
	"fmt"
	"github.com/crazygit/go-downloading-tool/util"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	url, output       string
	quitAfterDownload bool
	urlPromptContent  = util.NewPromptContent(
		"The url to download?",
		"Please provide a valid url. for example: <http://212.183.159.230/5MB.zip>",
		util.OptionSetValidator(util.UrlValidator),
		util.OptionSetDefault("http://212.183.159.230/5MB.zip"),
	)
	outputPromptContent = util.NewPromptContent(
		"Where to save the download?",
		"Please provide a valid path. for example: /path/to/download/filename.ext",
		util.OptionSetDefault("5MB.zip"),
	)
	quitPromptContent = util.NewPromptContent(
		"Quit?",
		"Please choose Y/N",
		util.OptionSetValidator(util.YesNoValidator),
	)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-downloading-tool",
	Short: "Download something from a given url",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd()
	},
	Version: "0.01",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.speedtest.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.MousetrapHelpText = "" // Disable prompt "This is a command line tool" on Windows system
	rootCmd.Flags().StringVarP(&url, "url", "u", "", urlPromptContent.Label)
	rootCmd.Flags().StringVarP(&output, "output", "o", "", outputPromptContent.Label)
	rootCmd.Flags().BoolVarP(&quitAfterDownload, "quit", "q", false, outputPromptContent.Label)
}

func downloadFile(url, output string) (err error) {
	// Create the file
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			util.Log.Error(err)
		}
	}(f)

	//Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			util.Log.Error(err)
		}
	}(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	// Writer the body to file
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return err
	}
	if err != nil {
		log.Fatal(err)
	}
	util.Log.WithField("costTimeInSeconds", bar.State().SecondsSince).Info("Download Success")
	return nil
}

func runCmd() {
START:
	if url == "" {
		url = util.PromptGetInput(urlPromptContent)
	}
	if output == "" {
		output = util.PromptGetInput(outputPromptContent)
	}
	util.Log.WithFields(logrus.Fields{
		"url":    url,
		"output": output,
	}).Info("Download Info")
	if err := downloadFile(url, output); err != nil {
		util.Log.Error(err)
	}
	if !quitAfterDownload {
		quit := util.PromptGetInput(quitPromptContent)
		if strings.ToLower(quit) == "y" {
			os.Exit(0)
		} else {
			url = ""
			output = ""
			goto START
		}
	}
}
