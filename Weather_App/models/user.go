package models

// Preferences holds a User’s weather settings (location, unit, verbosity, forecast)
type Preferences struct {
	Location  string
	Unit      string
	Verbosity string
	Forecast  string
}

// User represents an application user with credentials and Preferences
type User struct {
	UserID      string
	Name        string
	Password    string
	Preferences Preferences
}
