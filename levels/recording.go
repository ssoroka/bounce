package levels

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

var (
	recording  = false
	ffmpeg     *exec.Cmd
	ffmpegPipe io.WriteCloser
	filename   string
)

func (g *Game) StartRecording(width, height int) error {
	filename = findFreeFilename()
	ffmpeg = exec.Command("ffmpeg",
		"-y",             // overwrite output
		"-f", "rawvideo", // input format
		"-pix_fmt", "rgba", // pixel format
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-r", "60", // framerate
		"-i", "pipe:0", // read from stdin
		"-c:v", "libx264", // H.264 codec
		"-pix_fmt", "yuv420p", // output pixel format
		"-preset", "fast",
		filename,
	)

	var err error
	ffmpegPipe, err = ffmpeg.StdinPipe()
	if err != nil {
		return err
	}

	if err := ffmpeg.Start(); err != nil {
		return err
	}
	recording = true
	return nil
}

func (g *Game) StopRecording() string {
	if ffmpegPipe != nil {
		if err := ffmpegPipe.Close(); err != nil {
			fmt.Println("error closing ffmpeg pipe:", err)
		}
		go func() {
			if err := ffmpeg.Wait(); err != nil {
				fmt.Println("ffmpeg error:", err)
			}
		}()
		recording = false
	}
	return filename
}

func findFreeFilename() string {
	for i := 1; ; i++ {
		filename := fmt.Sprintf("recording_%d.mp4", i)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename
		}
	}
}
