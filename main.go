package main

import (
	"fmt"
	"os"
	"os/exec"
)

func mix(videopath, audioPath, outputpath string) bool {
	args := []string{}
	args = append(args, "-i", videopath, "-i", audioPath, "-filter_complex",
		`[0:a]volume=0.8[a0];[1:a]volume=2[a1];[a0][a1]amix=inputs=2[a]`,
		"-map", "0:v", "-map", `[a]`, "-c:v", "libx264", "-c:a", "aac",
		outputpath, "-y")

	var cmd *exec.Cmd

	cmd = exec.Command("ffmpeg", args...)

	//cmd.Dir = "/usr/local/bin"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Args)
	err := cmd.Run()
	if err != nil {
		panic("混流失败")
		return false
	}
	return true
}
func main() {
	mix("C:\\Users\\17805\\Desktop\\17-20-09.mp4",
		"C:\\Users\\17805\\Desktop\\7f966df47555_tmp.ts", "1.ts")

}
