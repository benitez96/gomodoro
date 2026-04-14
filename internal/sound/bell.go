package sound

import (
	"encoding/binary"
	"io"
	"math"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

var (
	otoCtx   *oto.Context
	initOnce sync.Once
)

func initContext() {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return
	}
	<-ready
	otoCtx = ctx
}

func Bell() {
	initOnce.Do(initContext)
	if otoCtx == nil {
		return
	}

	p := otoCtx.NewPlayer(&toneReader{
		freq:     660,
		duration: 400 * time.Millisecond,
		rate:     44100,
	})
	p.Play()
	time.Sleep(450 * time.Millisecond)
	p.Close()
}

type toneReader struct {
	freq     float64
	duration time.Duration
	rate     int
	pos      int
}

func (t *toneReader) Read(buf []byte) (int, error) {
	total := int(float64(t.rate) * t.duration.Seconds())
	n := 0

	for n+1 < len(buf) && t.pos < total {
		sample := math.Sin(2 * math.Pi * t.freq * float64(t.pos) / float64(t.rate))

		fadeStart := total * 8 / 10
		if t.pos > fadeStart {
			sample *= float64(total-t.pos) / float64(total-fadeStart)
		}

		val := int16(sample * 16000)
		binary.LittleEndian.PutUint16(buf[n:], uint16(val))
		n += 2
		t.pos++
	}

	if t.pos >= total {
		return n, io.EOF
	}
	return n, nil
}
