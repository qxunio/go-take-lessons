package comm

const (
	Enable              = 1
	Disable             = 0
	Admin               = 1
	Teacher             = 2
	Student             = 3
	TeacherId           = 1235808465281093632
	StudentId           = 1235808540581433344
	BatchFactor         = 500
	TimeFormatTime      = "2006-01-02 15:04:05"
	AuthorizationHeader = "Authorization"
	EventStatusDisable  = 0
	EventStatusEnable   = 1
	EventStatusHistory  = 2
)

const (
	RedisAuthEncryptionKey                        = "AUTH_ENCRYPTION:"
	RedisAuthTokenKey                             = "USER_TOKEN:"
	RedisPreAuthTokenKey                          = "PRE_USER_TOKEN:"
	EffectiveTakeLessonsActivityKey               = "effective_take_lessons_activity_key"
	ContextSessionUserKey                         = "sessionUser"
	RedisStuSelectConfigSubjectListKey            = "STU_SUB_LIST_EID:"
	RedisStuSelectConfigSubjectRemainingPlacesKey = "STU_SUB_RP_EID:"
	RedisStuFocusSelfConfigSubjectIdsKey          = "STU_FOCUS_SUB_IDS:"
)
