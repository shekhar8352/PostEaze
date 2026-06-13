package apiv1

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils/mediasign"
)

func StreamMediaHandler(c *gin.Context) {
	versionID, err := strconv.ParseInt(c.Param("versionId"), 10, 64)
	if err != nil || versionID <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	expStr := c.Query("exp")
	sig := c.Query("sig")
	expUnix, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil || !mediasign.Verify(versionID, expUnix, sig) {
		c.Status(http.StatusForbidden)
		return
	}

	v, ownerID, err := repositories.GetMediaVersionWithOwner(c.Request.Context(), versionID)
	if err != nil || v == nil {
		c.Status(http.StatusNotFound)
		return
	}

	rangeHeader := c.GetHeader("Range")

	if v.StorageProvider == entities.StorageProviderBlob && v.BlobURL != "" {
		proxyHTTPURL(c, v.BlobURL, rangeHeader, v.ContentType)
		return
	}

	body, contentType, contentLength, contentRange, statusCode, err := businessv1.StreamDriveVersion(
		c.Request.Context(), ownerID, v, rangeHeader,
	)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	defer body.Close()

	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", contentType)
	if contentRange != "" {
		c.Header("Content-Range", contentRange)
	}
	if contentLength > 0 {
		c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	}
	c.Status(statusCode)
	_, _ = io.Copy(c.Writer, body)
}

func proxyHTTPURL(c *gin.Context, blobURL, rangeHeader, contentType string) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, blobURL, nil)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = contentType
	}
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", ct)
	if cr := resp.Header.Get("Content-Range"); cr != "" {
		c.Header("Content-Range", cr)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		c.Header("Content-Length", cl)
	}
	c.Status(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)
}
