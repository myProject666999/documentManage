package controllers

import (
	"documentManage/database"
	"documentManage/models"
	"documentManage/utils"

	"github.com/gin-gonic/gin"
)

type CreateDocumentRequest struct {
	Title          string `json:"title" binding:"required"`
	Content        string `json:"content"`
	DocumentNumber string `json:"document_number"`
	DocumentTypeID uint   `json:"document_type_id"`
	DepartmentID   *uint  `json:"department_id"`
	FileURL        string `json:"file_url"`
	Author         string `json:"author"`
	Keyword        string `json:"keyword"`
}

type UpdateDocumentRequest struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	DocumentNumber string `json:"document_number"`
	DocumentTypeID uint   `json:"document_type_id"`
	DepartmentID   *uint  `json:"department_id"`
	FileURL        string `json:"file_url"`
	Author         string `json:"author"`
	Keyword        string `json:"keyword"`
	Status         *int   `json:"status"`
}

func AdminGetDocuments(c *gin.Context) {
	var documents []models.Document
	query := database.DB.Model(&models.Document{})

	title := c.Query("title")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	documentTypeID := c.Query("document_type_id")
	if documentTypeID != "" {
		query = query.Where("document_type_id = ?", documentTypeID)
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Preload("DocumentType").Preload("Department")

	var total int64
	query.Count(&total)

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&documents)

	utils.Success(c, gin.H{
		"list":      documents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AdminGetDocumentDetail(c *gin.Context) {
	id := c.Param("id")

	var document models.Document
	if err := database.DB.Preload("DocumentType").Preload("Department").First(&document, id).Error; err != nil {
		utils.NotFound(c, "档案不存在")
		return
	}

	utils.Success(c, document)
}

func AdminCreateDocument(c *gin.Context) {
	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	document := &models.Document{
		Title:          req.Title,
		Content:        req.Content,
		DocumentNumber: req.DocumentNumber,
		DocumentTypeID: req.DocumentTypeID,
		DepartmentID:   req.DepartmentID,
		FileURL:        req.FileURL,
		Author:         req.Author,
		Keyword:        req.Keyword,
		Status:         1,
	}

	if err := database.DB.Create(document).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, document)
}

func AdminUpdateDocument(c *gin.Context) {
	id := c.Param("id")

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var document models.Document
	if err := database.DB.First(&document, id).Error; err != nil {
		utils.NotFound(c, "档案不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.DocumentNumber != "" {
		updates["document_number"] = req.DocumentNumber
	}
	if req.DocumentTypeID != 0 {
		updates["document_type_id"] = req.DocumentTypeID
	}
	if req.DepartmentID != nil {
		updates["department_id"] = req.DepartmentID
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.Keyword != "" {
		updates["keyword"] = req.Keyword
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&document).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, document)
}

func AdminDeleteDocument(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Document{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}
