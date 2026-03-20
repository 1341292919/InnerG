package oss

import (
	"InnerG/config"
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func IsImage(data *multipart.FileHeader) error {
	file, err := data.Open()
	if err != nil {
		return fmt.Errorf("open file error " + err.Error())
	}
	defer file.Close()

	// 读取足够的文件头（512字节通常足够）
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read file error" + err.Error())
	}
	if n < 12 { // 最小需要读取一些字节来检测基本格式
		return fmt.Errorf("file too small" + err.Error())
	}

	// 使用可靠的文件类型检测库
	kind, _ := filetype.Match(buffer)
	if kind == filetype.Unknown {
		// 检查文件扩展名作为后备方案
		ext := strings.ToLower(filepath.Ext(data.Filename))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".svg":
			return nil
		default:
			return fmt.Errorf("file not image")
		}
	}

	if filetype.IsImage(buffer) {
		return nil
	}

	return fmt.Errorf("file not image")
}

func SaveFile(data *multipart.FileHeader, storePath, fileName string) (err error) {
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		err := os.MkdirAll(storePath, 0755) //0755 是一个八进制数，表示文件或目录的权限。它的二进制形式是 111 101 101，对应的权限
		if err != nil {
			return fmt.Errorf("mkdir error " + err.Error())
		}
	}

	//打开本地文件
	dist, err := os.OpenFile(filepath.Join(storePath, fileName), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("open file error " + err.Error())
	}
	defer func(dist *os.File) {
		_ = dist.Close()
	}(dist)

	//打开上传的文件
	src, err := data.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)
	// 复制文件内容
	_, err = io.Copy(dist, src)

	return nil
}

func Upload(localFile, filename, userid, origin string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", origin, userid, filename)

	putPolicy := storage.PutPolicy{
		Scope: config.Oss.Bucket,
	}

	mac := auth.New(config.Oss.AccessKey, config.Oss.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = getQiniuZone(config.Oss.Region)
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}

	recorder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		return "", fmt.Errorf("create file recorder failed" + err.Error())
	}

	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}
	err = resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return "", fmt.Errorf("upload file error：%v", err.Error())
	}
	defer func() {
		err = os.Remove(localFile)
	}()
	if err != nil {
		return "", fmt.Errorf("remove file error")
	}
	return storage.MakePublicURL(config.Oss.Domain, ret.Key), nil
}

func getQiniuZone(region string) *storage.Zone {
	switch region {
	case "z0":
		return &storage.Zone_z0
	case "z1":
		return &storage.Zone_z1
	case "z2":
		return &storage.Zone_z2
	case "na0":
		return &storage.Zone_na0
	case "as0":
		return &storage.Zone_as0
	default:
		return &storage.Zone_z0
	}
}
