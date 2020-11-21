package core

import (
	"errors"
	"testing"
)

type TestFetcher struct {
	tag         string
	shouldError bool
}

func (self TestFetcher) GetAssetTag(name string) (string, error) {
	return self.tag, nil
}

func (self TestFetcher) GetAsset(path string, name string) error {
	if self.shouldError {
		return errors.New("Some error")
	}
	return nil
}

func TestAssetRefresher_needsRefresh(t *testing.T) {
	type fields struct {
		fetcher      StaticAssetFetcher
		assetName    string
		localTag     string
		downloadPath string
	}
	tests := []struct {
		name       string
		fields     fields
		wasUpdated bool
		wantErr    bool
	}{
		{
			name: "Tags match, do nothing",
			fields: fields{
				fetcher:   TestFetcher{"abc123", false},
				assetName: "test-asset",
				localTag:  "abc123",
			},
			wasUpdated: false,
			wantErr:    false,
		},
		{
			name: "Tags don't match, update assets",
			fields: fields{
				fetcher:   TestFetcher{"def456", false},
				assetName: "test-asset",
				localTag:  "abc123",
			},
			wasUpdated: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := AssetRefresher{
				fetcher:      tt.fields.fetcher,
				assetName:    tt.fields.assetName,
				localTag:     tt.fields.localTag,
				downloadPath: "/",
			}
			wasUpdated, err := self.needsRefresh()

			if wasUpdated != tt.wasUpdated {
				t.Errorf("AssetRefresher.CheckForUpdate() wasUpdated = %v, wantErr %v", wasUpdated, tt.wasUpdated)
			}

			if err != nil && tt.wantErr {
				t.Errorf("AssetRefresher.CheckForUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssetRefresher_pullAssets(t *testing.T) {
	type fields struct {
		fetcher      StaticAssetFetcher
		assetName    string
		localTag     string
		downloadPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		newTag  string
	}{
		{
			name: "Downloads file, updates tag",
			fields: fields{
				fetcher:      TestFetcher{"def456", false},
				assetName:    "test",
				localTag:     "abc123",
				downloadPath: "/",
			},
			wantErr: false,
			newTag:  "def456",
		},
		{
			name: "Download throws error, tag not updated",
			fields: fields{
				fetcher:      TestFetcher{"def456", true},
				assetName:    "test",
				localTag:     "abc123",
				downloadPath: "/",
			},
			wantErr: true,
			newTag:  "abc123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &AssetRefresher{
				fetcher:      tt.fields.fetcher,
				assetName:    tt.fields.assetName,
				localTag:     tt.fields.localTag,
				downloadPath: tt.fields.downloadPath,
			}
			if err := self.pullAssets(); (err != nil) != tt.wantErr {
				t.Errorf("%v - AssetRefresher.pullAssets() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			if tt.wantErr == false && self.localTag != tt.newTag {
				t.Errorf("%v - Asset Tag was not updated. Wanted = %v, got = %v", tt.name, tt.newTag, self.localTag)
			}
		})
	}
}
