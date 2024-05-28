package models

/*
使用 First 方法代替 Select 方法:
因为你只需要获取单个字段的值,所以可以使用 First 方法,它会自动将结果映射到目标变量上。这样可以减少代码行数,并提高可读性。
示例:
go
复制
var reward int
err := models.DB.Model(&models.User{ID: userId}).Select("reward").First(&reward).Error
使用 pluck 方法:
如果你需要获取多个 User 模型的 reward 字段值,可以使用 pluck 方法。它可以将指定字段的值直接映射到一个切片中。
示例:
go
复制
var rewards []int
err := models.DB.Model(&models.User{}).Where("id = ?", userId).Pluck("reward", &rewards).Error
*/
