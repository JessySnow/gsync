package main

import (
	"gync/internal/client/github"
	"gync/internal/config"
	"gync/internal/helper"
	"gync/internal/index"
	"gync/internal/index/bptindex"
	"gync/internal/llog"
	"io"
	"path/filepath"
	"time"
)

// dispatch schedule run dispatchOnce and handle error
func dispatch(context *config.Config, indexer *bptindex.BptreeReleaseDirIndexer) {
	for {
		err := dispatchOnce(context, indexer)
		if err != nil {
			llog.Errorf("dispatch failed try dispatch again after %v minutes: %v", context.SyncInterval, err)
		}

		time.Sleep(time.Duration(context.SyncInterval) * time.Minute)
	}
}

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
			if 0 == len(r.Name) {
				r.Name = r.TagNme
			}
			rr := index.RepoRelease{RepoOwner: repo.Owner, RepoName: repo.Name, ReleaseTime: r.Time, ReleaseName: r.Name}
			rr2Release[rr] = r
			rrs = append(rrs, rr)
		}
		absents, err := indexer.GetAbsent(rrs)
		if err != nil {
			llog.Errorf("dispatchOnce on repo %v failed: %v", repo, err)
			return err
		}
		if 0 == len(absents) {
			llog.Infof("all release are downloaded")
			continue
		}

		// download absent release
		for _, absent := range absents {
			// 0. create repo release dir
			release := rr2Release[absent]
			repoReleaseDir := filepath.Join(context.RootDir, repo.Owner, repo.Name, release.Name)
			err := helper.MkdirIfAbsent(repoReleaseDir)
			if err != nil {
				llog.Errorf("failed to make repo release dir: %v", err)
				return err
			}

			// 1. download all assets
			for _, asset := range release.Assets {
				downloadedAsset, err := github.DownloadRelease(asset.DownloadUrl)
				if err != nil {
					llog.Errorf("failed to download repo{%v} downloadedAsset{%v}: %v", repo, downloadedAsset)
					return err
				}

				assetFileName := filepath.Join(context.RootDir, repo.Owner, repo.Name, release.Name, asset.Name)
				assetFile, err := helper.CreateFileIfAbsent(assetFileName)
				if err != nil {
					llog.Errorf("failed to create asset file %v : %v", assetFile, err)
					return err
				}

				_, err = io.Copy(assetFile, downloadedAsset.Body)
				downloadedAsset.Body.Close()
				assetFile.Close()
				if err != nil {
					llog.Errorf("failed to copy asset to file: %v", err)
					return err
				}
			}

			// 2. write metadata
			metaFileName := filepath.Join(context.RootDir, repo.Owner, repo.Name, release.Name, ".meta")
			metaFile, err := helper.CreateFileIfAbsent(metaFileName)
			if err != nil {
				llog.Errorf("failed to create metadata file for repo %v, release: %v: %v", repo, release, err)
				return err
			}
			_, err = metaFile.Write([]byte(release.Time))
			if err != nil {
				llog.Errorf("failed to write metadata to file : %v", err)
				return err
			}
			metaFile.Close()

			// 3. update indexer
			_, err = indexer.Add(&absent)
			if err != nil {
				llog.Errorf("failed to add new release %v to indexer: %v", absent, err)
				return err
			}
		}

	}

	return nil
}
