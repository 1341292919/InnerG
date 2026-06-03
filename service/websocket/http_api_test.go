package websocket

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
)

func newUploadHeader(t *testing.T, fileName string) *multipart.FileHeader {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write([]byte("fake-file-content-long-enough")); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(body.Len()) + 1024); err != nil {
		t.Fatalf("parse multipart form: %v", err)
	}

	files := req.MultipartForm.File["file"]
	if len(files) != 1 {
		t.Fatalf("expected one uploaded file, got %d", len(files))
	}
	return files[0]

}

func TestUploadImageUsesIMBucketAndImageOrigin(t *testing.T) {
	originalValidateImageFile := validateImageFile
	originalSaveUploadFile := saveUploadFile
	originalUploadLocalFile := uploadLocalFile
	originalNowUnix := nowUnix
	originalGetIMBucket := getIMBucket
	t.Cleanup(func() {
		validateImageFile = originalValidateImageFile
		saveUploadFile = originalSaveUploadFile
		uploadLocalFile = originalUploadLocalFile
		nowUnix = originalNowUnix
		getIMBucket = originalGetIMBucket
	})

	getIMBucket = func() string { return "im-bucket" }
	validateImageFile = func(*multipart.FileHeader) error { return nil }
	saveUploadFile = func(*multipart.FileHeader, string, string) error { return nil }
	nowUnix = func() int64 { return 123 }

	var gotBucket, gotOrigin, gotFilename string
	uploadLocalFile = func(localFile, filename, userid, origin, bucket string) (string, error) {
		gotBucket = bucket
		gotOrigin = origin
		gotFilename = filename
		return "https://cdn.example/ws-image", nil
	}

	ctx := ctl.NewContext(context.Background(), &ctl.UserInfo{Id: 7})
	_, err := (&WebSocketSrv{}).UploadImage(ctx, newUploadHeader(t, "image.png"))
	if err != nil {
		t.Fatalf("upload image returned error: %v", err)
	}

	if gotBucket != "im-bucket" {
		t.Fatalf("expected im bucket, got %q", gotBucket)
	}
	if gotOrigin != constants.WebsocketImageOssOrigin {
		t.Fatalf("expected image origin %q, got %q", constants.WebsocketImageOssOrigin, gotOrigin)
	}
	if gotFilename != "ws_image_7_123" {
		t.Fatalf("expected image filename %q, got %q", "ws_image_7_123", gotFilename)
	}
}

func TestUploadVideoUsesIMBucketAndVideoOrigin(t *testing.T) {
	originalValidateVideoFile := validateVideoFile
	originalSaveUploadFile := saveUploadFile
	originalUploadLocalFile := uploadLocalFile
	originalNowUnix := nowUnix
	originalGetIMBucket := getIMBucket
	t.Cleanup(func() {
		validateVideoFile = originalValidateVideoFile
		saveUploadFile = originalSaveUploadFile
		uploadLocalFile = originalUploadLocalFile
		nowUnix = originalNowUnix
		getIMBucket = originalGetIMBucket
	})

	getIMBucket = func() string { return "im-bucket" }
	validateVideoFile = func(*multipart.FileHeader) error { return nil }
	saveUploadFile = func(*multipart.FileHeader, string, string) error { return nil }
	nowUnix = func() int64 { return 456 }

	var gotBucket, gotOrigin, gotFilename string
	uploadLocalFile = func(localFile, filename, userid, origin, bucket string) (string, error) {
		gotBucket = bucket
		gotOrigin = origin
		gotFilename = filename
		return "https://cdn.example/ws-video", nil
	}

	ctx := ctl.NewContext(context.Background(), &ctl.UserInfo{Id: 7})
	_, err := (&WebSocketSrv{}).UploadVideo(ctx, newUploadHeader(t, "video.mp4"))
	if err != nil {
		t.Fatalf("upload video returned error: %v", err)
	}

	if gotBucket != "im-bucket" {
		t.Fatalf("expected im bucket, got %q", gotBucket)
	}
	if gotOrigin != constants.WebsocketVideoOssOrigin {
		t.Fatalf("expected video origin %q, got %q", constants.WebsocketVideoOssOrigin, gotOrigin)
	}
	if gotFilename != "ws_video_7_456" {
		t.Fatalf("expected video filename %q, got %q", "ws_video_7_456", gotFilename)
	}
}
