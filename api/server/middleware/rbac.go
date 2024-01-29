package middleware

/*func Rbac() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if AlwaysAllowPath.Has(path) {
			return
		}

		// 用户ID
		uid, exist := c.Get("userId")
		r := httputils.NewResponse()
		if !exist {
			r.SetCode(http.StatusUnauthorized)
			httputils.SetFailed(c, r, httpstatus.NoPermission)
			return
		}

		method := c.Request.Method
		enforcer := oats.CoreV1.Policy().GetEnforce()
		if enforcer == nil {
			log.Logger.Errorf("failed to get enforce.")
			return
		}
		uidStr := strconv.FormatInt(uid.(int64), 10)
		ok, err := enforcer.Enforce(uidStr, path, method)
		if err != nil {
			r.SetCode(http.StatusInternalServerError)
			httputils.SetFailed(c, r, httpstatus.InnerError)
			c.Abort()
			return
		}
		if !ok {
			r.SetCode(http.StatusUnauthorized)
			httputils.SetFailed(c, r, httpstatus.NoPermission)
			c.Abort()
			return
		}
	}
}
*/
