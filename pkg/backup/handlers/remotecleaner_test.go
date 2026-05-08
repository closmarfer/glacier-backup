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

func TestRemoteCleaner_Run(t *testing.T) {
	t.Run("should clean files that do not exist locally", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChecker := backup.NewMockExistentFilesChecker(ctrl)
		mockRepo := backup.NewMockRemoteFilesRepository(ctrl)

		cleaner := NewRemoteCleaner(mockChecker, mockRepo)

		tmpDir, err := ioutil.TempDir("", "glacier-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		existingFile := filepath.Join(tmpDir, "exists.txt")
		err = ioutil.WriteFile(existingFile, []byte("content"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		missingFilePath := filepath.Join(tmpDir, "missing.txt")

		filesMap := map[string]time.Time{
			existingFile:    time.Now(),
			missingFilePath: time.Now(),
		}

		mockChecker.EXPECT().Open(gomock.Any()).Return(nil)
		mockChecker.EXPECT().GetFiles().Return(filesMap)

		mockRepo.EXPECT().Delete(gomock.Any(), missingFilePath).Return(nil)

		mockChecker.EXPECT().Remove(missingFilePath)
		mockChecker.EXPECT().Close(gomock.Any()).Return(nil)

		cleaner.Run()
	})

	t.Run("should not panic if open fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChecker := backup.NewMockExistentFilesChecker(ctrl)
		mockRepo := backup.NewMockRemoteFilesRepository(ctrl)

		cleaner := NewRemoteCleaner(mockChecker, mockRepo)

		mockChecker.EXPECT().Open(gomock.Any()).Return(fmt.Errorf("open error"))

		cleaner.Run()
	})
}
