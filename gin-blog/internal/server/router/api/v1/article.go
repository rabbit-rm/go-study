package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"blog/internal/server/code"
	"blog/internal/server/models"

	"github.com/gin-gonic/gin"
	"github.com/rabbit-rm/rabbit/validateKit"
	"gorm.io/gorm"
)

func GetArticle(ctx *gin.Context) {
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
	if models.ExistArticleByID(id) {
		article := models.GetArticle(id)
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": article,
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

func GetArticles(ctx *gin.Context) {
	var tagId uint64
	var state uint64
	var err error

	stateParam := ctx.Query("state")
	if strings.Compare("", stateParam) != 0 {
		if state, err = strconv.ParseUint(stateParam, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
		if err = validateKit.ValidateVar(state, "oneof=0 1"); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "state 必须是0|1",
			})
			return
		}
	}
	tagIdParam := ctx.Query("tagId")
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
	var pageIndex uint64
	pageIndexParam := ctx.Query("page")
	if strings.Compare("", pageIndexParam) != 0 {
		if pageIndex, err = strconv.ParseUint(pageIndexParam, 10, 64); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	if err = validateKit.ValidateVar(pageIndex, "gt=0"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "page index 必须大于0",
		})
		return
	}
	conditionParam := ctx.Query("condition")
	var conditionMap map[string]string
	if strings.Compare("", conditionParam) != 0 {
		if err = json.Unmarshal([]byte(conditionParam), &conditionMap); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  err.Error(),
			})
			return
		}
	}
	list, i := models.GetArticleList(nil, 5, int(pageIndex))
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.Success,
		"msg":  code.GetMsg(code.Success),
		"data": list,
		"size": i,
	})
}

func AddArticle(ctx *gin.Context) {
	var tagId uint64
	var state uint64
	var err error
	tagIdParam := ctx.PostForm("tagId")
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
	title := ctx.PostForm("title")
	if err = validateKit.ValidateVar(title, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "title 不能为空",
		})
		return
	}
	description := ctx.PostForm("description")
	if err = validateKit.ValidateVar(description, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "description 不能为空",
		})
		return
	}
	content := ctx.PostForm("content")
	if err = validateKit.ValidateVar(content, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "content 不能为空",
		})
		return
	}
	var createBy string
	if value, exist := ctx.Get("username"); exist {
		createBy = value.(string)
		if err = validateKit.ValidateVar(createBy, "required"); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "createBy 不能为空",
			})
			return
		}
	}
	if models.ExistTagById(tagId) {
		success := models.AddArticle(models.Article{
			TagID:       uint(tagId),
			Title:       title,
			Description: description,
			Content:     content,
			CreateBy:    createBy,
			ModifiedBy:  createBy,
			State:       uint(state),
		})
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": success,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.ErrorNotExistTag,
		"msg":  code.GetMsg(code.ErrorNotExistTag),
	})
}

func EditArticle(ctx *gin.Context) {
	var id uint64
	var tagId uint64
	var state uint64
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
	tagIdParam := ctx.PostForm("tagId")
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
	title := ctx.PostForm("title")
	if err = validateKit.ValidateVar(title, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "title 不能为空",
		})
		return
	}
	description := ctx.PostForm("description")
	if err = validateKit.ValidateVar(description, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "description 不能为空",
		})
		return
	}
	content := ctx.PostForm("content")
	if err = validateKit.ValidateVar(content, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "content 不能为空",
		})
		return
	}
	var modifiedBy string
	if value, exist := ctx.Get("username"); exist {
		modifiedBy = value.(string)
		if err = validateKit.ValidateVar(modifiedBy, "required"); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Error,
				"msg":  "modifiedBy 不能为空",
			})
			return
		}
	}
	if models.ExistArticleByID(id) {
		success := models.EditArticle(id, models.Article{
			Model: gorm.Model{
				ID: uint(id),
			},
			TagID:       uint(tagId),
			Title:       title,
			Description: description,
			Content:     content,
			ModifiedBy:  modifiedBy,
			State:       uint(state),
		})
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": success,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.ErrorNotExistArticle,
		"msg":  code.GetMsg(code.ErrorNotExistArticle),
	})
}

func DeleteArticle(ctx *gin.Context) {
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
	if models.ExistArticleByID(id) {
		success := models.DeleteArticle(id)
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Success,
			"msg":  code.GetMsg(code.Success),
			"data": success,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.ErrorNotExistArticle,
		"msg":  code.GetMsg(code.ErrorNotExistArticle),
	})
}
