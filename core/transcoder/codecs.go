package transcoder

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var supportedCodecs = map[string]string{
	"libx264":    "libx264",
	"h264_omx":   "omx",
	"h264_vaapi": "vaapi",
	"h264_nvenc": "NVIDEA nvenc",
	"h264_qsv":   "Intel Quicksync",
}

type Libx264Codec struct {
}

func (c *Libx264Codec) Name() string {
	return "libx264"
}

func (c *Libx264Codec) GlobalFlags() string {
	return ""
}

func (c *Libx264Codec) PixelFormat() string {
	return "yuv420p"
}

func (c *Libx264Codec) ExtraArguments() string {
	return strings.Join([]string{
		// "-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
	}, " ")
}

type OmxCodec struct {
}

func (c *OmxCodec) Name() string {
	return "h264_omx"
}

func (c *OmxCodec) GlobalFlags() string {
	return ""
}

func (c *OmxCodec) PixelFormat() string {
	return "yuv420p"
}

func (c *OmxCodec) ExtraArguments() string {
	return strings.Join([]string{
		"-tune", "zerolatency", // Option used for good for fast encoding and low-latency streaming (always includes iframes in each segment)
	}, " ")
}

type VaapiCodec struct {
}

func (c *VaapiCodec) Name() string {
	return "h264_vaapi"
}

func (c *VaapiCodec) GlobalFlags() string {
	flags := []string{
		// "-hwaccel", "vaapi",
		// "-hwaccel_output_format", "vaapi",
		"-vaapi_device", "/dev/dri/renderD128",
	}

	return strings.Join(flags, " ")
}

func (c *VaapiCodec) PixelFormat() string {
	return "vaapi_vld"
}

func (c *VaapiCodec) ExtraArguments() string {
	return "-vf 'format=nv12,hwupload'"
}

type NvencCodec struct {
}

func (c *NvencCodec) Name() string {
	return "h264_nvenc"
}

func (c *NvencCodec) GlobalFlags() string {
	flags := []string{
		"-hwaccel cuda",
	}

	return strings.Join(flags, " ")
}

func (c *NvencCodec) PixelFormat() string {
	return "yuv420p"
}

func (c *NvencCodec) ExtraArguments() string {
	return ""
}

type QuicksyncCodec struct {
}

func (c *QuicksyncCodec) Name() string {
	return "h264_qsv"
}

func (c *QuicksyncCodec) GlobalFlags() string {
	return ""
}

func (c *QuicksyncCodec) PixelFormat() string {
	return "nv12"
}

func (c *QuicksyncCodec) ExtraArguments() string {
	return ""
}

type Video4Linux struct{}

func (c *Video4Linux) Name() string {
	return "h264_v4l2m2m"
}

func (c *Video4Linux) GlobalFlags() string {
	return ""
}

func (c *Video4Linux) PixelFormat() string {
	return "yuv420p"
}

func (c *Video4Linux) ExtraArguments() string {
	return ""
}

// Codec represents a supported codec on the system.
type Codec interface {
	Name() string
	GlobalFlags() string
	PixelFormat() string
	ExtraArguments() string
}

// GetCodecs will return the supported codecs available on the system.
func GetCodecs(ffmpegPath string) []string {
	codecs := make([]string, 0)

	cmd := exec.Command(ffmpegPath, "-encoders")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln(err)
		return codecs
	}

	response := string(out)
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "H.264") {
			fields := strings.Fields(line)
			codec := fields[1]
			if _, supported := supportedCodecs[codec]; supported {
				codecs = append(codecs, codec)
			}
		}
	}

	return codecs
}
