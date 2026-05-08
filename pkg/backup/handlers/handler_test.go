package handlers

import (
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"go.uber.org/mock/gomock"
)

func TestHandler_Run(t *testing.T) {
	t.Run("should upload files and exit gracefully on signal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChecker := backup.NewMockExistentFilesChecker(ctrl)
		mockBackuper := backup.NewMockBackuper(ctrl)
		mockChecker.EXPECT().Close(gomock.Any()).Return(nil).AnyTimes()
		mockChecker.EXPECT().Uploaded().Return(10).AnyTimes()
		mockChecker.EXPECT().Ignored().Return(2).AnyTimes()

		h := NewHandler(mockChecker, mockBackuper)

		mockBackuper.EXPECT().
			Upload(gomock.Any(), gomock.Any()).
			Return(nil)

		done := make(chan bool)
		go func() {
			h.Run()
			done <- true
		}()

		time.Sleep(100 * time.Millisecond)

		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("handler did not return after SIGINT")
		}
	})

	t.Run("should handle upload error and exit gracefully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChecker := backup.NewMockExistentFilesChecker(ctrl)
		mockBackuper := backup.NewMockBackuper(ctrl)
		mockChecker.EXPECT().Close(gomock.Any()).Return(nil).AnyTimes()
		mockChecker.EXPECT().Uploaded().Return(0).AnyTimes()
		mockChecker.EXPECT().Ignored().Return(0).AnyTimes()

		h := NewHandler(mockChecker, mockBackuper)

		expectedErr := fmt.Errorf("upload error")

		mockBackuper.EXPECT().
			Upload(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		done := make(chan bool)
		go func() {
			h.Run()
			done <- true
		}()

		time.Sleep(100 * time.Millisecond)

		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("handler did not return after SIGINT")
		}
	})
}
