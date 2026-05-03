package routers

import (
	"documentManage/controllers"
	"documentManage/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		api.GET("/banners", controllers.GetBanners)
		api.GET("/news", controllers.GetNews)
		api.GET("/news/:id", controllers.GetNewsDetail)
		api.GET("/announcements", controllers.GetAnnouncements)
		api.GET("/announcements/:id", controllers.GetAnnouncementDetail)
		api.GET("/documents", controllers.GetDocuments)
		api.GET("/documents/:id", controllers.GetDocumentDetail)
		api.GET("/document-types", controllers.GetDocumentTypes)
		api.GET("/departments", controllers.GetDepartments)

		user := api.Group("").Use(middleware.Auth())
		{
			user.GET("/user/info", controllers.GetCurrentUser)
			user.PUT("/user/profile", controllers.UpdateProfile)
			user.PUT("/user/password", controllers.ChangePassword)

			user.GET("/favorites", controllers.GetFavorites)
			user.POST("/favorites", controllers.AddFavorite)
			user.DELETE("/favorites/:document_id", controllers.RemoveFavorite)
		}

		admin := api.Group("/admin").Use(middleware.Auth(), middleware.AdminOrSuperAdminRequired())
		{
			admin.GET("/documents", controllers.AdminGetDocuments)
			admin.GET("/documents/:id", controllers.AdminGetDocumentDetail)
			admin.POST("/documents", controllers.AdminCreateDocument)
			admin.PUT("/documents/:id", controllers.AdminUpdateDocument)
			admin.DELETE("/documents/:id", controllers.AdminDeleteDocument)
		}

		superAdmin := api.Group("/super-admin").Use(middleware.Auth(), middleware.SuperAdminRequired())
		{
			superAdmin.GET("/users", controllers.SuperAdminGetUsers)
			superAdmin.GET("/users/:id", controllers.SuperAdminGetUserDetail)
			superAdmin.POST("/users", controllers.SuperAdminCreateUser)
			superAdmin.PUT("/users/:id", controllers.SuperAdminUpdateUser)
			superAdmin.DELETE("/users/:id", controllers.SuperAdminDeleteUser)
			superAdmin.POST("/users/:id/reset-password", controllers.SuperAdminResetPassword)

			superAdmin.GET("/departments", controllers.SuperAdminGetDepartments)
			superAdmin.POST("/departments", controllers.SuperAdminCreateDepartment)
			superAdmin.PUT("/departments/:id", controllers.SuperAdminUpdateDepartment)
			superAdmin.DELETE("/departments/:id", controllers.SuperAdminDeleteDepartment)

			superAdmin.GET("/document-types", controllers.SuperAdminGetDocumentTypes)
			superAdmin.POST("/document-types", controllers.SuperAdminCreateDocumentType)
			superAdmin.PUT("/document-types/:id", controllers.SuperAdminUpdateDocumentType)
			superAdmin.DELETE("/document-types/:id", controllers.SuperAdminDeleteDocumentType)

			superAdmin.GET("/news", controllers.SuperAdminGetNews)
			superAdmin.GET("/news/:id", controllers.SuperAdminGetNewsDetail)
			superAdmin.POST("/news", controllers.SuperAdminCreateNews)
			superAdmin.PUT("/news/:id", controllers.SuperAdminUpdateNews)
			superAdmin.DELETE("/news/:id", controllers.SuperAdminDeleteNews)

			superAdmin.GET("/banners", controllers.SuperAdminGetBanners)
			superAdmin.POST("/banners", controllers.SuperAdminCreateBanner)
			superAdmin.PUT("/banners/:id", controllers.SuperAdminUpdateBanner)
			superAdmin.DELETE("/banners/:id", controllers.SuperAdminDeleteBanner)

			superAdmin.GET("/announcements", controllers.SuperAdminGetAnnouncements)
			superAdmin.POST("/announcements", controllers.SuperAdminCreateAnnouncement)
			superAdmin.PUT("/announcements/:id", controllers.SuperAdminUpdateAnnouncement)
			superAdmin.DELETE("/announcements/:id", controllers.SuperAdminDeleteAnnouncement)
		}
	}

	return r
}
