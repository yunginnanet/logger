package logger

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("native file", func(t *testing.T) {
		log := New()
		log.C().SetPrefix(t.Name())
		log.C().Info("yeet")
		log.Z().Info().Msg("yeet")

		tmp := t.TempDir()
		tmpF := filepath.Join(tmp, "test.log")
		f, err := os.Create(tmpF)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		log.AddWriter(f)
		log.C().Info("yeet 1")
		log.Z().Info().Msg("yeet 2")
		if err = f.Sync(); err != nil {
			t.Fatalf("failed to sync file: %v", err)
		}
		if err = f.Close(); err != nil {
			t.Fatalf("failed to close file: %v", err)
		}
		f, err = os.Open(tmpF)
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}
		defer func() {
			_ = f.Close()
		}()
		xerox := bufio.NewScanner(f)
		i := 0
		for xerox.Scan() {
			if xerox.Err() != nil {
				t.Fatalf("error occurred during scan: %v", xerox.Err())
			}

			t.Log(xerox.Text())

			if !strings.Contains(xerox.Text(), t.Name()) {
				t.Errorf("missing test name prefix")
			}

			switch i {
			case 0:
				if !strings.Contains(xerox.Text(), "yeet 1") {
					t.Fatalf("unexpected line: %s", xerox.Text())
				}
			case 1:
				if !strings.Contains(xerox.Text(), "yeet 2") {
					t.Fatalf("unexpected line: %s", xerox.Text())
				}
			default:
				if strings.TrimSpace(xerox.Text()) != "" {
					t.Fatal("too many lines in file")
				}
			}

			i++
		}
	})

	t.Run("assisted file", func(t *testing.T) {
		log := New()
		log.C().SetPrefix(t.Name())

		tmpDir := t.TempDir()
		f, err := CreateDatedLogFile(tmpDir, "yeet")
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		lDir, lName := filepath.Split(f.Name())
		lDir = filepath.Dir(lDir)
		if lDir != tmpDir {
			t.Fatalf("unexpected log directory: %s, expected %s", lDir, tmpDir)
		}

		split := strings.Split(lName, "-")
		if split[0] != "yeet" {
			t.Fatalf("incorrect filename: %s", split[0])
		}
		log.AddWriter(f)
		log.C().Info("yeet NOW")
		log.Z().Info().Msg("or now, whatever")
		time.Sleep(10 * time.Millisecond)
		dat, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		xerox := bufio.NewScanner(bytes.NewReader(dat))
		i := 0
		for xerox.Scan() {
			if xerox.Err() != nil {
				t.Fatalf("failed to scan: %v", xerox.Err())
			}
			switch i {
			case 0:
				sought := "yeet NOW"
				if !strings.Contains(xerox.Text(), sought) {
					t.Fatalf("got: %s, missing: %s", xerox.Text(), sought)
				}
			case 1:
				sought := "or now, whatever"
				if !strings.Contains(xerox.Text(), sought) {
					t.Fatalf("got: %s, missing: %s", xerox.Text(), sought)
				}
			default:
				t.Fatal("too many lines in file")
			}
			i++
		}
	})

	t.Run("assisted file with context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		tmpDir := t.TempDir()

		f, err := CreateDatedLogFileCtx(ctx, tmpDir, "ctx")

		t.Run("start periodic sync", func(t *testing.T) {
			StartPeriodicSync(ctx, f, time.Millisecond)
		})
		defer func() {
			t.Log("cancelling context")
			cancel()
			time.Sleep(5 * time.Millisecond)
			t.Run("closed after context closure", func(t *testing.T) {
				if _, err = f.Write([]byte{0x05}); err == nil {
					t.Errorf("write should have failed after context closure")
				}
			})
		}()

		log := New(f)
		log.C().SetPrefix(t.Name())
		log.C().Info("dated yeet")
		log.Z().Info().Msg("dated yeet")

		time.Sleep(10 * time.Millisecond)

		fDir, fName := filepath.Split(f.Name())

		t.Logf("dir: %s - name: %s", fDir, fName)

		if !strings.HasPrefix(fName, "ctx") {
			t.Errorf("bad filename: %s", fName)
		}

		if strings.TrimSuffix(fDir, "/") != strings.TrimSuffix(tmpDir, "/") {
			t.Errorf("bad directory: %s", fDir)
		}

		var dat []byte
		if dat, err = os.ReadFile(f.Name()); err != nil {
			t.Fatalf("failed to read dated file: %v", err)
		}

		datStr := string(dat)
		if !strings.Contains(datStr, "dated yeet") {
			t.Fatalf("missing expected contents in %s", filepath.Join(tmpDir, f.Name()))
		}
	})

	t.Run("test global logger", func(t *testing.T) {
		log := New()
		log.C().SetPrefix(t.Name())
		log.Z().Info().Msg("we out here")
		if Global() != nil {
			t.Fatalf("before we set the global logger it should be nil")
		}
		log.WithGlobalPackageAccess()
		if Global() == nil {
			t.Fatalf("after setting the global logger it should not be nil")
		}
		Global().Z().Info().Msg("yeeterson mcgeeterson")
	})
}
