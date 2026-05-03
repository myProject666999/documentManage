package controllers

import (
	"documentManage/database"
	"documentManage/models"
	"documentManage/utils"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RealName     string `json:"real_name"`
	Role         string `json:"role" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
}

type UpdateUserRequest struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RealName     string `json:"real_name"`
	Role         string `json:"role"`
	DepartmentID *uint  `json:"department_id"`
	Status       *int   `json:"status"`
}

func SuperAdminGetUsers(c *gin.Context) {
	var users []models.User
	query := database.DB.Model(&models.User{})

	username := c.Query("username")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	role := c.Query("role")
	if role != "" {
		query = query.Where("role = ?", role)
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Preload("Department")

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
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users)

	utils.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func SuperAdminGetUserDetail(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.Preload("Department").First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func SuperAdminCreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingUser models.User
	if database.DB.Where("username = ?", req.Username).First(&existingUser).Error == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	user := &models.User{
		Username:     req.Username,
		Password:     req.Password,
		Email:        req.Email,
		Phone:        req.Phone,
		RealName:     req.RealName,
		Role:         models.Role(req.Role),
		DepartmentID: req.DepartmentID,
		Status:       1,
	}

	if err := user.HashPassword(); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Create(user).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, user)
}

func SuperAdminUpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
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
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.DepartmentID != nil {
		updates["department_id"] = req.DepartmentID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, user)
}

func SuperAdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	currentUserID := c.GetUint("user_id")

	if currentUserID == 0 {
		utils.BadRequest(c, "无法获取当前用户信息")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if user.ID == currentUserID {
		utils.BadRequest(c, "不能删除自己")
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func SuperAdminResetPassword(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user.Password = req.Password
	if err := user.HashPassword(); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "重置密码失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"`
	Sort     int    `json:"sort"`
}

func SuperAdminGetDepartments(c *gin.Context) {
	var departments []models.Department
	query := database.DB.Model(&models.Department{})

	name := c.Query("name")
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Order("sort ASC, id ASC").Find(&departments)

	utils.Success(c, departments)
}

func SuperAdminCreateDepartment(c *gin.Context) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	department := &models.Department{
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   1,
	}

	if err := database.DB.Create(department).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, department)
}

func SuperAdminUpdateDepartment(c *gin.Context) {
	id := c.Param("id")

	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var department models.Department
	if err := database.DB.First(&department, id).Error; err != nil {
		utils.NotFound(c, "部门不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ParentID != nil {
		updates["parent_id"] = req.ParentID
	}
	updates["sort"] = req.Sort

	if err := database.DB.Model(&department).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, department)
}

func SuperAdminDeleteDepartment(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Department{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreateDocumentTypeRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code"`
	ParentID *uint  `json:"parent_id"`
	Sort     int    `json:"sort"`
}

func SuperAdminGetDocumentTypes(c *gin.Context) {
	var types []models.DocumentType
	query := database.DB.Model(&models.DocumentType{})

	name := c.Query("name")
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Order("sort ASC, id ASC").Find(&types)

	utils.Success(c, types)
}

func SuperAdminCreateDocumentType(c *gin.Context) {
	var req CreateDocumentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	documentType := &models.DocumentType{
		Name:     req.Name,
		Code:     req.Code,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   1,
	}

	if err := database.DB.Create(documentType).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, documentType)
}

func SuperAdminUpdateDocumentType(c *gin.Context) {
	id := c.Param("id")

	var req CreateDocumentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var documentType models.DocumentType
	if err := database.DB.First(&documentType, id).Error; err != nil {
		utils.NotFound(c, "档案类型不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.ParentID != nil {
		updates["parent_id"] = req.ParentID
	}
	updates["sort"] = req.Sort

	if err := database.DB.Model(&documentType).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, documentType)
}

func SuperAdminDeleteDocumentType(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.DocumentType{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreateNewsRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Author  string `json:"author"`
	IsTop   int    `json:"is_top"`
}

func SuperAdminGetNews(c *gin.Context) {
	var news []models.News
	query := database.DB.Model(&models.News{})

	title := c.Query("title")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

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
	query.Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&news)

	utils.Success(c, gin.H{
		"list":      news,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func SuperAdminGetNewsDetail(c *gin.Context) {
	id := c.Param("id")

	var news models.News
	if err := database.DB.First(&news, id).Error; err != nil {
		utils.NotFound(c, "新闻不存在")
		return
	}

	utils.Success(c, news)
}

func SuperAdminCreateNews(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	news := &models.News{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
		IsTop:   req.IsTop,
		Status:  1,
	}

	if err := database.DB.Create(news).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, news)
}

type UpdateNewsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	IsTop   *int   `json:"is_top"`
	Status  *int   `json:"status"`
}

func SuperAdminUpdateNews(c *gin.Context) {
	id := c.Param("id")

	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var news models.News
	if err := database.DB.First(&news, id).Error; err != nil {
		utils.NotFound(c, "新闻不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&news).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, news)
}

func SuperAdminDeleteNews(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.News{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreateBannerRequest struct {
	Title   string `json:"title" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
	LinkURL  string `json:"link_url"`
	Sort     int    `json:"sort"`
}

func SuperAdminGetBanners(c *gin.Context) {
	var banners []models.Banner
	query := database.DB.Model(&models.Banner{})

	title := c.Query("title")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

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
	query.Order("sort ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&banners)

	utils.Success(c, gin.H{
		"list":      banners,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func SuperAdminCreateBanner(c *gin.Context) {
	var req CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	banner := &models.Banner{
		Title:    req.Title,
		ImageURL: req.ImageURL,
		LinkURL:  req.LinkURL,
		Sort:     req.Sort,
		Status:   1,
	}

	if err := database.DB.Create(banner).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, banner)
}

type UpdateBannerRequest struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	Sort     *int   `json:"sort"`
	Status   *int   `json:"status"`
}

func SuperAdminUpdateBanner(c *gin.Context) {
	id := c.Param("id")

	var req UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var banner models.Banner
	if err := database.DB.First(&banner, id).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.LinkURL != "" {
		updates["link_url"] = req.LinkURL
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&banner).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, banner)
}

func SuperAdminDeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Banner{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

type CreateAnnouncementRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	IsTop   int    `json:"is_top"`
}

func SuperAdminGetAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	query := database.DB.Model(&models.Announcement{})

	title := c.Query("title")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

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
	query.Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&announcements)

	utils.Success(c, gin.H{
		"list":      announcements,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func SuperAdminCreateAnnouncement(c *gin.Context) {
	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	announcement := &models.Announcement{
		Title:   req.Title,
		Content: req.Content,
		IsTop:   req.IsTop,
		Status:  1,
	}

	if err := database.DB.Create(announcement).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, announcement)
}

type UpdateAnnouncementRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	IsTop   *int   `json:"is_top"`
	Status  *int   `json:"status"`
}

func SuperAdminUpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&announcement).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, announcement)
}

func SuperAdminDeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Announcement{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}
