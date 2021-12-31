package api

const (
	CodeNotFound = 10404
	CodeDbError = 10500
	CodeRequestError  = 10400
	CodeConflict  = 10409

)

func mapCodeToHttp(code int)int{
	switch code {
	case CodeNotFound:
		return 404
	case CodeDbError:
		return 500
	case CodeRequestError:
		return 400
	case CodeConflict:
		return 409
	}
	return 500
}
