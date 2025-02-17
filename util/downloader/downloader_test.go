package downloader

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownloader_Chunked(t *testing.T) {
	chunks, totalBytes, err := chunked("https://qhrenderpicoss.kujiale.com/r/2023/11/25/L3D327S41ENDPB7KJ2IUWIEGGLUFX7DW5YA8_1024x576.jpg")
	assert.Nil(t, err)
	assert.Equal(t, DefaultChunkNum, len(chunks))
	assert.Equal(t, 224971, totalBytes)
	assert.Equal(t, 224971, chunks[0].End)

	chunks, totalBytes, err = chunked("https://qhrenderstorage-oss.kujiale.com/beautify/2024/04/30/MYYGOAAKTIUF2AABAAAAACY8.jpg")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chunks))
	assert.Equal(t, 1691754, totalBytes)
	assert.Equal(t, 1691754, chunks[1].End)

	chunks, totalBytes, err = chunked("https://qhrenderpicoss.kujiale.com/r/2023/12/17/L3D124S41ENDPBQX37QUWJD4KLUFX7FBVMI8_4000x3000.jpg")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[2].End)
}

func TestMustChunked(t *testing.T) {
	chunks, totalSize := mustChunked("httx")
	assert.Equal(t, math.MaxInt, totalSize)
	assert.Equal(t, 1, len(chunks))
	assert.Equal(t, math.MaxInt, chunks[0].End)
}

func TestChunkedWithNumber(t *testing.T) {
	url := "https://qhrenderpicoss.kujiale.com/r/2023/12/17/L3D124S41ENDPBQX37QUWJD4KLUFX7FBVMI8_4000x3000.jpg"
	chunks, totalBytes, err := chunked(url, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[1].End)

	chunks, totalBytes = mustChunked(url, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[1].End)

	chunks, totalBytes, err = chunked(url, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[2].End)

	chunks, totalBytes, err = chunked(url, 5)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[4].End)
	for i, chunk := range chunks {
		fmt.Println(i, chunk.Start, chunk.End)
	}
}

func TestDownload(t *testing.T) {
	url := "https://qhrenderpicoss.kujiale.com/r/2023/12/17/L3D124S41ENDPBQX37QUWJD4KLUFX7FBVMI8_4000x3000.jpg"
	chunkNums := []int{0, 1, 3, 5}

	for _, chunkNum := range chunkNums {
		start := time.Now()
		err := DownloadWithChunks(url, chunkNum, fmt.Sprintf("test-%d.jpg", chunkNum))
		assert.Nil(t, err)
		fmt.Println(chunkNum, time.Since(start))
	}
}

func TestDownloadVideo(t *testing.T) {
	WithProgress(true)

	videoUrl := "https://sns-video-hw.xhscdn.com/1000g00g2ku4lp5cj20005o0d5ij09fivqklvkug"
	start := time.Now()
	err := Download(videoUrl, "四个动作，有效拉伸✨ ₊⁺细腿直腿一起get.mp4")
	fmt.Println("耗时", time.Since(start))
	assert.Nil(t, err)

}
