package utils

import (
	"io"
	"net/http"
	"regexp"
)

func ExtractImageUrlFromFlatMap(data map[string]interface{}) string {
	imageURLs := []string{}

	re := regexp.MustCompile(`^attr_\d+_image$`)

	for key, value := range data {
		if re.MatchString(key) {
			if str, ok := value.(string); ok && str != "" {
				imageURLs = append(imageURLs, str)
			}
		}
	}

	if len(imageURLs) == 0 {
		return ""
	}

	return imageURLs[0]
}
func GetImageBytesFromFlatMap(data map[string]interface{}) []byte {
	res := ExtractImageUrlFromFlatMap(data)
	resp, err := http.Get(res)
	if err != nil {
		return []byte{}
	}
	defer resp.Body.Close()

	imgBytes, _ := io.ReadAll(resp.Body)
	return imgBytes
}
