package got

import (
	"testing"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
)

func sum(values []int, resultChan chan int) {
	sum := 0
	for _, value := range values {
		sum += value
	}
	// 将计算结果发送到channel中
	resultChan <- sum
}
func TestT(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	resultChan := make(chan int, 3)
	go sum(values[:len(values) / 2], resultChan)
	go sum(values[len(values) / 2:], resultChan)
	go sum(values[len(values) / 3:], resultChan)

	//close(resultChan)
	for ch := <-resultChan {
		fmt.Println(ch);
	}
	for i := 0; i<3;i++ {
		fmt.Println(<-resultChan)
	}
	//for ch := range (resultChan) {
	//	fmt.Println(ch)
	//}
	//fmt.Println("a")
	//close(resultChan)
	//for sum := range resultChan  {
	//	fmt.Println(sum)
	//}
	//sum1, sum2, sum3 := <-resultChan, <-resultChan, <-resultChan
	//fmt.Println("Result:", sum1, sum2, sum3)
}