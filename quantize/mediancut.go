// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package quantize

import (
	"image"
	"image/color"
	"sync"
)

type bucketPool struct {
	sync.Pool
	maxCap int
	m      sync.Mutex
}

func (p *bucketPool) getBucket(c int) colorBucket {
	p.m.Lock()
	if p.maxCap > c {
		p.maxCap = p.maxCap * 99 / 100
	}
	if p.maxCap < c {
		p.maxCap = c
	}
	maxCap := p.maxCap
	p.m.Unlock()
	val := p.Pool.Get()
	if val == nil || cap(val.(colorBucket)) < c {
		return make(colorBucket, maxCap)[0:c]
	}
	slice := val.(colorBucket)
	slice = slice[0:c]
	for i := range slice {
		slice[i] = colorPriority{}
	}
	return slice
}

var bpool bucketPool

// AggregationType specifies the type of aggregation to be done
type AggregationType uint8

const (
	Mode AggregationType = iota // pick the highest priority value
	Mean                        // weighted average all values
)

// MedianCutQuantizer implements the go draw.Quantizer interface using the Median Cut method
type MedianCutQuantizer struct {
	// The type of Aggregation to be used to find final colors
	Aggregation AggregationType
	// The Weighting function to use on each pixel
	Weighting func(image.Image, int, int) uint32
	// Whether need to add a transparent entry after conversion
	ReserveTransparent bool
}

// bucketize takes a bucket and performs median cut on it to obtain the target number of grouped buckets
func bucketize(colors colorBucket, num int) (buckets []colorBucket) {
	if len(colors) == 0 || num == 0 {
		return nil
	}
	bucket := colors
	buckets = make([]colorBucket, 1, num*2)
	buckets[0] = bucket

	for len(buckets) < num && len(buckets) < len(colors) { // Limit to palette capacity or number of colors
		bucket, buckets = buckets[0], buckets[1:]
		if len(bucket) < 2 {
			buckets = append(buckets, bucket)
			continue
		} else if len(bucket) == 2 {
			buckets = append(buckets, bucket[:1], bucket[1:])
			continue
		}

		left, right := bucket.partition()
		buckets = append(buckets, left, right)
	}
	return
}

// palettize finds a single color to represent a set of color buckets
func (q *MedianCutQuantizer) palettize(p color.Palette, buckets []colorBucket) color.Palette {
	for _, bucket := range buckets {
		switch q.Aggregation {
		case Mean:
			mean := bucket.mean()
			p = append(p, mean)
		case Mode:
			var best colorPriority
			for _, c := range bucket {
				if c.p > best.p {
					best = c
				}
			}
			p = append(p, best.RGBA)
		}
	}
	return p
}

// quantizeSlice expands the provided bucket and then palettizes the result
func (q *MedianCutQuantizer) quantizeSlice(p color.Palette, colors []colorPriority) color.Palette {
	numColors := cap(p) - len(p)
	reserveTransparent := q.ReserveTransparent
	if reserveTransparent {
		numColors--
	}
	buckets := bucketize(colors, numColors)
	p = q.palettize(p, buckets)
	return p
}

func colorAt(m image.Image, x int, y int) color.RGBA {
	switch i := m.(type) {
	case *image.YCbCr:
		yi := i.YOffset(x, y)
		ci := i.COffset(x, y)
		c := color.YCbCr{
			Y:  i.Y[yi],
			Cb: i.Cb[ci],
			Cr: i.Cr[ci],
		}
		return color.RGBA{R: c.Y, G: c.Cb, B: c.Cr, A: 255}
	case *image.RGBA:
		ci := i.PixOffset(x, y)
		return color.RGBA{R: i.Pix[ci+0], G: i.Pix[ci+1], B: i.Pix[ci+2], A: i.Pix[ci+3]}
	default:
		return color.RGBAModel.Convert(i.At(x, y)).(color.RGBA)
	}
}

// buildBucketMultiple creates a prioritized color slice with all the colors in
// the images.
func (q *MedianCutQuantizer) buildBucketMultiple(ms []image.Image) (bucket colorBucket) {
	if len(ms) < 1 {
		return colorBucket{}
	}

	bounds := ms[0].Bounds()
	size := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y) * 2
	sparseBucket := bpool.getBucket(size)

	for _, m := range ms {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				priority := uint32(1)
				if q.Weighting != nil {
					priority = q.Weighting(m, x, y)
				}
				c := colorAt(m, x, y)
				if c.A == 0 {
					if !q.ReserveTransparent {
						q.ReserveTransparent = true
					}
					continue
				}
				if priority != 0 {
					index := int(c.R)<<16 | int(c.G)<<8 | int(c.B)
					for i := 1; ; i++ {
						p := &sparseBucket[index%size]
						if p.p == 0 || p.RGBA == c {
							*p = colorPriority{p.p + priority, c}
							break
						}
						index += 1 + i
					}
				}
			}
		}
	}

	bucket = sparseBucket[:0]
	switch ms[0].(type) {
	case *image.YCbCr:
		for _, p := range sparseBucket {
			if p.p != 0 {
				r, g, b := color.YCbCrToRGB(p.R, p.G, p.B)
				bucket = append(bucket, colorPriority{p.p, color.RGBA{R: r, G: g, B: b, A: p.A}})
			}
		}
	default:
		for _, p := range sparseBucket {
			if p.p != 0 {
				bucket = append(bucket, p)
			}
		}
	}
	return
}

// QuantizeMultiple quantizes several images at once to a palette and returns
// the palette
func (q *MedianCutQuantizer) QuantizeMultiple(p color.Palette, m []image.Image) color.Palette {
	bucket := q.buildBucketMultiple(m)
	defer bpool.Put(&bucket)
	return q.quantizeSlice(p, bucket)
}
