package course

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type CourseLength struct {
	mu     sync.Mutex
	length float64

	FolderPath string
}

func NewCourseLength(folderPath string) CourseLength {
	return CourseLength{
		mu:         sync.Mutex{},
		length:     0,
		FolderPath: folderPath,
	}
}

func (courseLen *CourseLength) computeLenght() error {
	var videos []string
	var wg sync.WaitGroup
	err := filepath.Walk(courseLen.FolderPath, func(path string, info os.FileInfo, err error) error {
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

	for _, video := range videos {
		wg.Add(1)
		courseLen.computeVideoLen(video, &wg)
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

func (cl *CourseLength) GetCourseLength() (string, error) {
	err := cl.computeLenght()
	if err != nil {
		return "", err
	}

	return FormatDuration(cl.length), nil
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

func (courseLen *CourseLength) computeVideoLen(dir string, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Printf("video: %s\n", dir)
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", dir)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var result FFProbeOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}

	if result.Format.Duration == "" {
		return err
	}

	durationFloat, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return err
	}

	courseLen.mu.Lock()
	courseLen.length += durationFloat
	courseLen.mu.Unlock()

	return nil
}
