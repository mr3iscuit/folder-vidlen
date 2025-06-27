package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   getCommandName(),
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		pth := args[0]

		length, err := GetCourseLength(pth)
		if err != nil {
			return err
		}
		fmt.Printf("%s", length)

		return nil
	},
}

func GetCourseLength(dir string) (string, error) {
	var videos []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".mp4") {
			videos = append(videos, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}

	var lenght float64 = 0

	for _, video := range videos {
		vidLenght, err := GetVideoLen(video)
		lenght += vidLenght
		if err != nil {
			return "", err
		}
	}

	return FormatDuration(lenght), nil
}

func FormatDuration(seconds float64) string {
	totalSeconds := int(math.Floor(seconds))
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d\n", hours, minutes, secs)
}

type FFProbeFormat struct {
	Duration string `json:"duration"`
}

type FFProbeOutput struct {
	Format FFProbeFormat `json:"format"`
}

func GetVideoLen(dir string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", dir)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var result FFProbeOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, err
	}

	if result.Format.Duration == "" {
		return 0, nil
	}

	durationFloat, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, err
	}

	return durationFloat, nil
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func getCommandName() string {
	return os.Args[0]
}

func init() {
}
