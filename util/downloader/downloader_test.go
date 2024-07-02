package downloader

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloader_Chunked(t *testing.T) {
	chunks, totalBytes, err := Chunked("https://qhrenderpicoss.kujiale.com/r/2023/11/25/L3D327S41ENDPB7KJ2IUWIEGGLUFX7DW5YA8_1024x576.jpg")
	assert.Nil(t, err)
	assert.Equal(t, DefaultChunkNum, len(chunks))
	assert.Equal(t, 224971, totalBytes)
	assert.Equal(t, 224971, chunks[0].End)

	chunks, totalBytes, err = Chunked("https://qhrenderstorage-oss.kujiale.com/beautify/2024/04/30/MYYGOAAKTIUF2AABAAAAACY8.jpg")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chunks))
	assert.Equal(t, 1691754, totalBytes)
	assert.Equal(t, 1691754, chunks[1].End)

	chunks, totalBytes, err = Chunked("https://qhrenderpicoss.kujiale.com/r/2023/12/17/L3D124S41ENDPBQX37QUWJD4KLUFX7FBVMI8_4000x3000.jpg")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[2].End)
}

func TestMustChunked(t *testing.T) {
	chunks, totalSize := MustChunked("httx")
	assert.Equal(t, math.MaxInt, totalSize)
	assert.Equal(t, 1, len(chunks))
	assert.Equal(t, math.MaxInt, chunks[0].End)
}

func TestChunkedWithNumber(t *testing.T) {
	url := "https://qhrenderpicoss.kujiale.com/r/2023/12/17/L3D124S41ENDPBQX37QUWJD4KLUFX7FBVMI8_4000x3000.jpg"
	chunks, totalBytes, err := Chunked(url, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[1].End)

	chunks, totalBytes, err = Chunked(url, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[2].End)

	chunks, totalBytes, err = Chunked(url, 5)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(chunks))
	assert.Equal(t, 2334896, totalBytes)
	assert.Equal(t, 2334896, chunks[4].End)
	for i, chunk := range chunks {
		fmt.Println(i, chunk.Start, chunk.End)
	}
}
