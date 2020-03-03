package utils

import "fmt"

/**
 * Throw InvalidLogLine when try to parse a line cannot be parsed.
 * It happends sometimes as logs may contain many unformatted lines.
 */
type InvalidLogLine struct{}
func (e *InvalidLogLine) Error() string{
	return "Invalid log line"
}

/**
 * Throw ParseFailed when when parse some string failed.
 * It always implies wrong regulation expression
 */
type ParseFailed struct{
	Expr string
	Str string
}

func (e *ParseFailed) Error() string {
	return fmt.Sprintf("Failed to parse %s use %s", e.Str, e.Expr)
}

/**
 *
 */
type UnknownReg struct{
	Reg string
}
func (e *UnknownReg) Error() string{
	return fmt.Sprintf("Unknown regular expression %s.", e.Reg)
}

type WrongEventType struct{
	eventType string
}
func (e *WrongEventType) Error() string{
	return fmt.Sprintf("Get wrong event type %s.", e.eventType)
}

type UnMatchedSelfPeer struct{
	selfpeer string
	peer string
}
func (e *UnMatchedSelfPeer) Error() string {
	return fmt.Sprintf("Unmatched self peer %s. %s", e.selfpeer, e.peer)
}