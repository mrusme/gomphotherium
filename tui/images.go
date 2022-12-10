package tui

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"image/color"
	"strings"

	"github.com/peterbourgon/diskv/v3"

	"github.com/eliukblau/pixterm/pkg/ansimage"
)

type Images struct {
	diskCache *diskv.Diskv
}

func cacheKey(url string, y int, x int) string {
	sha := sha256.New()
	sha.Write([]byte(url))

	return fmt.Sprintf("%x-%d-%d", sha.Sum(nil), y, x)
}

func NewImages(directory string) (*Images, error) {
	images := &Images{}

	transform := func(s string) []string {
		parts := []string{}
		parts = append(parts, string(s[0:2]))
		parts = append(parts, string(s[2:4]))
		return parts
	}

	images.diskCache = diskv.New(diskv.Options{
		BasePath:     directory,
		Transform:    transform,
		CacheSizeMax: 1024 * 1024,
	})

	test := make([]string, 1)
	test[0] = "test"
	err := images.saveToDiskCache("test", &test)
	if err != nil {
		return nil, err
	} else {
		return images, nil
	}
}

func (i *Images) loadFromDiskCache(key string) (*[]string, error) {
	cacheStream, err := i.diskCache.ReadStream(key, false)
	if err != nil {
		return nil, err
	}

	var storedSlices []string
	gob.NewDecoder(cacheStream).Decode(&storedSlices)
	return &storedSlices, nil
}

func (i *Images) saveToDiskCache(key string, data *[]string) error {
	var inputBuffer bytes.Buffer
	gob.NewEncoder(&inputBuffer).Encode(data)
	return i.diskCache.Write(key, inputBuffer.Bytes())
}

func (i *Images) ImageAtSize(url string, y int, x int, afterLoad func(loaded *[]string) *[]string) *[]string {
	if afterLoad == nil {
		afterLoad = func(loaded *[]string) *[]string {
			return loaded
		}
	}
	key := cacheKey(url, y, x)

	cachedSlice, err := i.loadFromDiskCache(key)
	if err == nil {
		return afterLoad(cachedSlice)
	}

	pix, err := ansimage.NewScaledFromURL(
		url,
		y,
		x,
		color.Transparent,
		ansimage.ScaleModeResize,
		ansimage.NoDithering,
	)

	if err == nil {
		imageString := pix.RenderExt(false, false)
		split := strings.Split(imageString, "\n")
		i.saveToDiskCache(key, &split)

		return afterLoad(&split)
	} else {
		placeholder := make([]string, 0)
		for i := 0; i < y; i++ {
			placeholder = append(placeholder, strings.Repeat("?", x))
		}
		i.saveToDiskCache(key, &placeholder)

		return afterLoad(&placeholder)
	}
}
