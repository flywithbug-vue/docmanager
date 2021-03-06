package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flywithbug/file"
	"github.com/flywithbug/log4go"
	"github.com/nfnt/resize"
	"github.com/pborman/uuid"
	"golang.org/x/image/bmp"

	"doc-manager/web_server/model"

	"github.com/gin-gonic/gin"
)

const localImageFilePath = "./image/"

func uploadImageHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()

	//gin将het/http包的FormFile函数封装到c.Request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	defer file.Close()
	//获取文件名
	ext := filepath.Ext(header.Filename)

	//获取文件的md5值
	data, err := ioutil.ReadAll(file)
	h := md5.New()
	h.Write(data)
	value := h.Sum(nil)
	fileName := hex.EncodeToString(value) + ext

	//文件夹创建管理
	month := time.Now().Format("2006-01")
	localPath := localImageFilePath + month + "/"
	//文件路径
	localFilePath := localPath + fileName
	bExit, err := PathExists(localFilePath)

	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, fmt.Sprintf("system err : %s", err.Error()))
		return
	}
	if bExit {
		log4go.Info("fileExit: %s", fileName)
		avatarPath := fmt.Sprintf("filename=%s&dir=%s", fileName, month)
		aRes.SetResponseDataInfo("imagePath", avatarPath)
		return
	}
	out, err := os.Create(localFilePath)
	if err != nil {
		log4go.Info("创建文件失败：%s", err.Error())
		//判断文件夹是否存在
		bExit, err = PathExists(localPath)
		if err != nil {
			log4go.Info(err.Error())
			aRes.SetErrorInfo(http.StatusInternalServerError, fmt.Sprintf("get folder err : %s", err.Error()))
			return
		}
		//文件夹不存在创建文件夹
		if !bExit {
			err = os.Mkdir(localPath, os.ModePerm)
			if err != nil {
				log4go.Info(err.Error())
				aRes.SetErrorInfo(http.StatusInternalServerError, fmt.Sprintf("make folder err : %s", err.Error()))
				return
			}
		}
		//重新启动out
		out, err = os.Create(localFilePath)
		if err != nil {
			log4go.Info(err.Error())
			aRes.SetErrorInfo(http.StatusInternalServerError, fmt.Sprintf("make file err : %s", err.Error()))
			return
		}
	}
	defer out.Close()

	_, err = out.Write(data)
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, fmt.Sprintf("write file err : %s", err.Error()))
		return
	}
	avatarPath := fmt.Sprintf("filename=%s&dir=%s", fileName, month)
	aRes.SetResponseDataInfo("imagePath", avatarPath)
}

func makeFilePath(ext string) (string, error) {
	month := time.Now().Format("2006-01")
	local := localImageFilePath + month + "/"
	name := uuid.New()
	fileName := name + ext
	//判断文件夹是否存在
	bExit, err := PathExists(local)
	if err != nil {
		return "", err
	}
	if !bExit {
		err = os.Mkdir(local, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return local + fileName, nil
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	//return jpeg.Encode(out, canvas, &jpeg.Options{quality})

	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, canvas)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

func getImageHandler(c *gin.Context) {
	filename := c.Query("filename")
	dir := c.Query("dir")
	size := c.Query("size")
	fileOrigin := localImageFilePath + dir + "/" + filename
	if len(size) == 0 {
		http.ServeFile(c.Writer, c.Request, fileOrigin)
		return
	}
	ext := filepath.Ext(filename)
	if strings.EqualFold(ext, ".gif") {
		http.ServeFile(c.Writer, c.Request, fileOrigin)
		return
	}
	filePath := localImageFilePath + dir + "/" + size + "-" + filename
	if !file.FileExists(filePath) {
		if !file.FileExists(fileOrigin) {
			c.Writer.Write([]byte("Error: Image Not found."))
			return
		}
		fIn, _ := os.Open(fileOrigin)
		//log4go.Info(fileOrigin)
		defer fIn.Close()
		fOut, _ := os.Create(filePath)
		//log4go.Info(filename)
		defer fOut.Close()
		err := scale(fIn, fOut, 120, 120, 100)
		if err != nil {
			log4go.Info(err.Error())
			http.ServeFile(c.Writer, c.Request, fileOrigin)
			return
		}
	}
	http.ServeFile(c.Writer, c.Request, filePath)
}
