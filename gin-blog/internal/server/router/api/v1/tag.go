package v1

import (
	"net/http"
	"strconv"
	"strings"

	"blog/internal/server/code"
	"blog/internal/server/models"

	"github.com/gin-gonic/gin"
	"github.com/rabbit-rm/rabbit/validateKit"
	"gorm.io/gorm"
)

func GetTag(ctx *gin.Context) {
	var id uint64
	var err error
	param := ctx.Param("id")
	if strings.Compare(param, "") != 0 {
		if id, err = strconv.ParseUint(param, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	if err := validateKit.ValidateVar(id, "gt=0"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "id 必须大于0",
		})
		return
	}
	if models.ExistTagById(id) {
		tag, err := models.GetTag(id)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": tag,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.ErrorNotExistArticle,
			"msg":  code.GetMsg(code.ErrorNotExistArticle),
		})
		return
	}
}

func AddTag(ctx *gin.Context) {
	var err error
	tagName := ctx.PostForm("name")
	if err = validateKit.ValidateVar(tagName, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "tag name 不能为空",
		})
		return
	}
	var state uint64
	stateParam := ctx.PostForm("state")
	if strings.Compare("", stateParam) != 0 {
		if state, err = strconv.ParseUint(stateParam, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	if err = validateKit.ValidateVar(state, "oneof=0 1"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "state 必须是0|1",
		})
		return
	}
	var createBy string
	if value, exists := ctx.Get("username"); exists {
		if err = validateKit.ValidateVar(value.(string), "required"); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "createBy 不能为空",
			})
			return
		}
		createBy = value.(string)
	}
	if err = models.AddTag(&models.Tag{
		Name:       tagName,
		CreateBy:   createBy,
		ModifiedBy: createBy,
		State:      1,
	}); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.Success,
		"msg":  code.GetMsg(code.Success),
		"data": true,
	})
}

func EditTag(ctx *gin.Context) {
	var tagId uint64
	var state uint64
	var err error
	tagIdParam := ctx.Param("id")
	if strings.Compare("", tagIdParam) != 0 {
		if tagId, err = strconv.ParseUint(tagIdParam, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	if err = validateKit.ValidateVar(tagId, "gt=0"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "tag id 必须大于0",
		})
		return
	}
	stateParam := ctx.PostForm("state")
	if strings.Compare("", stateParam) != 0 {
		if state, err = strconv.ParseUint(stateParam, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	if err = validateKit.ValidateVar(state, "oneof=0 1"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "state 必须是0|1",
		})
		return
	}
	name := ctx.PostForm("name")
	if err = validateKit.ValidateVar(name, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "name 不能为空",
		})
		return
	}
	var modifiedBy string
	if value, exists := ctx.Get("username"); exists {
		if err = validateKit.ValidateVar(value.(string), "required"); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "modifiedBy 不能为空",
			})
			return
		}
		modifiedBy = value.(string)
	}
	if models.ExistTagById(tagId) {
		if err = models.EditTag(tagId, &models.Tag{
			Model: gorm.Model{
				ID: uint(tagId),
			},
			Name:       name,
			ModifiedBy: modifiedBy,
			State:      uint(state),
		}); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": true,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.ErrorNotExistTag,
		"msg":  code.GetMsg(code.ErrorNotExistTag),
	})
}

func DeleteTag(ctx *gin.Context) {
	var id uint64
	var err error
	idParam := ctx.Param("id")
	if strings.Compare("", idParam) != 0 {
		if id, err = strconv.ParseUint(idParam, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	if err = validateKit.ValidateVar(id, "gt=0"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "id 必须大于0",
		})
		return
	}
	if models.ExistTagById(id) {
		if err = models.DeleteTag(id); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": true,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.ErrorNotExistArticle,
		"msg":  code.GetMsg(code.ErrorNotExistArticle),
	})
}
