package utils

const (
	GET_AUTH_CONFIG_FROM_USER_AND_USER_DIRECTORY = `SELECT auth_server.auth_config_json AS config
        FROM user
        JOIN user_directory ON user.user_directory_id = user_directory.id
        JOIN auth_server ON auth_server.user_directory_id = user_directory.id
        WHERE user.id = ?`
)
