package api

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"path"
	"strings"

	"blog/internal/config"
	"blog/internal/logger"
	"blog/internal/server/code"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UploadImage(ctx *gin.Context) {
	_, image, err := ctx.Request.FormFile("image")
	if err != nil {
		logger.L().Error("解析文件失败", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  err.Error(),
		})
		return
	}
	if image != nil {
		imageExt := path.Ext(image.Filename)
		// check ext
		match := false
		for _, t := range config.UploadConf().ImageType {
			if strings.Compare(strings.TrimPrefix(imageExt, "."), t) == 0 {
				match = true
				break
			}
		}
		if !match {
			logger.L().Error("不支持的文件后缀")
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "不支持的文件后缀",
			})
			return
		}
		// check size
		imageSize := image.Size
		maxSize := config.UploadConf().ImageMaxSize << 20
		if int64(maxSize) < imageSize {
			logger.L().Error("不支持的文件大小")
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "不支持的文件大小",
			})
			return
		}
		imageName := strings.TrimSuffix(image.Filename, imageExt)
		hash := md5.New()
		hash.Write([]byte(imageName))
		imageNewName := hex.EncodeToString(hash.Sum(nil)) + imageExt
		imageNewPath := path.Join(config.UploadConf().Path, imageNewName)
		// check exist
		exists := true
		_, err := os.Stat(imageNewPath)
		if err != nil {
			if os.IsNotExist(err) {
				exists = false
			}
		}
		if exists {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "保存文件成功",
				"data": gin.H{
					"url": "/images/" + imageNewName,
				},
			})
			return
		}
		err = ctx.SaveUploadedFile(image, imageNewPath)
		if err != nil {
			logger.L().Error("保存文件失败", zap.Error(err))
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "保存文件失败",
			})
			return
		}
		logger.L().Info("保存文件成功", zap.String("path", imageNewName))
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "保存文件成功",
			"data": gin.H{
				"url": "/images/" + imageNewName,
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.Error,
		"msg":  "文件无效",
	})
}
