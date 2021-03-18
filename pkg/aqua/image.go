package aqua

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (cli *APIClient) getMostRecentImageTag(images []ImageResult) (mostRecentTag string) {
	var mostRecentDate time.Time
	for _, image := range images {
		if image.ScanStatus == finished && image.Created.After(mostRecentDate) {
			mostRecentTag = image.Tag
			mostRecentDate = image.Created
		}
	}

	return mostRecentTag
}

func (cli *APIClient) GetImagesForRepository(ctx context.Context, registry, repository string) (images []ImageResult, err error) {
	images = make([]ImageResult, 0)
	page := 1

	endpoint := strings.Replace(getImages, "$REPOSITORYNAME", repository, 1)
	endpoint = strings.Replace(endpoint, "$REGISTRYNAME", registry, 1)

	endpoint = strings.Replace(endpoint, " ", "%20", -1)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var request *http.Request
		if request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s&page=%d&pagesize=50", cli.baseURL, endpoint, page), nil); err == nil {
			var body []byte
			if body, err = cli.executeRequest(request); err == nil {
				imagePage := &ImagePage{}
				if err = json.Unmarshal(body, imagePage); err == nil {
					if len(imagePage.Result) == 0 {
						break
					} else {
						images = append(images, imagePage.Result...)
						page++
					}
				} else {
					err = fmt.Errorf("error while parsing vulnerabilities from response - %s", err.Error())
					break
				}
			} else {
				err = fmt.Errorf("error while gathering images - %s", err.Error())
				break
			}
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
			break
		}
	}

	return images, err
}

func (cli *APIClient) StartImageScan(registryName string, imageName string) (err error) {
	endpoint := strings.Replace(postStartImageScan, "$REGISTRYNAME", registryName, 1)
	endpoint = strings.Replace(endpoint, "$IMAGENAME", imageName, 1)

	var request *http.Request
	if request, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", cli.baseURL, endpoint), nil); err == nil {
		var body []byte
		if body, err = cli.executeRequest(request); err == nil {
			fmt.Println(string(body))
		} else {
			err = fmt.Errorf("error while creating image scan - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while making request - %s", err.Error())
	}

	return err
}

type fullImageRescanReq struct {
	FullRescan bool             `json:"full_rescan"`
	Images     []ImageRescanReq `json:"images"`
}

type ImageRescanReq struct {
	Registry   string `json:"registry"`
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

func (cli *APIClient) StartFullImageRescan(registryName, imageName string) (err error) {
	endpoint := "/api/v1/images/rescan"

	reqBody := &fullImageRescanReq{
		FullRescan: true,
		Images: []ImageRescanReq{
			{
				Registry:   registryName,
				Name:       fmt.Sprintf("%s:master", imageName),
				Repository: imageName,
				Tag:        "master",
			},
		},
	}

	var body []byte
	if body, err = json.Marshal(reqBody); err == nil {
		bodyReader := bytes.NewReader(body)

		var request *http.Request
		if request, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", cli.baseURL, endpoint), bodyReader); err == nil {
			if body, err = cli.executeRequest(request); err == nil {

			} else {
				err = fmt.Errorf("error while creating image scan - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while building request")
		}
	} else {
		err = fmt.Errorf("error while marshalling body")
	}

	return err
}

// TODO getting auth issue on this endpoint alone for some reason?
//func (cli *APIClient) GetImageScanStatus(registryName string, imageName string) (status string, err error) {
//	endpoint := strings.Replace(getImageScanStatus, "$REGISTRYNAME", registryName, 1)
//	endpoint = strings.Replace(endpoint, "$IMAGENAME", imageName, 1)
//	fmt.Println(endpoint)
//
//	var request *http.Request
//	if request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", cli.baseURL, endpoint), nil); err == nil {
//		var body []byte
//		if body, err = cli.executeRequest(request); err == nil {
//			fmt.Println(string(body))
//		} else {
//			err = fmt.Errorf("error while getting image scan status - %s", err.Error())
//		}
//	} else {
//		err = fmt.Errorf("error while making request - %s", err.Error())
//	}
//
//	return status, err
//}
