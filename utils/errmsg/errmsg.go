package errmsg

const (
	SUCCESS = 200
	ERROR   = 500
	NOT_LOGIN = 100001
	ADD_FAILED = 100002
	LOGIN_BY_LOCAL_FAILED = 100003
	UNKNOW_LOGIN_TYPE = 100004
	CREATE_USER_DIRECTORY_FAILED = 100005
	CREATE_AUTH_SERVER_FAIL = 100006
	GET_AUTH_SERVER_FAIL = 100007
	LOGIN_BY_HTTP_FAIL = 100008
	LOGIN_BY_OAUTH2_FAIL = 100009
	GET_AUTH_CONFIG_FAIL = 1000010
)

var codeMsg = map[int]string{
	SUCCESS: "OK",
	ERROR: "系统内部错误",
	NOT_LOGIN: "当前用户未登录",
	ADD_FAILED: "添加用户失败",
}

var ErrMsg = map[string]string{
	"ADD_SUCCESS": "添加用户成功",
	"LOGIN_BY_LOCAL_FAILED": "本地密码登录失败",
	"UNKNOW_LOGIN_TYPE": "未知的登录类型",
	"LOGIN_BY_LOCAL_SUCCESS": "本地密码登录成功",
	"CREATE_USER_DIRECTORY_FAILED": "用户目录创建失败",
	"CREATE_USER_DIRECTORY_SUCCESS": "创建用户目录成功",
	"CREATE_AUTH_SERVER_FAIL": "创建认证服务器失败",
	"CREATE_AUTH_SERVER_SUCCESS": "创建认证服务器成功",
	"GET_AUTH_SERVER_FAIL": "获取认证服务器失败",
	"LOGIN_BY_HTTP_SUCCESS": "HTTP登录成功",
	"LOGIN_BY_HTTP_FAIL": "HTTP登录失败",
	"LOGIN_BY_OAUTH2_SUCCESS": "HTTP登录成功",
	"LOGIN_BY_OAUTH2_FAIL": "HTTP登录失败",
	"GET_AUTH_CONFIG_FAIL": "获取认证服务器配置失败",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}