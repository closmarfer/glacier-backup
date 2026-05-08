package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"go.uber.org/mock/gomock"
)

func TestSizeCounter_Run(t *testing.T) {
	t.Run("should verify size of existing files", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChecker := backup.NewMockExistentFilesChecker(ctrl)

		counter := NewSizeCounter(mockChecker)

		tmpDir, err := ioutil.TempDir("", "glacier-size-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		file1 := filepath.Join(tmpDir, "file1.txt")
		err = ioutil.WriteFile(file1, []byte("12345"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		file2 := filepath.Join(tmpDir, "file2.txt")
		err = ioutil.WriteFile(file2, []byte("1234567890"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		filesMap := map[string]time.Time{
			file1: time.Now(),
			file2: time.Now(),
		}

		mockChecker.EXPECT().Open(gomock.Any()).Return(nil)
		mockChecker.EXPECT().GetFiles().Return(filesMap)

		counter.Run()
	})

	t.Run("should not panic if open fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChecker := backup.NewMockExistentFilesChecker(ctrl)

		counter := NewSizeCounter(mockChecker)

		mockChecker.EXPECT().Open(gomock.Any()).Return(fmt.Errorf("open error"))

		counter.Run()
	})
}
