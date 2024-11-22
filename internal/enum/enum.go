package enum

type Commands string

const (
	Start            Commands = "/start"
	Help             Commands = "/help"
	AdminReport      Commands = "/adminreport"
	AwaitingPassword Commands = "/awaitingpassword"
	AdminLoggedIn    Commands = "/adminloggedIn"
)

type QueryOperation string

const (
	// QueryOperation
	Equal        QueryOperation = "="
	NoEqual      QueryOperation = "!="
	GreaterThan  QueryOperation = ">"
	GreaterEqual QueryOperation = "=>"
	SmallerThan  QueryOperation = "<"
	SmallerEqual QueryOperation = "=<"
	LIKE         QueryOperation = "LIKE"
	NoLIKE       QueryOperation = "NoLIKE"
)

const (
	// Token
	UnExpectedError string = "UnExpectedError"
	ClaimsNotFound  string = "ClaimsNotFound"
	TokenRequired   string = "TokenRequired"
	TokenExpired    string = "TokenExpired"
	TokenInvalid    string = "TokenInvalid"

	// OTP
	OptExists   string = "OtpExists"
	OtpUsed     string = "OtpUsed"
	OtpNotValid string = "OtpInvalid"

	// User
	EmailExists               string = "EmailExists"
	UsernameExists            string = "UsernameExists"
	PermissionDenied          string = "PermissionDenied"
	UsernameOrPasswordInvalid string = "UsernameOrPasswordInvalid"

	// DB
	RecordNotFound string = "RecordNotFound"
)
