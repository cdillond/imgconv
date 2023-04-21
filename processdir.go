package main

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Versioner struct {
	seen map[string]int
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
	//errs := make([]error, 0)
	tdir, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	// gather list of all files to parse
	var files []string

	v := Versioner{
		seen: make(map[string]int),
		m:    sync.Mutex{},
	}
	err = filepath.WalkDir(tdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// skip directory only if d is a directory and d is not the root directory and recursive is false
		if d.IsDir() && path != tdir && !recursive {
			return fs.SkipDir
		}
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
			// this is done synchronously, so no need to bother with mutex
			// this will protect against overwriting already existing files
			v.seen[path] = 0
		}
		return nil
	})

	workerChan := make(chan struct{}, maxProcs)
	var wg sync.WaitGroup
	ec := ErrCounter{
		count: 0,
		m:     sync.Mutex{},
	}
	for _, file := range files {
		// add struct{} into workerChan; this should block if maxProcs is reached
		workerChan <- struct{}{}
		go func(file string) {
			wg.Add(1)
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
			// this part will have to be serialized
			v.m.Lock()
			version, collision := v.seen[dstPath]
			if collision {
				base := filepath.Base(dstPath)
				dir := filepath.Dir(dstPath)
				baseNameExt := strings.Split(base, ".")
				if len(baseNameExt) < 1 {
					v.m.Unlock()
					ec.Increment()
					return
				}
				v.seen[dstPath] = version + 1
				dstPath = filepath.Join(dir, baseNameExt[0]+"_v"+strconv.Itoa(version+1)+"."+baseNameExt[1])
			} else {
				v.seen[dstPath] = 0
			}
			v.m.Unlock()
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
