package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetTokenToCookie(c *fiber.Ctx, name, token, origin string, duration int) {
	forceProduction := strings.Contains(origin, "localhost")

	cookie := &fiber.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		MaxAge:   duration,
	}

	if forceProduction {
		cleanDomain := extractCleanDomain(origin)
		cookie.Secure = true

		if cleanDomain != "" && cleanDomain != "localhost" && cleanDomain != "127.0.0.1" {
			if strings.Contains(cleanDomain, "udatopup.com") {
				if cleanDomain != "udatopup.com" {
					cookie.Domain = ".udatopup.com"
				}
				cookie.SameSite = "None"
			} else {
				cookie.Domain = cleanDomain
				cookie.SameSite = "Lax"
			}
		} else {
			cookie.SameSite = "Lax"
			cookie.Secure = false
		}
	} else {
		cookie.Secure = false
		cookie.SameSite = "Lax"
	}

	c.Cookie(cookie)
}

func DeleteTokenCookie(c *fiber.Ctx, name, origin string, isProduction bool) {
	forceProduction := isProduction || strings.Contains(origin, "udatopup.com")

	cookie := &fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		MaxAge:   -1,
	}

	if forceProduction {
		cleanDomain := extractCleanDomain(origin)
		cookie.Secure = true

		if cleanDomain != "" && cleanDomain != "localhost" && cleanDomain != "127.0.0.1" {
			if strings.Contains(cleanDomain, "udatopup.com") {
				// PERBAIKAN: Sama seperti SetTokenToCookie
				if cleanDomain != "udatopup.com" {
					cookie.Domain = ".udatopup.com"
				}
				cookie.SameSite = "None"
			} else {
				cookie.Domain = cleanDomain
				cookie.SameSite = "Lax"
			}
		} else {
			cookie.SameSite = "Lax"
			cookie.Secure = false
		}
	} else {
		cookie.Secure = false
		cookie.SameSite = "Lax"
	}

	c.Cookie(cookie)
}

func extractCleanDomain(origin string) string {
	if origin == "" {
		return ""
	}

	domain := strings.TrimPrefix(origin, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	if colonIndex := strings.Index(domain, ":"); colonIndex != -1 {
		domain = domain[:colonIndex]
	}

	if slashIndex := strings.Index(domain, "/"); slashIndex != -1 {
		domain = domain[:slashIndex]
	}

	return domain
}
