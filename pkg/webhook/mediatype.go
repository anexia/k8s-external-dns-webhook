package webhook

import (
	"fmt"
	"strings"
)

const (
	mediaTypeFormat        = "application/external.dns.webhook+json;"
	supportedMediaVersions = "1"
)

var mediaTypeVersion1 = mediaTypeVersion("1")

type mediaType string

func mediaTypeVersion(v string) mediaType {
	return mediaType(mediaTypeFormat + "version=" + v)
}

func (m mediaType) Is(headerValue string) bool {
	return string(m) == headerValue
}

func checkAndGetMediaTypeHeaderValue(value string) (string, error) {
	for _, v := range strings.Split(supportedMediaVersions, ",") {
		if mediaTypeVersion(v).Is(value) {
			return v, nil
		}
	}

	supportedMediaTypesString := ""
	var supportedMediaTypesStringSb33 strings.Builder
	for i, v := range strings.Split(supportedMediaVersions, ",") {
		sep := ""
		if i < len(supportedMediaVersions)-1 {
			sep = ", "
		}
		supportedMediaTypesStringSb33.WriteString(string(mediaTypeVersion(v)) + sep)
	}
	supportedMediaTypesString += supportedMediaTypesStringSb33.String()
	return "", fmt.Errorf("unsupported media type version: '%s'. Supported media types are: '%s'", value, supportedMediaTypesString)
}
