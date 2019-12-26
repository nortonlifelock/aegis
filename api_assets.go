package nexpose

import (
	"context"
	"fmt"
	"net/http"
	"github.com/nortonlifelock/log"
	"strconv"
)

// GetAssets loads all assets from nexpose through the api.
// The returned data can be sorted by passing a sort string that uses the formatting available in the nexpose API.
// Sort for this method is defaulted to IP when no sort is passed
func (a *Session) GetAssets(ctx context.Context, sort string) (<-chan *Asset, error) {
	var devices = make(chan *Asset)
	var err error

	go func(devices chan<- *Asset) {
		defer handleRoutinePanic(a.lstream)
		defer close(devices)

		fields := map[string]string{"sort": sort}
		// Set the default sorting for the endpoint
		if len(sort) == 0 {
			fields["sort"] = "ip"
		}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfAsset{}, apiGetAssets, fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var device Asset
						if device, ok = d.(Asset); ok {
							devices <- &device
						} else {
							a.lstream.Send(log.Error("unable to cast paged return as device", err))
						}
					} else {
						return
					}
				}
			}
		} else {
			a.lstream.Send(log.Error("error executing paged call against nexpose", err))
		}
	}(devices)

	return devices, err
}

// GetAsset loads a specific asset from the nexpose API by ID
func (a *Session) GetAsset(id int) (asset *Asset, err error) {

	asset = &Asset{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetAsset, encode(strconv.Itoa(id))), nil, nil, asset)

	return asset, err
}

// GetAssetsForSite loads all assets for a given site from the nexpose api using the site ID that is passed.
// The returned data can be sorted by passing a sort string that uses the formatting available in the nexpose API.
// Sort for this method is defaulted to IP when no sort is passed
func (a *Session) GetAssetsForSite(ctx context.Context, siteID string, sort string) (<-chan *Asset, error) {
	var devices = make(chan *Asset)
	var err error

	go func(devices chan<- *Asset) {
		defer handleRoutinePanic(a.lstream)
		defer close(devices)

		fields := map[string]string{"sort": sort}
		// Set the default sorting for the endpoint
		if len(sort) == 0 {
			fields["sort"] = "ip"
		}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfAsset{}, fmt.Sprintf(apiGetSiteAssets, encode(siteID)), fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var device Asset
						if device, ok = d.(Asset); ok {
							devices <- &device
						} else {
							a.lstream.Send(log.Error("unable to cast paged return as device", err))
						}
					} else {
						return
					}
				}
			}
		} else {
			a.lstream.Send(log.Error("error executing paged call against nexpose", err))
		}
	}(devices)

	return devices, err
}
