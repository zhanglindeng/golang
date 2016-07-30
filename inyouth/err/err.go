package err

import "errors"

var (
    ErrEmail = errors.New("邮箱格式错误")
)