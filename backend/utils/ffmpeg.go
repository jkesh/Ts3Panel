// backend/utils/ffmpeg.go
package utils

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

// AudioStream 代表一个正在播放的音频流
type AudioStream struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
}

// StartFFmpeg 启动 FFmpeg 进程，返回 PCM 音频流
// url: 音乐链接或文件路径
// qn: TS3 通常需要 48kHz, 双声道(stereo), 16bit 小端序(s16le) 的 PCM 数据
func StartFFmpeg(url string) (*AudioStream, error) {
	// 构造 FFmpeg 命令
	// -i url       : 输入源
	// -f s16le     : 输出格式为 signed 16-bit little endian (PCM)
	// -ar 48000    : 采样率 48kHz (Opus 标准)
	// -ac 2        : 声道数 2 (立体声)
	// -vn          : 不处理视频
	// -loglevel error : 减少日志输出
	// pipe:1       : 输出到标准输出
	args := []string{
		"-re", // (可选) 按原速读取，如果是直播流或本地文件建议加上，避免读取过快
		"-i", url,
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"-vn",
		"-loglevel", "error",
		"pipe:1",
	}

	cmd := exec.Command("ffmpeg", args...)

	// 获取标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	// 启动进程
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	log.Printf("[FFmpeg] Started streaming: %s", url)
	return &AudioStream{cmd: cmd, stdout: stdout}, nil
}

// Read 读取音频数据 (实现 io.Reader 接口)
func (as *AudioStream) Read(p []byte) (n int, err error) {
	return as.stdout.Read(p)
}

// Stop 停止播放并杀掉 FFmpeg 进程
func (as *AudioStream) Stop() {
	if as.cmd != nil && as.cmd.Process != nil {
		_ = as.cmd.Process.Kill()
		log.Println("[FFmpeg] Process killed")
	}
}
