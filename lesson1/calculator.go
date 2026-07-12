package main

import "fmt"

/*
实现加减乘除运算
支持多个数字运算
返回结果和错误信息
*/
type Calculator struct {
	history []string
}

func NewCalculator() *Calculator {
	return &Calculator{
		history: []string{},
	}
}

func (c *Calculator) Add(nums ...float64) (float64, error) {
	result := 0.0
	for _, num := range nums {
		result += num
	}
	c.history = append(c.history, "Add")
	return result, nil
}

func (c *Calculator) Subtract(nums ...float64) (float64, error) {
	if len(nums) == 0 {
		return 0, nil
	}
	result := nums[0]
	for _, num := range nums[1:] {
		result -= num
	}
	c.history = append(c.history, "Subtract")
	return result, nil
}

func (c *Calculator) Multiply(nums ...float64) (float64, error) {
	result := 1.0
	for _, num := range nums {
		result *= num
	}
	c.history = append(c.history, "Multiply")
	return result, nil
}

func (c *Calculator) Divide(nums ...float64) (float64, error) {
	if len(nums) == 0 {
		return 0, nil
	}
	result := nums[0]
	for _, num := range nums[1:] {
		if num == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result /= num
	}
	c.history = append(c.history, "Divide")
	return result, nil
}

func (c *Calculator) History() []string {
	return c.history
}

func (c *Calculator) ClearHistory() {
	c.history = []string{}
}

func main() {
	calc := NewCalculator()
	result, _ := calc.Add(1, 2, 3)
	fmt.Println("Add result:", result)
	result, _ = calc.Subtract(10, 5, 2)
	fmt.Println("Subtract result:", result)
	result, _ = calc.Multiply(2, 3, 4)
	fmt.Println("Multiply result:", result)
	result, err := calc.Divide(20, 2, 2)
	if err != nil {
		fmt.Println("Divide error:", err)
	} else {
		fmt.Println("Divide result:", result)
	}
	fmt.Println("History:", calc.History())
	calc.ClearHistory()
	fmt.Println("History after clear:", calc.History())
}
