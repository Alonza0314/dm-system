package constant

// log
const (
	CFG_LOG  = "CFG"
	ACC_LOG  = "ACC"
	BCK_LOG  = "BCK"
	PROC_LOG = "PROC"
	CTX_LOG  = "CTX"
	DB_LOG   = "DB"
	CAT_LOG  = "CAT"
	DEV_LOG  = "DEV"
	QRD_LOG  = "QRD"
	API_LOG  = "API"
	PWD_LOG  = "PWD"
	SET_LOG  = "SET"
)

// db
const (
	COLL_ID       = "ID"
	COLL_CATEGORY = "category"
	COLL_ACCOUNT  = "account"

	COLL_CATEGORY_TAG = "category-"

	ID_KEY_CATEGORY = "category"
	ID_KEY_DEVICE   = "device"
)

// device status
const (
	STATUS_IDLE  = "idle"
	STATUS_USING = "using"
)

// password
const (
	PWD_MEMORY      = 19 * 1024 // 19 MiB
	PWD_ITERATIONS  = 2
	PWD_PARALLELISM = 1
	PWD_SALT_LENGTH = 16
	PWD_KEY_LENGTH  = 32
)
