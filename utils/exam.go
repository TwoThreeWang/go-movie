package utils

import (
	"fmt"
	"math/rand"
	"time"
)

/**
在这个示例代码中，我们首先定义了一个名为callAPI的函数，它接受一个外部接口的名称和一个bool类型的channel。在函数中，我们使用select语句来监听channel，
如果接收到了信号，就直接返回。否则，我们调用外部接口的代码，并在调用完成后向channel发送信号。

在主函数中，我们首先定义了一个bool类型的channel，然后遍历了三个外部接口。对于每个接口，我们使用goroutine来异步调用callAPI函数，并在调用完成后等待3秒。
然后，我们使用select语句来监听channel。如果接收到了信号，说明有一个接口返回数据了，那么我们就输出提示信息并直接返回。否则，我们继续调用下一个外部接口。

如果10秒内没有任何一个接口返回数据，那么我们就输出超时提示信息，并向channel发送信号以停止所有goroutine。
*/

func callAPI(api string, ch chan bool) {
	select {
	case <-ch:
		return
	default:
		// 调用接口的代码
		fmt.Printf("Calling %s...\n", api)
		// 生成一个 0 到 10 秒之间的随机数
		rand.Seed(time.Now().UnixNano()) // 初始化随机数生成器
		waitTime := rand.Intn(10)
		time.Sleep(time.Duration(waitTime) * time.Second)
		fmt.Println(api + " 返回数据了\n")
		ch <- true
	}
}

// 全局的超时监控
func timeOutMouit(ch chan bool) {
	select {
	case <-ch:
		return
	default:
		<-time.After(10 * time.Second) // 等待十秒
		fmt.Println("Timeout!")
		ch <- true // 往通道发送一个超时结束信号
	}
}

func Wb() {
	// 定义了一个bool类型的channel，用来传输结束信号
	ch := make(chan bool)
	// 异步执行一个超时监控，超时后往通道发送一个结束信号
	go timeOutMouit(ch)
	// 定义三个不同的外部接口
	apis := []string{"api1", "api2", "api3", "api4", "api5"}
	// 遍历三个接口
	for _, api := range apis {
		// 协程的方式去调用接口
		go callAPI(api, ch)
		// 等待三分钟
		<-time.After(3 * time.Second)
		select {
		case <-ch:
			// 如果channel变成true了，说明有结果返回了或者超时了，直接return不再继续调用剩余接口
			fmt.Println("Data received!")
			return
		default:
			continue
		}
	}
}
