package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	"github.com/gen0cide/gscript"
	"github.com/google/go-github/github"
	update "github.com/inconshreveable/go-update"
	"github.com/urfave/cli"
)

func UpdateCLI(c *cli.Context) error {
	gse := gscript.New("UPDATER")
	gse.EnableLogging()
	gse.Logger.DisabledInfo = true

	ctx := context.Background()

	client := github.NewClient(nil)

	repRel, _, err := client.Repositories.GetLatestRelease(ctx, "gen0cide", "gscript")
	if err != nil {
		gse.LogCritf("Github Error: " + err.Error())
	}

	tagName := strings.TrimLeft(repRel.GetTagName(), "v")

	if c.App.Version == tagName {
		gse.LogInfof("Running Lastest Version (no update available): %s", c.App.Version)
		return nil
	}

	gse.LogInfof("New Version Found: %s", tagName)

	fileName := fmt.Sprintf("gscript_%s_%s_amd64.zip", tagName, runtime.GOOS)

	assetID := 0
	fileSize := 0

	for _, asset := range repRel.Assets {
		if asset.GetName() == fileName {
			assetID = asset.GetID()
			fileSize = asset.GetSize()
			gse.LogInfof("Downloading New Binary: %s", asset.GetBrowserDownloadURL())
		}
	}

	if assetID == 0 {
		gse.LogCritf("No release found for your OS! WTF? Report this.")
	}

	_, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, "gen0cide", "gscript", assetID)
	if err != nil {
		gse.LogCritf("Github Error: %s", err.Error())
	}

	if len(redirectURL) == 0 {
		gse.LogCritf("There was an error retriving the release from Github.")
	}

	resp, err := http.Get(redirectURL)
	if err != nil {
		gse.LogCritf("Error Retrieving Binary: %s", err.Error())
	}
	defer resp.Body.Close()
	compressedFile, err := ioutil.ReadAll(resp.Body)

	gse.LogInfof("Uncompressing Binary...")

	var binary bytes.Buffer
	writer := bufio.NewWriter(&binary)
	compressedReader := bytes.NewReader(compressedFile)

	r, err := zip.NewReader(compressedReader, int64(fileSize))
	if err != nil {
		gse.LogCritf("Error Buffering Zip File: %s", err.Error())
	}

	for _, zf := range r.File {
		if zf.Name != "gscript" && zf.Name != "gscript.exe" {
			continue
		}
		src, err := zf.Open()
		if err != nil {
			gse.LogCritf("Error Unzipping File: %s", err.Error())
		}
		defer src.Close()

		io.Copy(writer, src)
	}

	reader := bufio.NewReader(&binary)

	err = update.Apply(reader, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			gse.LogCritf("Failed to rollback from bad update: %v", rerr)
		}
		gse.LogCritf("Update Failed - original version rolled back successfully.")
	}

	gse.LogInfof("Successfully updated to gscript v%s", tagName)
	return nil
}
