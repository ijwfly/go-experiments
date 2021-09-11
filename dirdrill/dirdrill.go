package dirdrill

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
)

const FileSizesChanSize = 1024

func readDirectory(path string) ([]os.DirEntry, error) {
	// TODO: open files via file descriptor pool to avoid problems with directories with many directories
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := file.ReadDir(0)
	if err != nil {
		return nil, err
	}
	_ = file.Close()
	return files, nil
}

func extractFileInfo(dirEntry os.DirEntry) (string, int64) {
	fileName := dirEntry.Name()
	fileInfo, _ := dirEntry.Info()
	return fileName, fileInfo.Size()
}

func prepareSelectCases(resultChannels []chan int64) []reflect.SelectCase {
	cases := make([]reflect.SelectCase, len(resultChannels))
	for i, ch := range resultChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	return cases
}

func drillDirStructure(path string, root *FSDirectory, wg *sync.WaitGroup) chan int64 {
	fileSizesChan := make(chan int64)
	go func() {
		defer close(fileSizesChan)
		defer wg.Done()

		files, err := readDirectory(path)
		if err != nil {
			if !errors.Is(err, fs.ErrPermission) {
				panic(err.Error())
			}
			return
		}

		childResultChannels := make([]chan int64, 0, FileSizesChanSize)

		sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
		for _, dirEntry := range files {
			fileName, fileSize := extractFileInfo(dirEntry)
			if dirEntry.IsDir() {
				result := FSDirectory{
					name: fileName,
				}
				root.AddObject(&result)

				// Start recursive goroutine for directory processing
				newPath := filepath.Join(path, dirEntry.Name())
				resultChan := drillDirStructure(newPath, &result, wg)
				childResultChannels = append(childResultChannels, resultChan)
			} else {
				result := FSFile{
					name: fileName,
					size: fileSize,
				}
				fileSizesChan <- fileSize
				root.AddObject(&result)
			}
		}

		// Process childResultChannels for child directories
		cases := prepareSelectCases(childResultChannels)
		openChannels := len(cases)
		for openChannels > 0 {
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				cases[chosen].Chan = reflect.ValueOf(nil)
				openChannels -= 1
				continue
			}

			fileSize := value.Int()

			root.AddFileSize(fileSize)
			fileSizesChan <- fileSize
		}
	}()
	wg.Add(1)
	return fileSizesChan
}

func GetDirStructure(path string) FSDirectory {
	result := FSDirectory{}
	wg := sync.WaitGroup{}
	resultChan := drillDirStructure(path, &result, &wg)
	for _ = range resultChan {
	}
	wg.Wait()
	return result
}
