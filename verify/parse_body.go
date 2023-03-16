package verify

import (
	"bufio"
	"strings"

	"github.com/gin-gonic/gin"
)

// split request body with ':' as key value
// then append to context
func ParseBodyMiddleware(ctx *gin.Context) {
	scanner := bufio.NewScanner(ctx.Request.Body)
	for scanner.Scan() {
		arr := strings.SplitN(scanner.Text(), ":", 2)
		if len(arr) == 2 {
			ctx.Set(GenDataKey(arr[0]), arr[1])
		}
	}
}

func GenDataKey(str string) string {
	return "DATAVERSE_" + strings.ReplaceAll(strings.ToUpper(str), " ", "_")
}
