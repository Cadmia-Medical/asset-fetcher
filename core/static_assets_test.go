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
		Fetcher      StaticAssetFetcher
		AssetName    string
		LocalTag     string
		DownloadPath string
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
				Fetcher:   TestFetcher{"abc123", false},
				AssetName: "test-asset",
				LocalTag:  "abc123",
			},
			wasUpdated: false,
			wantErr:    false,
		},
		{
			name: "Tags don't match, update assets",
			fields: fields{
				Fetcher:   TestFetcher{"def456", false},
				AssetName: "test-asset",
				LocalTag:  "abc123",
			},
			wasUpdated: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := AssetRefresher{
				Fetcher:      tt.fields.Fetcher,
				AssetName:    tt.fields.AssetName,
				LocalTag:     tt.fields.LocalTag,
				DownloadPath: "/",
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
		Fetcher      StaticAssetFetcher
		AssetName    string
		LocalTag     string
		DownloadPath string
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
				Fetcher:      TestFetcher{"def456", false},
				AssetName:    "test",
				LocalTag:     "abc123",
				DownloadPath: "/",
			},
			wantErr: false,
			newTag:  "def456",
		},
		{
			name: "Download throws error, tag not updated",
			fields: fields{
				Fetcher:      TestFetcher{"def456", true},
				AssetName:    "test",
				LocalTag:     "abc123",
				DownloadPath: "/",
			},
			wantErr: true,
			newTag:  "abc123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &AssetRefresher{
				Fetcher:      tt.fields.Fetcher,
				AssetName:    tt.fields.AssetName,
				LocalTag:     tt.fields.LocalTag,
				DownloadPath: tt.fields.DownloadPath,
			}
			if err := self.pullAssets(); (err != nil) != tt.wantErr {
				t.Errorf("%v - AssetRefresher.pullAssets() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			if tt.wantErr == false && self.LocalTag != tt.newTag {
				t.Errorf("%v - Asset Tag was not updated. Wanted = %v, got = %v", tt.name, tt.newTag, self.LocalTag)
			}
		})
	}
}
