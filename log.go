package log

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mkideal/log/logger"
	"github.com/mkideal/log/provider"
)

const (
	LvFATAL = logger.FATAL
	LvERROR = logger.ERROR
	LvWARN  = logger.WARN
	LvINFO  = logger.INFO
	LvDEBUG = logger.DEBUG
	LvTRACE = logger.TRACE
)

func ParseLevel(s string) (logger.Level, bool) { return logger.ParseLevel(s) }
func MustParseLevel(s string) logger.Level     { return logger.MustParseLevel(s) }

// global logger
var glogger = logger.NewStdLogger()

func Uninit(err error) {
	glogger.Quit()
}

func InitWithLogger(l logger.Logger) error {
	glogger = l
	glogger.Run()
	return nil
}

func InitWithProvider(p logger.Provider) error {
	glogger = logger.New(p)
	glogger.SetLevel(LvINFO)
	glogger.Run()
	return nil
}

// Init inits global logger with providerType and opts
// * providerType: providerType should be one of {file, console}
// * opts        : opts is a json string or empty
func Init(providerType, opts string) error {
	pcreator := logger.Lookup(providerType)
	if pcreator == nil {
		return errors.New("unsupported provider type: " + providerType)
	}
	return InitWithProvider(pcreator(opts))
}

// InitFile inits with file provider by log file fullpath
func InitFile(fullpath string) error {
	return Init("file", makeFileOpts(fullpath))
}

func makeFileOpts(fullpath string) string {
	dir, filename := filepath.Split(fullpath)
	if dir == "" {
		dir = "."
	}
	return fmt.Sprintf(`{"dir":"%s","filename":"%s"}`, dir, filename)
}

// InitConsole inits with console provider by toStderrLevel
func InitConsole(toStderrLevel logger.Level) error {
	return Init("console", makeConsoleOpts(toStderrLevel))
}

func makeConsoleOpts(toStderrLevel logger.Level) string {
	return fmt.Sprintf(`{"tostderrlevel":%d}`, toStderrLevel)
}

// InitFileAndConsole inits with console and file providers
func InitFileAndConsole(fullpath string, toStderrLevel logger.Level) error {
	fileOpts := makeFileOpts(fullpath)
	consoleOpts := makeConsoleOpts(toStderrLevel)
	p := provider.NewMixProvider(provider.NewFile(fileOpts), provider.NewConsole(consoleOpts))
	return InitWithProvider(p)
}

func GetLevel() logger.Level                   { return glogger.GetLevel() }
func SetLevel(level logger.Level)              { glogger.SetLevel(level) }
func Trace(format string, args ...interface{}) { glogger.Trace(1, format, args...) }
func Debug(format string, args ...interface{}) { glogger.Debug(1, format, args...) }
func Info(format string, args ...interface{})  { glogger.Info(1, format, args...) }
func Warn(format string, args ...interface{})  { glogger.Warn(1, format, args...) }
func Error(format string, args ...interface{}) { glogger.Error(1, format, args...) }
func Fatal(format string, args ...interface{}) { glogger.Fatal(1, format, args...) }