package filex

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
)

// CleanName 清洗文件名称。
func CleanName(text string) (string, error) {
	invalidChars := "<>:\"/\\|?*"
	regexPattern := "[" + regexp.QuoteMeta(invalidChars) + "]"
	re := regexp.MustCompile(regexPattern)
	return re.ReplaceAllString(text, "_"), nil
}

// GetNameFromURL 从网址中提取可用的文件名称。
func GetNameFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "/" || fileName == "." {
		return "", fmt.Errorf("无法从网址中提取文件名: %s", rawURL)
	}

	return CleanName(fileName)
}
