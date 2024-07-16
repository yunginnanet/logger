package logger

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const MaximumSyncErrors = 10000

var manageFileSyncErrors bool

// DisableSyncErrorAccounting disables storing and panicking upon exceeding the maximum sync error constant with regard
// to log files created with [CreateDatedLogfile] and managed by [StartPeriodicSync]. By default, this behavior is enabled.
func DisableSyncErrorAccounting() {
	manageFileSyncErrors = false
}

// EnableSyncErrorAccounting enables storing and panicking upon exceeding the maximum sync error constant with regard to log files
// created with [CreateDatedLogfile] and managed by [StartPeriodicSync]. By default, this behavior is enabled.
func EnableSyncErrorAccounting() {
	manageFileSyncErrors = true
}

// CreateDatedLogfile creates a new logfile with the current date appended to the filename.
func createDatedLogfile(ctx context.Context, directory, prefix string, format ...string) (*os.File, error) {
	if len(format) != 0 && len(format) != 1 {
		return nil, errors.New("specify either zero arguments or a single argument for time format")
	}

	var tnow = strconv.Itoa(int(time.Now().UnixMilli()))

	if len(format) == 1 {
		tnow = time.Now().Format(format[0])
		if tnow == format[0] {
			return nil, errors.New("invalid time format")
		}
	}

	fname := prefix + "-" + tnow + ".log"
	path := filepath.Join(directory, fname)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	go func() {
		if ctx == nil {
			return
		}
		<-ctx.Done()
		_ = f.Close()
	}()

	return f, nil
}

func CreateDatedLogFile(directory, prefix string) (*os.File, error) {
	ctx := context.Background()
	return createDatedLogfile(ctx, directory, prefix)
}

func CreateDatedLogFileFormatted(directory, prefix string, format string) (*os.File, error) {
	ctx := context.Background()
	return createDatedLogfile(ctx, directory, prefix, format)
}

func CreateDatedLogFileCtx(ctx context.Context, directory, prefix string) (*os.File, error) {
	return createDatedLogfile(ctx, directory, prefix)
}

func CreateDatedLogFileFormattedCtx(ctx context.Context, directory, prefix string, format string) (*os.File, error) {
	return createDatedLogfile(ctx, directory, prefix, format)
}

func syncFile(ctx context.Context, f *os.File) {
	select {
	case <-ctx.Done():
		return
	default:
		//
	}
	var err error
	if err = f.Sync(); err == nil || !manageFileSyncErrors {
		return
	}
	if ctx.Value("sync_err") == nil {
		ctx = context.WithValue(ctx, "sync_err", make([]error, 0, 1))
	}
	var errs []error
	var ok bool
	if errs, ok = ctx.Value("sync_err").([]error); !ok {
		panic("context sync_err values is not a slice of errors...!? the sky is falling")
	}
	errs = append(errs, err)
	if len(errs) > MaximumSyncErrors && manageFileSyncErrors {
		panic("logger: too many log file sync errors")
	}
}

func StartPeriodicSync(ctx context.Context, f *os.File, dur time.Duration) {
	ticker := time.NewTicker(dur)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				syncFile(ctx, f)
			}
		}
	}()
}
