package main

import (
	"gync/internal/client/github"
	"gync/internal/config"
	"gync/internal/index"
	"gync/internal/index/bptindex"
	"gync/internal/llog"
	"io"
	"os"
	"path/filepath"
)

// TODO update indexer
// dispatchOnce scan missing release and download missing release, meantime update indexer
func dispatchOnce(context *config.Config, indexer *bptindex.BptreeReleaseDirIndexer) error {
	repos := context.Repos

	for _, repo := range repos {
		// list all release
		allReleases, err := github.ListRelease(repo.Owner, repo.Name)
		if err != nil {
			llog.Errorf("dispatchOnce on repo %v failed: %v", repo, err)
			return err
		}

		if len(allReleases) == 0 {
			continue
		}

		// find release absents
		rr2Release := make(map[index.RepoRelease]github.Release)
		rrs := make([]index.RepoRelease, 0)
		for _, r := range allReleases {
			rr := index.RepoRelease{RepoOwner: repo.Owner, RepoName: repo.Name, ReleaseTime: r.Time, ReleaseName: r.Name}
			rr2Release[rr] = r
			rrs = append(rrs, rr)
		}

		llog.Debugln(rrs)

		absents, err := indexer.GetAbsent(rrs)
		if err != nil {
			llog.Errorf("dispatchOnce on repo %v failed: %v", repo, err)
			return err
		}

		for _, absent := range absents {
			release := rr2Release[absent]
			err := mkdirIfAbsent(context, &repo, &release)
			if err != nil {
				return err
			}
			for _, asset := range release.Assets {
				downloadedAsset, err := github.DownloadRelease(asset.DownloadUrl)
				if err != nil {
					llog.Errorf("failed to download repo{%v} downloadedAsset{%v}: %v", repo, downloadedAsset)
					return err
				}

				assetFileName := filepath.Join(context.RootDir, repo.Name, release.Name, asset.Name)
				assetFile, err := os.Create(assetFileName)
				if err != nil {
					llog.Errorf("failed to create asset file %v : %v", assetFile, err)
					return err
				}

				_, err = io.Copy(assetFile, downloadedAsset.Body)
				if err != nil {
					llog.Errorf("failed to copy asset to file: %v", err)
					return err
				}
			}
		}
	}

	return nil
}

// mkdir if release dir is absent
func mkdirIfAbsent(context *config.Config, repo *config.Repo, release *github.Release) error {
	rrDir := filepath.Join(context.RootDir, repo.Name, release.Name)
	if _, err := os.Stat(rrDir); os.IsNotExist(err) {
		err := os.MkdirAll(rrDir, os.ModePerm)
		if err != nil {
			llog.Errorf("failed to mkdir %v: %v", rrDir, err)
			return err
		}
	}

	return nil
}
