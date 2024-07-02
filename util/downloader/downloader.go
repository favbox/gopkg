// Package downloader 网络资源下载器。
//
// 1. 支持手动或自动设置资源块数
//
// 2. 支持并发下载资源的字节分块
//
// 3. 支持下载分块出错时自动重试
//
// 4. 自动安全地合并分块
//
// TODO 暂不支持断点续传的资源。
package downloader

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"github.com/bytedance/gopkg/util/xxhash3"
	"github.com/favbox/gopkg/util/filex"
	"github.com/imroc/req/v3"
	"github.com/schollz/progressbar/v3"
	"github.com/sourcegraph/conc/pool"
)

const (
	// ChunkBytes 分块的内容长度(1MB)
	ChunkBytes int = 1048576

	// DefaultChunkNum 单个资源的默认分块数
	DefaultChunkNum int = 1

	// MaxChunkNum 单个资源的最大分块数量(即并发数)
	MaxChunkNum int = 64
)

var progress bool

type (
	// Chunk 待下载数据块的字节区间。
	Chunk struct {
		Start int
		End   int
	}

	// ChunkedFile 已下载的分块文件及序号
	ChunkedFile struct {
		Index int
		Name  string
	}
)

// WithProgress 设置是否使用进度条。全局变量，默认不启用。
func WithProgress(b bool) {
	progress = b
}

// Download 并行分块下载指定网址的资源到本地。
func Download(url string, filename ...string) error {
	return DownloadWithChunks(url, 0, filename...)
}

// DownloadWithChunks 并行分块下载指定网址的资源到本地文件。
func DownloadWithChunks(url string, chunkNum int, filename ...string) (err error) {
	// 设置文件名称
	var name string
	if len(filename) > 0 && len(filename[0]) > 0 {
		name = filename[0]
	} else {
		name, err = filex.GetNameFromURL(url)
		if err != nil {
			return err
		}
	}

	chunks, _ := mustChunked(url, chunkNum)
	//log.Println("字节数:", totalBytes, "并行数:", len(chunks))
	var bar = progressbar.New(len(chunks))
	ch := make(chan int, len(chunks))
	go func(ch <-chan int) {
		for range ch {
			err := bar.Add(1)
			if err != nil {
				log.Fatal(fmt.Errorf("bar add 失败: %v", err))
			}
		}
	}(ch)

	g := pool.NewWithResults[*ChunkedFile]().WithErrors().WithFirstError()
	for i, chunk := range chunks {
		g.Go(func() (*ChunkedFile, error) {
			return downChunk(url, i, chunk, ch)
		})
	}

	// 处理并发结果
	chunkFiles, err := g.Wait()
	if err != nil {
		return err
	}

	// 重排乱序的并发块文件
	sort.Slice(chunkFiles, func(i, j int) bool {
		return chunkFiles[i].Index < chunkFiles[j].Index
	})

	// 合并临时文件块
	dst, err := os.Create(name)
	if err != nil {
		return err
	}
	// 拷贝分块文件至目标文件
	for _, cf := range chunkFiles {
		in, _ := os.Open(cf.Name)
		_, err := io.Copy(dst, in)
		in.Close()
		os.Remove(in.Name())
		if err != nil {
			return fmt.Errorf("无法合并分块文件: %s, %v", cf.Name, err)
		}
	}
	return nil
}

func downChunk(url string, i int, chunk Chunk, ch chan<- int) (*ChunkedFile, error) {
	defer func() {
		if progress {
			ch <- 1
		}
	}()

	resp, err := req.R().
		SetHeader("Range", fmt.Sprintf("bytes=%v-%v", chunk.Start, chunk.End)).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("资源响应状态异常(%s: %s)", resp.GetStatus(), url)
	}

	// 写入分块临时文件
	name := makeFilename(url, i)
	file, err := os.CreateTemp("", name)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, err
	}
	//fmt.Println("分块长度", i, chunk.Start, chunk.End)

	chunkFile := &ChunkedFile{
		Index: i,
		Name:  file.Name(),
	}
	file.Close()

	return chunkFile, nil
}

func makeFilename(url string, i int) string {
	name := fmt.Sprintf("%d-%d", xxhash3.HashString(url), i)
	return name
}

// mustChunked 返回资源的分块下载信息。
func mustChunked(url string, chunkNumber ...int) ([]Chunk, int) {
	chunks, totalBytes, err := chunked(url, chunkNumber...)
	if err != nil {
		return []Chunk{
			{Start: 0, End: math.MaxInt},
		}, math.MaxInt
	}
	return chunks, totalBytes
}

// chunked 根据内容长度，对资源自动进行分块，并返回分块数量。 TODO 不支持断点续传的资源。
func chunked(url string, chunkNumber ...int) ([]Chunk, int, error) {
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
