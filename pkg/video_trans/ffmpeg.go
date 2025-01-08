package video_trans

import (
	"fmt"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"math/rand"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func TransWithProgress(inFileName, outFileName string, onProgress func(progress string), onErr func(err error)) {
	a, err := ffmpeg.Probe(inFileName)
	if err != nil {
		onErr(err)
	}
	totalDuration := gjson.Get(a, "format.duration").Float()

	err = ffmpeg.Input(inFileName).
		Output(outFileName,
			ffmpeg.KwArgs{
				"profile:v":     "baseline", // 设置视频配置文件为baseline
				"level":         "3.0",      // 设置视频编码级别
				"start_number":  "0",        // 设置分片开始的序号
				"hls_time":      "10",       // 每个HLS分片的时长为10秒
				"hls_list_size": "0",        // 不限制播放列表中的M3U8片段数量
				"f":             "hls",      // 设置输出格式为HLS
			},
		).
		GlobalArgs("-progress", "unix://"+TempSock(totalDuration, onProgress, onErr)).
		OverWriteOutput().
		Run()
	if err != nil {
		onErr(err)
	}
}

func TempSock(totalDuration float64, onProgress func(progress string), onErr func(err error)) string {
	// serve

	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		onErr(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			onErr(err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f%", float64(c)/totalDuration/1000000*100)
			}
			if strings.Contains(data, "progress=end") {
				cp = "100.00%"
			}
			if cp == "" {
				cp = "0.00%"
			}
			if cp != progress {
				progress = cp
				onProgress(cp)
			}
		}
	}()

	return sockFileName
}
