// +build windows
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca"
)

type WAVEFormat struct {
	FormatTag      uint16
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
	BitsPerSample  uint16
	DataTag        [4]byte // 'data'
	DataSize       uint32
	RawData        []byte
}

func readFile(filename string) (audio *WAVEFormat, err error) {
	var file []byte
	if file, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	audio = &WAVEFormat{}
	reader := bytes.NewReader(file)
	binary.Read(io.NewSectionReader(reader, 20, 2), binary.LittleEndian, &audio.FormatTag)
	binary.Read(io.NewSectionReader(reader, 22, 2), binary.LittleEndian, &audio.Channels)
	binary.Read(io.NewSectionReader(reader, 24, 4), binary.LittleEndian, &audio.SamplesPerSec)
	binary.Read(io.NewSectionReader(reader, 28, 4), binary.LittleEndian, &audio.AvgBytesPerSec)
	binary.Read(io.NewSectionReader(reader, 32, 2), binary.LittleEndian, &audio.BlockAlign)
	binary.Read(io.NewSectionReader(reader, 34, 2), binary.LittleEndian, &audio.BitsPerSample)
	binary.Read(io.NewSectionReader(reader, 36, 4), binary.LittleEndian, &audio.DataTag)
	binary.Read(io.NewSectionReader(reader, 40, 4), binary.LittleEndian, &audio.DataSize)

	buf := new(bytes.Buffer)
	io.Copy(buf, io.NewSectionReader(reader, 44, int64(audio.DataSize)))
	audio.RawData = buf.Bytes()

	if len(audio.RawData) == 0 {
		err = fmt.Errorf("empty data")
	}
	return
}

func main() {
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
	return
}
func run(args []string) (err error) {
	var filenameFlag string
	var audio *WAVEFormat

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&filenameFlag, "f", "", "Specify WAVE format audio (e.g. music.wav)")
	f.Parse(args[1:])

	if filenameFlag == "" {
		return
	}
	if audio, err = readFile(filenameFlag); err != nil {
		return
	}
	return renderSharedTimerDriven(audio)
}

func renderSharedTimerDriven(audio *WAVEFormat) (err error) {
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}

	var de *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &de); err != nil {
		return
	}
	defer de.Release()

	var mmd *wca.IMMDevice
	if err = de.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if err = ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return
	}
	fmt.Printf("Capturing what you hear from\n%s\n", pv.String())

	var ac *wca.IAudioClient
	if err = mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &ac); err != nil {
		return
	}
	defer ac.Release()

	var wfx *wca.WAVEFORMATEX
	if err = ac.GetMixFormat(&wfx); err != nil {
		return
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(wfx)))

	if wfx.WFormatTag != wca.WAVE_FORMAT_PCM {
		wfx.WFormatTag = 1
		wfx.NSamplesPerSec = audio.SamplesPerSec
		wfx.WBitsPerSample = audio.BitsPerSample
		wfx.NChannels = audio.Channels
		wfx.NBlockAlign = audio.BlockAlign
		wfx.NAvgBytesPerSec = audio.AvgBytesPerSec
		wfx.CbSize = 0
	}

	fmt.Println("--------")
	fmt.Printf("Format: PCM %d bit signed integer\n", wfx.WBitsPerSample)
	fmt.Printf("Rate: %d Hz\n", wfx.NSamplesPerSec)
	fmt.Printf("Channels: %d\n", wfx.NChannels)
	fmt.Println("--------")

	var defaultPeriod int64
	var minimumPeriod int64
	var capturingPeriod time.Duration
	if err = ac.GetDevicePeriod(&defaultPeriod, &minimumPeriod); err != nil {
		return
	}
	capturingPeriod = time.Duration(int(defaultPeriod) * 100)
	fmt.Printf("Default capturing period: %d ms\n", capturingPeriod/time.Millisecond)

	if err = ac.Initialize(wca.AUDCLNT_SHAREMODE_SHARED, 0, 500*10000, 0, wfx, nil); err != nil {
		return
	}

	var bufferFrameSize uint32
	if err = ac.GetBufferSize(&bufferFrameSize); err != nil {
		return
	}
	fmt.Printf("Allocated buffer size: %d\n", bufferFrameSize)

	var arc *wca.IAudioRenderClient
	if err = ac.GetService(wca.IID_IAudioRenderClient, &arc); err != nil {
		return
	}
	defer arc.Release()

	if err = ac.Start(); err != nil {
		return
	}
	fmt.Println("Start rendering audio with shared-timer-driven mode")

	time.Sleep(capturingPeriod)

	var data *byte
	var offset int
	for {
		if offset >= int(audio.DataSize) {
			break
		}
		var padding uint32
		var availableFrameSize uint32
		if err = ac.GetCurrentPadding(&padding); err != nil {
			return
		}
		availableFrameSize = bufferFrameSize - padding
		if availableFrameSize == 0 {
			continue
		}
		if err = arc.GetBuffer(availableFrameSize, &data); err != nil {
			return
		}
		start := unsafe.Pointer(data)
		lim := int(availableFrameSize) * int(wfx.NBlockAlign)
		remaining := int(audio.DataSize) - offset
		if remaining < lim {
			lim = remaining
		}
		var b *byte
		for n := 0; n < lim; n++ {
			b = (*byte)(unsafe.Pointer(uintptr(start) + uintptr(n)))
			*b = audio.RawData[offset+n]
		}
		offset += lim
		if err = arc.ReleaseBuffer(availableFrameSize, 0); err != nil {
			return
		}
		time.Sleep(80 * time.Millisecond)
	}

	if err = ac.Stop(); err != nil {
		return
	}
	fmt.Println("Stopping capturing loopback audio")
	return
}
