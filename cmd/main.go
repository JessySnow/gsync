package main

import (
	"github.com/common-nighthawk/go-figure"
	config2 "gync/internal/config"
	"gync/internal/index"
	"gync/internal/index/bptindex"
	"gync/internal/llog"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var context *config2.Config
var indexer *bptindex.BptreeReleaseDirIndexer

func init() {
	initConfig()
	initIndexer()
}

func main() {
	gsync := figure.NewFigure("GSYNC", "", true)
	gsync.Print()
	dispatch(context, indexer)
}

// initConfig read config file and parse it to global config -- context
func initConfig() {
	path, err := os.Executable()
	if err != nil {
		llog.Fatalf("can't get executable path: %v", err)
	}

	// open config file, get config reference
	dir := filepath.Dir(path)
	fileName := "gsync.json"
	filePath := filepath.Join(dir, fileName)
	configFile, err := os.Open(filePath)
	if err != nil {
		llog.Fatalf("can't open config file %s: %v", filePath, err)
	}
	defer configFile.Close()

	// read and parse config file
	bytes, err := io.ReadAll(configFile)
	if err != nil {
		llog.Fatalf("can't read config file %s: %v", filePath, err)
	}
	context, err = config2.Parse(bytes)
	if err != nil {
		llog.Fatalf("can't parse config file %s: %v", filePath, err)
	}
}

// initIndexer new indexer and bulkLoad reporelease to indexer
// ---- RootDir - DirDepth 0
// ---- ---- RepoOwnerA - DirDepth 1
// ---- ---- ---- Repo1 - DirDepth 2
// ---- ---- ---- ---- Release_1 - DirDepth 3
// ---- ---- ---- ---- Release_2 - DirDepth 3
// ---- ---- RepoOwnerB - DirDepth 1
// ---- ---- ---- Repo1 - DirDepth 2
// ---- ---- ---- ---- Release1 - DirDepth 3
func initIndexer() {
	rootDir := context.RootDir
	if len(rootDir) == 0 {
		llog.Fatalln("gsync root dir is empty")
	}
	var err error
	indexer, err = bptindex.New()
	if err != nil {
		llog.Fatalf("can't init indexer: %v", err)
	}

	// bulk load rootDir to indexer
	doInitIndexer(rootDir, rootDir, 0)
}

func doInitIndexer(rootDir string, currentDir string, dirDepth int) {
	if dirDepth == 3 {
		// read metadata
		metadata, err := os.ReadFile(filepath.Join(currentDir, ".meta"))
		dirName, err := filepath.Rel(rootDir, currentDir)
		if err != nil {
			llog.Fatalf("can't get relative path: %v", err)
		}
		parts := strings.Split(dirName, string(os.PathSeparator))
		repoOwner := parts[0]
		repoName := parts[1]
		releaseName := parts[2]
		releaseTime := string(metadata)
		rr := index.RepoRelease{RepoOwner: repoOwner, RepoName: repoName, ReleaseName: releaseName, ReleaseTime: releaseTime}
		llog.Debugf("adding repo release to indexer: %v", rr)
		_, err = indexer.Add(&rr)
		if err != nil {
			llog.Fatalf("can't add repo release: %v", err)
		}
		return
	}

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		llog.Fatalf("can't read rootDir %s: %v", rootDir, err)
	}

	for _, entry := range entries {
		doInitIndexer(rootDir, filepath.Join(currentDir, entry.Name()), dirDepth+1)
	}
}
