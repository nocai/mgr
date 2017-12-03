package conf

const DEFAULT_LAYOUT = "2006-01-02 15:04:05"

const (
	Page int64 = 1
	Rows int64 = 10
)

const (
	// ResultMsg
	SuccessCode = 0
	SuccessMsg = "操作成功"

	FailCode = -1
	FailMsg = "系统异常"
)

const (
	MsgQuery           = "查询失败"
	MsgInsert          = "添加失败"
	MsgUpdate          = "更新失败"
	MsgDelete          = "删除失败"
	MsgArgument        = "无效参数"
	MsgDataDuplication = "数据重复"
)

