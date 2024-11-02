package main

import (
	
	"errors"
	"fmt"
	"strings"
	"url-shorenet/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type request struct {
  Long_url string `json:"long_url"`
}
	func short(long_url string, id int) string {
	last_part := strings.Split(long_url, "/")
	store := last_part[len(last_part)-1]
	return fmt.Sprintf("%s%d", store, id)
}
func url_shorten (c *gin.Context){
	var req request
	if err:=c.BindJSON(&req);err!=nil{
		c.JSON(400,gin.H{
			"error":err.Error(),
		})
		return
}
var url model.Url
url.Url=req.Long_url
err:=model.DB.Where("url=?",url.Url).First(&url).Error
if err != nil {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = model.DB.Create(&url).Error
		if err != nil {
			return
		}
		url.ShortUrl = short(url.Url, url.ID)
		err = model.DB.Save(&url).Error
		if err != nil {
			return
		}
		c.JSON(200, gin.H{
			"short_url": url.ShortUrl,
		})
	} else {
		fmt.Println(err.Error())
		return
	}
	c.JSON(200, gin.H{
		"short_url": url.ShortUrl})
}


}
func get_long_url(c *gin.Context) {
    var req request
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var url model.Url
    err := model.DB.Where("short_url = ?", req.Long_url).First(&url).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, gin.H{"error": "Short URL not found"})
        } else {
            c.JSON(500, gin.H{"error": err.Error()})
        }
        return
    }


    c.JSON(200, gin.H{
        "long_url": url.Url,
    })
}


func main(){
	model.Init()
	r:=gin.Default()
	db:=model.DB
	
	model.Migarete(db)
	r.POST("/shorten",url_shorten)
	r.GET("/get",get_long_url)
	r.Run()

}