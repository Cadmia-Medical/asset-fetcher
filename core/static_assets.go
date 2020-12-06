package core

import (
	"fmt"
	"os"
)

type StaticAssetFetcher interface {
	GetAssetTag(name string) (string, error)
	GetAsset(path string, name string) error
}

type AssetRefresher struct {
	Fetcher       StaticAssetFetcher
	AssetName     string
	LocalTag      string
	DownloadPath  string
	SymlinkTarget string
}

func (self *AssetRefresher) needsRefresh() (bool, error) {
	remoteTag, err := self.Fetcher.GetAssetTag(self.AssetName)
	if err != nil {
		return false, err
	}

	if remoteTag == self.LocalTag {
		return false, nil
	}

	return true, nil
}

func (self *AssetRefresher) pullAssets() error {
	newTag, err := self.Fetcher.GetAssetTag(self.AssetName)
	if err != nil {
		return err
	}
	err = self.Fetcher.GetAsset(self.DownloadPath, self.AssetName)
	if err != nil {
		return err
	}

	self.LocalTag = newTag
	return nil
}

func (self *AssetRefresher) relinkAssets() error {
	if err := os.Symlink(fmt.Sprintf("%v/%v", self.DownloadPath, self.AssetName), fmt.Sprintf("%v/%v", self.SymlinkTarget, self.AssetName)); err != nil {
		return err
	}
	return nil
}

func (self *AssetRefresher) Refresh() error {
	shouldRefresh, err := self.needsRefresh()
	if err != nil || !shouldRefresh {
		return err
	}

	if err := self.pullAssets(); err != nil {
		return err
	}

	if err := self.relinkAssets(); err != nil {
		return err
	}

	return nil
}
