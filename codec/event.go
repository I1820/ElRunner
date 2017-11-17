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
type DecodeEvent []byte

func (d DecodeEvent) Data() []byte {
	return d
}

func (d DecodeEvent) Type() int {
	return runner.UserEventType
}

// EncodeEvent generates for encode call
type EncodeEvent string

func (e EncodeEvent) Data() []byte {
	return []byte(e)
}

func (e EncodeEvent) Type() int {
	return runner.UserEventType
}
