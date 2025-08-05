package tama

import "github.com/upmaru/tama-go/motor"

type MotorService struct{ *motor.Service }

func newMotorService(c *Client) *MotorService { return &MotorService{motor.NewService(c.httpClient)} }
