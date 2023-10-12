package main

import "errors"

var ErrQueueEmpty = errors.New("queue is empty")
var ErrOnePerUser = errors.New("user already applied task")
var ErrQueueLimit = errors.New("reached queue limit")
var ErrNoUserTask = errors.New("user task not found")