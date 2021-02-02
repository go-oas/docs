package docs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	defaultRoute     = "/api"
	defaultDirectory = "./internal/dist"
	defaultIndexPath = "/index.html"
	fwSlashSuffix    = "/"
)

type ConfigSwaggerUI struct {
	Route string
	Port  string

	httpServer *http.Server
	initFS     FileSystem
	stopper    chan os.Signal
}

// ServeSwaggerUI does what its name implies - runs Swagger UI on mentioned set port and route.
func ServeSwaggerUI(conf *ConfigSwaggerUI) error {
	if conf == nil {
		return errors.New("swagger config is required")
	}

	if conf.Route == "" {
		conf.Route = defaultRoute
	}

	if conf.initFS.isNil() {
		conf.initFS = *initializeStandardFS()
	}

	if conf.httpServer == nil {
		conf.initializeDefaultHTTPServer()
	}

	log.Printf("Serving SwaggerIU on HTTP port: %s\n", conf.Port)
	conf.sigCont()

	select {
	case val := <-conf.stopper:
		if shouldStop(val) {
			err := conf.httpServer.Shutdown(context.Background())
			if err != nil {
				return fmt.Errorf("an error occurred while shutting down SwaggerUI: %w", err)
			}
		} else {
			err := conf.httpServer.ListenAndServe()
			if err != nil {
				return fmt.Errorf("an error occurred while serving SwaggerUI: %w", err)
			}
		}

	}

	return nil
}

// FileSystem represents a wrapper for http.FileSystem, with relevant type func implementations.
type FileSystem struct {
	fileSysInit http.FileSystem

	fsOpenFn  fsOpenFn
	getStatFn getStatFn
	getIsDir  getIsDirFn
}

func initializeStandardFS() *FileSystem {
	fsInit := http.Dir(defaultDirectory)

	return &FileSystem{
		fileSysInit: fsInit,
		fsOpenFn:    newFSOpen(fsInit),
		getStatFn:   newGetStatFn(),
		getIsDir:    newGetIsDirFn(),
	}
}

func (c *ConfigSwaggerUI) initializeDefaultHTTPServer() {
	fileServer := http.FileServer(c.initFS)

	c.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", c.Port),
		Handler: http.StripPrefix(strings.TrimRight(c.Route, fwSlashSuffix), fileServer),
	}
}

func (c *ConfigSwaggerUI) sigCont() {
	if c.stopper == nil {
		var osSignal = make(chan os.Signal)
		c.stopper = osSignal

		go func() {
			time.Sleep(20 * time.Millisecond)
			osSignal <- syscall.SIGCONT
		}()
	}
}

func shouldStop(val os.Signal) bool {
	if val == syscall.SIGINT || val == syscall.SIGKILL {
		return true
	}

	return false
}

func newFSOpen(fis http.FileSystem) fsOpenFn {
	return func(name string) (http.File, error) {
		return fis.Open(name)
	}
}

func newGetStatFn() getStatFn {
	return func(file http.File) fileStatFn {
		return func() (fs.FileInfo, error) {
			return file.Stat()
		}
	}
}

func newGetIsDirFn() getIsDirFn {
	return func(file fs.FileInfo) fsIsDirFn {
		return func() bool {
			return file.IsDir()
		}
	}
}

func (fis *FileSystem) isNil() bool {
	if fis == nil {
		return true
	}

	if fis.getStatFn == nil ||
		fis.getIsDir == nil ||
		fis.fsOpenFn == nil ||
		fis.fileSysInit == nil {
		return true
	}

	return false
}

// Open opens file. Returns http.File, and error if there is any.
func (fis FileSystem) Open(path string) (http.File, error) {
	f, err := fis.fsOpenFn(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file in path %s :%w", path, err)
	}

	fileInfo, err := fis.getStatFn(f)()
	if err != nil {
		return f, fmt.Errorf("failed to fetch file info :%w", err)
	}

	if fis.getIsDir(fileInfo)() {
		index := strings.TrimSuffix(path, fwSlashSuffix) + defaultIndexPath
		if _, err = fis.fileSysInit.Open(index); err != nil {
			return nil, fmt.Errorf("failed trimming path sufix :%w", err)
		}
	}

	return f, nil
}

type (
	fsOpenFn   func(name string) (http.File, error)
	fsIsDirFn  func() bool
	fileStatFn func() (fs.FileInfo, error)
	getStatFn  func(file http.File) fileStatFn
	getIsDirFn func(file fs.FileInfo) fsIsDirFn
)
