/*
 * @Author: zhaobo
 * @Date: 2020-07-24 14:35:59
 * @Last Modified by: zhaobo
 * @Last Modified time: 2020-07-27 17:29:09
 */

package filesystem

import (
	"filesystem/config"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

// UPloadForm validte
type UPloadForm struct {
	File           *multipart.FileHeader `form:"file" binding:"required"`
	RemoteFilePath string                `form:"remoteFilePath" binding:"required"`
}

// DeleteForm validte
type DeleteForm struct {
	RemoteFilePath string `form:"remoteFilePath" binding:"required"`
}

// CopyForm validte
type CopyForm struct {
	SrcFilePath string `form:"srcFilePath" binding:"required"`
	DstFilePath string `form:"dstFilePath" binding:"required"`
}

// Upload godoc
// @Summary file system upload
// @Description file system upload
// @ID filesystem-upload
// @Tags filesystem
// @Accept  json
// @Produce  json
// @Param  file formData file true "file data"
// @Param  remoteFilePath body string true "remote file path"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /upload [post]
func Upload(c *gin.Context) {
	// 单个文件上传
	var form UPloadForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	remoteFilePath := c.PostForm("remoteFilePath")
	str, _ := os.Getwd()
	dst := str + form.File.Filename
	err = c.SaveUploadedFile(form.File, dst)

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "unknown error")
		return
	}

	// 上传到云端
	fsStorage, err := config.LoadStorage() // 加载存储端配置
	if err != nil {
		fmt.Println("Load storage error:", err)
	}

	err = fsStorage.UploadLocalFile(dst, remoteFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(200, gin.H{
		"status":  form.File.Filename,
		"message": "success",
	})
}

// Delete godoc
// @Summary file system delete
// @Description file system delete
// @ID filesystem-delete
// @Tags filesystem
// @Accept  json
// @Produce  json
// @Param  remoteFilePath body string true "remote file path"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /delete [post]
func Delete(c *gin.Context) {
	var form DeleteForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	fsStorage, err := config.LoadStorage() // 加载存储端配置
	if err != nil {
		fmt.Println("Load storage error:", err)
	}
	err = fsStorage.DeleteSingleFile(form.RemoteFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(200, gin.H{
		"status":  10000,
		"message": "success",
	})
}

// Copy godoc
// @Summary file system copy
// @Description file system copy
// @ID filesystem-copy
// @Tags filesystem
// @Accept  json
// @Produce  json
// @Param  srcFilePath body string true "remote src file path"
// @Param  dstFilePath body string true "remote dst file path"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /copy [post]
func Copy(c *gin.Context) {
	var form CopyForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	fsStorage, err := config.LoadStorage() // 加载存储端配置
	if err != nil {
		fmt.Println("Load storage error:", err)
	}
	err = fsStorage.CopyFile(form.SrcFilePath, form.DstFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(200, gin.H{
		"status":  10000,
		"message": "success",
	})
}

// UploadPolicyForm validte
type UploadPolicyForm struct {
	RemoteFilePath string `json:"remoteFilePath" form:"remoteFilePath" binding:"required"`
	CallbackURL    string `json:"callbackURL" form:"callbackURL" binding:"required"`
	CallbackBody   string `json:"callbackBody" form:"callbackBody" binding:"required"`
}

// UoloadPolicy godoc
// @Summary file system upload policy
// @Description file system get upload policy
// @ID filesystem-upload-policy
// @Tags filesystem
// @Accept  json
// @Produce  json
// @Param  remoteFilePath body string true "remote src file path"
// @Param  callbackURL body string true "upload callback url"
// @Param  callbackBody body string true "callback body"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /upload_policy [post]
func UoloadPolicy(c *gin.Context) {
	var form UploadPolicyForm
	// body, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println(string(body))

	err := c.ShouldBind(&form)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	fsStorage, err := config.LoadStorage() // 加载存储端配置
	if err != nil {
		fmt.Println("Load storage error:", err)
	}
	// body, _ := json.Marshal(form.CallbackBody)
	policyString, err := fsStorage.GetUploadPolicy(form.RemoteFilePath, form.CallbackURL, form.CallbackBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, policyString)
	// c.JSON(200, gin.H{
	// 	"message": policyString,
	// })
}
