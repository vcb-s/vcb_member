package service

// AuthMiddleware 登录检查以及token重签
// func AuthMiddleware(c *gin.Context) {
// 	var j JSONData
// 	tokenstring := c.GetHeader("token")

// 	if tokenstring == "" {
// 		j.Unauthorized(c)
// 		return
// 	}

// 	// c.Set("uid", tokenM.Uid)
// 	// c.Set("phone", phone)
// 	// c.Set("version", version)
// 	// c.Set("clientType", clientType)
// 	c.Next()
// }
