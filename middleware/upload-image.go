package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("uploadImage")

		if err != nil {
			fmt.Println("err1=", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		src, err := file.Open();

		if err != nil {
			fmt.Println("err2=", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		defer src.Close()

		tempFile, err := ioutil.TempFile("upload", "image-*.png")

		if err != nil {
			fmt.Println("err3=", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		defer tempFile.Close();

		if _, err = io.Copy(tempFile, src); err != nil {	
			fmt.Println("err4=",err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		data := tempFile.Name();
		fileName := data[7:]

		c.Set("dataFile", fileName)

		return next(c)
	}
}