package controllers

import (
	"documentManage/database"
	"documentManage/models"
	"documentManage/utils"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !user.CheckPassword(req.OldPassword) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "修改密码失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func GetDocuments(c *gin.Context) {
	var documents []models.Document
	query := database.DB.Model(&models.Document{}).Where("status = 1")

	title := c.Query("title")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	documentTypeID := c.Query("document_type_id")
	if documentTypeID != "" {
		query = query.Where("document_type_id = ?", documentTypeID)
	}

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("keyword LIKE ?", "%"+keyword+"%")
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

func GetDocumentDetail(c *gin.Context) {
	id := c.Param("id")

	var document models.Document
	if err := database.DB.Preload("DocumentType").Preload("Department").First(&document, id).Error; err != nil {
		utils.NotFound(c, "档案不存在")
		return
	}

	if document.Status != 1 {
		utils.NotFound(c, "档案不存在")
		return
	}

	database.DB.Model(&document).Update("view_count", document.ViewCount+1)

	utils.Success(c, document)
}

func GetAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	query := database.DB.Model(&models.Announcement{}).Where("status = 1")

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&announcements)

	utils.Success(c, gin.H{
		"list":      announcements,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetAnnouncementDetail(c *gin.Context) {
	id := c.Param("id")

	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if announcement.Status != 1 {
		utils.NotFound(c, "公告不存在")
		return
	}

	utils.Success(c, announcement)
}

func GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var favorites []models.Favorite
	query := database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Document").Preload("Document.DocumentType").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&favorites)

	utils.Success(c, gin.H{
		"list":      favorites,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		DocumentID uint `json:"document_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingFavorite models.Favorite
	if database.DB.Where("user_id = ? AND document_id = ?", userID, req.DocumentID).First(&existingFavorite).Error == nil {
		utils.BadRequest(c, "已收藏")
		return
	}

	favorite := &models.Favorite{
		UserID:     userID,
		DocumentID: req.DocumentID,
	}

	if err := database.DB.Create(favorite).Error; err != nil {
		utils.InternalServerError(c, "收藏失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	documentID := c.Param("document_id")

	if err := database.DB.Where("user_id = ? AND document_id = ?", userID, documentID).Delete(&models.Favorite{}).Error; err != nil {
		utils.InternalServerError(c, "取消收藏失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func GetBanners(c *gin.Context) {
	var banners []models.Banner
	database.DB.Where("status = 1").Order("sort ASC, created_at DESC").Find(&banners)
	utils.Success(c, banners)
}

func GetNews(c *gin.Context) {
	var news []models.News
	query := database.DB.Model(&models.News{}).Where("status = 1")

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&news)

	utils.Success(c, gin.H{
		"list":      news,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetNewsDetail(c *gin.Context) {
	id := c.Param("id")

	var news models.News
	if err := database.DB.First(&news, id).Error; err != nil {
		utils.NotFound(c, "新闻不存在")
		return
	}

	if news.Status != 1 {
		utils.NotFound(c, "新闻不存在")
		return
	}

	utils.Success(c, news)
}

func GetDocumentTypes(c *gin.Context) {
	var types []models.DocumentType
	database.DB.Where("status = 1").Order("sort ASC, id ASC").Find(&types)
	utils.Success(c, types)
}

func GetDepartments(c *gin.Context) {
	var departments []models.Department
	database.DB.Where("status = 1").Order("sort ASC, id ASC").Find(&departments)
	utils.Success(c, departments)
}

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return defaultVal
		}
	}
	return result
}
