package constants

type userIdKey string
type userRoleKey string

type sessionKey string

// UserCtx is the key for the user in the context of the request. don't use string directly and DO NOT MODIFY
const UserIdCtxKey userIdKey = "UserID"
const UserRoleCtxKey userRoleKey = "UserRole"
const SessionCtxKey sessionKey = "Session"
