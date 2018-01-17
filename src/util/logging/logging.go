package logging

import (
	"io"
	"os"
	logging "github.com/sirupsen/logrus"
	//"fmt"
)

// Logger wraps logrus.Logger
type Logger struct {
	*logging.Logger
}

type Entry struct {
	*logging.Entry
}

// Log levels.
const (
	PANIC = logging.PanicLevel
	FATAL = logging.FatalLevel
	ERROR = logging.ErrorLevel
	WARNING = logging.WarnLevel
	INFO = logging.InfoLevel
	DEBUG = logging.DebugLevel
)


// LogConfig logger configurations
type LogConfig struct {
	// Level convertes to level during initialization
	Level string
	// list of all modules
	Modules []string
	// format
	Format logging.Formatter
	// output
	Output io.Writer
}



// TODO:
// DefaultLogConfig vs (DevLogConfig + ProdLogConfig) ?

// DevLogConfig default development config for logging
func DevLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		Level:   "DEBUG",
		//Level:   logging.DebugLevel, // string
		Modules: modules,
		Format:  new(logging.TextFormatter),
		//Colors:  true,
		Output:  os.Stdout,
	}
}

// ProdLogConfig Default production config for logging
func ProdLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		Level:   "ERROR",
		Modules: modules,
		Format:  new(logging.TextFormatter),
		Output:  os.Stdout,
	}
}


// InitLogger initialize logging using this LogConfig;
// it panics if l.Format is invalid or l.Level is invalid
func (l *LogConfig) InitLogger() {

	logging.SetFormatter(l.Format)

	level,_ := logging.ParseLevel(l.Level)

	logging.SetLevel(level)
	logging.SetOutput(l.Output)

	//fileHook := NewLogrusWriterHook("log1.txt")
	//
	//logging.AddHook(fileHook)



}

// MustGetLogger safe initialize global logger
func MustGetLoggera(module string) *Logger {
	return &Logger{logging.New()}
}

// MustGetLogger safe initialize global logger
func MustGetLogger(module string) *Entry {

	entry := logging.WithField("module", module)

	return &Entry{entry}
}

// Disable disables the logger completely
func Disable() {
	//logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
}


// hook
//
//type LogrusWriterHook struct {
//	writter       io.Writer
//	formatter     *logging.TextFormatter
//}
//
//func NewLogrusWriterHook(writter io.Writer, disableColors bool) (*LogrusWriterHook) {
//	plainFormatter := &logging.TextFormatter{DisableColors: disableColors}
//
//	return &LogrusWriterHook{&writter, plainFormatter}
//}
//
//// Fire event
//func (hook *LogrusWriterHook) Fire(entry *logging.Entry) error {
//
//	plainformat, err := hook.formatter.Format(entry)
//	line := string(plainformat)
//	_, err = hook.writter.Write([]byte(line))
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "unable to write on writerhook(entry.String)%v", err)
//		return err
//	}
//
//	return nil
//}
