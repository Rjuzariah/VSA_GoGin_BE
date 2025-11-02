package services

type ServiceResponse struct {
	Status  int         // e.g., http.StatusOK, http.StatusBadRequest, etc.
	Data    interface{} // your data or result
	Message string      // human-readable message
	Error   error       // actual Go error (for logging if needed)
}
