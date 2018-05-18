/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 19-05-2018
 * |
 * | File Name:     requests.go
 * +===============================================
 */

package main

type codeReq struct {
	ID   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}
