package core

import (
	"os"
)

type StaticAssetFetcher interface {
	GetAssetTag(name string) (string, error)
	GetAsset(path string, name string) error
}

type AssetRefresher struct {
	Fetcher      StaticAssetFetcher
	AssetName    string
	LocalTag     string
	DownloadPath string
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
	if err := os.Symlink(self.DownloadPath, "/public"); err != nil {
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

/*
static asset cache
 - has a config the outlines which assert group (blob folder) to check
 - in a goroutine, checks to see if the asset group has been updated
   - checks ETag data?
   - eventually we can use this mechanism for client specific releases, feature flags, etc. (read from database which uuid to download)
 - If the ETag changes, pull down zip file, and unzip it to /available-assets directory
 - Once downloaded, replace symlink /static/appname -> /available-assets/appname-ETag
 - Webapp pulls static data from static directory
 - (Future Work) if error rate increases, roll back symlink change

*/
