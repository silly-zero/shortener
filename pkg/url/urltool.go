package urltool

import (
	"net/url"
	"path"

	"errors"
)

// GetBasePath 获取url的最后一段
func GetBasePath(targetUrl string) (string, error) {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("no host in targeturl")
	}
	return path.Base(myUrl.Path), nil
}
