package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Versioner struct {
	seen map[string]struct{}
	m    sync.Mutex
}

type ErrCounter struct {
	count int
	m     sync.Mutex
}

func (e *ErrCounter) Increment() {
	defer e.m.Unlock()
	e.m.Lock()
	e.count += 1
}

func ProcessDir(targetDir, dstDir string, maxProcs uint, recursive bool, encCfg EncodeCfg, rsmplCfg ResampleCfg) error {
	tdir, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	// gather list of all files to parse
	var files []string

	// use this to prevent file path name collisions
	v := Versioner{
		seen: make(map[string]struct{}),
		m:    sync.Mutex{},
	}
	err = filepath.WalkDir(tdir, func(path string, d fs.DirEntry, err error) error {
		// needed because errors are filterd through the fn passed to filepath.WalkDir
		if err != nil {
			return err
		}
		// skip directory only if d is a directory and d is not the root directory and recursive is false
		if d.IsDir() && path != tdir && !recursive {
			return fs.SkipDir
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	workerChan := make(chan struct{}, maxProcs)
	var wg sync.WaitGroup
	ec := ErrCounter{
		count: 0,
		m:     sync.Mutex{},
	}
	for _, file := range files {
		// add struct{} into workerChan; this should block if maxProcs is reached
		workerChan <- struct{}{}
		wg.Add(1)
		go func(file string) {
			defer func() {
				// remove a struct{} once done processing file
				<-workerChan
				wg.Done()
			}()
			b, likelySrcFmt, err := GetBytesAndFileTypeLocal(file)
			if err != nil {
				ec.Increment()
				return
			}
			img, err := BytesToImage(b, likelySrcFmt)
			if err != nil {
				ec.Increment()
				return
			}
			if rsmplCfg.IsUsed {
				img = Rescale(img, rsmplCfg)
			}
			dstPath, err := GetDstFilePath("", dstDir, file, false, encCfg.FileType)
			if err != nil {
				ec.Increment()
				return
			}

			// check if file already exists
			var version int
			origDstPath := dstPath
			for {
				dstPath = origDstPath
				// start by checking for conflicts with existing files in the dst directory
				if _, err = os.Stat(dstPath); err == nil {
					fdir := filepath.Dir(dstPath)
					fNameExt := filepath.Base(dstPath)
					fNameExtSlice := strings.Split(fNameExt, ".")
					if len(fNameExt) < 1 {
						ec.Increment()
						return
					}
					for ; err == nil; version++ {
						dstNameVersionExt := fNameExtSlice[0] + "_v" + strconv.Itoa(version+1) + "." + fNameExtSlice[1]
						dstPath = filepath.Join(fdir, dstNameVersionExt)
						_, err = os.Stat(dstPath)
					}
				}
				// do one final check to avoid a race with file writes in other go routines
				v.m.Lock()
				_, collision := v.seen[dstPath]
				if !collision {
					v.seen[dstPath] = struct{}{}
					v.m.Unlock()
					break
				}
				v.m.Unlock()
			}
			err = SaveFile(img, dstPath, encCfg)
			if err != nil {
				ec.Increment()
			}
		}(file)
	}
	wg.Wait()
	if ec.count == 0 {
		return nil
	}
	return errors.New("ignored " + strconv.Itoa(ec.count) + " error(s)")
}
