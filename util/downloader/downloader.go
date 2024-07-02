// Package downloader 网络资源下载器。
//
// 1. 支持手动或自动设置资源块数
//
// 2. 支持并发下载资源的字节分块
//
// 3. 支持下载分块出错时自动重试
//
// 4. 自动安全地合并分块
package downloader

import (
	"math"

	"github.com/imroc/req/v3"
)

const (
	// ChunkBytes 分块的内容长度(1MB)
	ChunkBytes int = 1048576

	// DefaultChunkNum 单个资源的默认分块数
	DefaultChunkNum int = 1

	// MaxChunkNum 单个资源的最大分块数量(即并发数)
	MaxChunkNum int = 64
)

// Chunk 下载块的字节区间。
type Chunk struct {
	Start int
	End   int
}

// Download 并行分块下载指定网址的资源到本地文件。
func Download(url string, chunkNum int, filename string) error {

	return nil
}

// MustChunked 返回资源的分块下载信息。
func MustChunked(url string, chunkNumber ...int) ([]Chunk, int) {
	chunks, totalBytes, err := Chunked(url, chunkNumber...)
	if err != nil {
		return []Chunk{
			{Start: 0, End: math.MaxInt},
		}, math.MaxInt
	}
	return chunks, totalBytes
}

// Chunked 根据内容长度，对资源自动进行分块，并返回分块数量。 TODO 不支持断点续传的资源。
func Chunked(url string, chunkNumber ...int) ([]Chunk, int, error) {
	// 获取内容长度
	resp, err := req.Head(url)
	if err != nil {
		return nil, 0, err
	}
	totalBytes := int(resp.ContentLength)

	// 设置块数
	var chunkNum int
	if len(chunkNumber) > 0 && chunkNumber[0] > 0 {
		chunkNum = chunkNumber[0]
	} else {
		// 不足分块的内容，采用默认分块数
		if totalBytes <= ChunkBytes {
			return []Chunk{
				{Start: 0, End: totalBytes},
			}, totalBytes, nil
		}

		// 自动计算分块数量
		chunkNum = int(math.Ceil(float64(totalBytes) / float64(ChunkBytes)))
	}

	// 计算分块大小
	chunkBytes := totalBytes / chunkNum

	// 计算每个块的字节区间
	chunks := make([]Chunk, chunkNum)
	for i := 0; i < chunkNum; i++ {
		switch i {
		case 0:
			// 第一块
			chunks[i].Start = 0
			chunks[i].End = chunkBytes
		case chunkNum - 1:
			// 最后一块
			chunks[i].Start = chunks[i-1].End + 1
			//chunks[i].End = totalBytes - 1
			chunks[i].End = totalBytes
		default:
			//	中间块
			chunks[i].Start = chunks[i-1].End + 1
			chunks[i].End = chunks[i].Start + chunkBytes
		}
	}

	return chunks, totalBytes, nil
}
