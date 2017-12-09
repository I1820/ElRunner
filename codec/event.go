/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-11-2017
 * |
 * | File Name:     codec/event.go
 * +===============================================
 */

package codec

import "github.com/aiotrc/GoRunner/runner"

// DecodeEvent generates for decode call
type DecodeEvent string

// Data returns data associated with event
func (d DecodeEvent) Data() string {
	return string(d)
}

// Type returns type of event
func (d DecodeEvent) Type() int {
	return runner.UserEventType
}

// EncodeEvent generates for encode call
type EncodeEvent string

// Data returns data associated with event
func (e EncodeEvent) Data() string {
	return string(e)
}

// Type returns type of event
func (e EncodeEvent) Type() int {
	return runner.UserEventType
}
