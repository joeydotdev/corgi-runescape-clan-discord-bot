package plugins

import "errors"

var TooFewArgumentsError error = errors.New("Too few arguments")
var InvalidOperationError error = errors.New("Invalid operation. Valid operations are: add, remove, update")
var NoDiscordUsernameAndDiscriminatorError error = errors.New("No Discord username and discriminator provided.")
var ActiveOngoingEventError error = errors.New("An event is already active. Please stop the current event before starting a new one.")
var NoEventError error = errors.New("No event is currently active. Please start an event before trying to stop it.")
