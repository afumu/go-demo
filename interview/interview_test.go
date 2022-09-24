package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
type U struct {
	Name string
}

// 断言的使用
func Test1(t *testing.T) {
	//type MyInt int
	//var i int = 1
	//var j MyInt = i.(MyInt)

	var i interface{}
	i = U{Name: "zhangsan"}
	//s := i.(string)
	s, ok := i.(string)
	fmt.Println(ok)
	fmt.Println(s)

	switch value := i.(type) {
	case string:
		fmt.Println(value)
	case U:
		fmt.Println(value)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 嵌套函数
func O(a int) func(int) int {
	return func(b int) int {
		return a + b
	}
}

func Test2(t *testing.T) {
	i := O(1)(2)
	fmt.Println(i)
}

// ---------------------------------------------------------------------------------------------------------------------
// 函子的使用,类似于Java的lambda表达式
type Functor interface {
	Fmap(func(int) int) Functor
}

type FunctorImpl struct {
	Arr []int
}

func (f FunctorImpl) Fmap(fn func(int) int) Functor {
	newArr := make([]int, len(f.Arr))
	for i, v := range f.Arr {
		r := fn(v)
		newArr[i] = r
	}
	return FunctorImpl{newArr}
}

func fun1(a int) int {
	return a + 1
}

func fun2(a int) int {
	return a + 1
}

func Test3(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6}
	functorImpl := FunctorImpl{Arr: arr}
	newFunctor1 := functorImpl.Fmap(fun1)
	newFunctor2 := newFunctor1.Fmap(fun2)
	newFunctor12 := functorImpl.Fmap(fun1).Fmap(fun2)
	fmt.Println("原来的值：", arr)
	fmt.Printf("fun1后：%+v\n", newFunctor1)
	fmt.Printf("fun2后：%+v\n", newFunctor2)
	fmt.Printf("fun1,fun2后：%+v\n", newFunctor12)
}

// ---------------------------------------------------------------------------------------------------------------------
// 测试接口赋值和值的传递问题
type Z struct {
	Name string
}

func Test4(t *testing.T) {
	// go 里面一些的的传递都是值传递，没有引用传递。
	var i interface{}
	u := U{"zhangsan"}
	i = &u
	u2 := i.(*U)
	u2.Name = "lisi"
	fmt.Println(u)
	fmt.Println(u2)
}

// ---------------------------------------------------------------------------------------------------------------------
// 判断两个结构体的类型是否一致
// 判断两个接口类型是否一致，需要判断他们的类型和他们对于的数据是否一致

func Test5(t *testing.T) {
	type A struct {
		Name string
	}
	var a A
	a.Name = "123"
	var i interface{}
	i = a
	var i2 interface{}
	a.Name = "126"
	i2 = a
	fmt.Println(i == i2)
}

// ---------------------------------------------------------------------------------------------------------------------
// 继承
type Person struct {
	Name string
	age  int
}

type Student struct {
	*Person
}

func Test6(t *testing.T) {
	var s Student
	s.Person = &Person{}
	s.Name = "123"
	fmt.Println(s)
}

// ---------------------------------------------------------------------------------------------------------------------
// 继承的实现

type IAnimal interface {
	Eat()
}

type Animal struct {
}

func (animal *Animal) Eat() {
	fmt.Println("animal eat")
}

func NewAnimal() *Animal {
	return &Animal{}
}

type Cat struct {
	*Animal
}

func (cat Cat) Eat() {
	fmt.Println("cat eat")
}

func NewCat() *Cat {
	return &Cat{Animal: NewAnimal()}
}

func Test7(t *testing.T) {
	cat := NewCat()
	cat.Eat()
}

// ---------------------------------------------------------------------------------------------------------------------
// 测试
// receiver类型什么时候使用T和*T
// 1.从性能上来说
// 2.从方法的超集来说
type PMall interface {
	Get1()
	Get2()
}

type Mall struct {
	Name string
}

func (m *Mall) Get1() {
	m.Name = "lishi"
	fmt.Println("mall")
}

func (m *Mall) Get2() {
	fmt.Println("mall")
}

func DumpMethodSet(i interface{}) {
	v := reflect.TypeOf(i)
	elemTyp := v.Elem()
	n := elemTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", elemTyp)
		return
	}
	fmt.Printf("%s's method set:\n", elemTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", elemTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}

func Test8(t *testing.T) {
	var pmall PMall
	var m1 Mall
	//pmall = m1
	var m2 *Mall
	pmall = m2
	DumpMethodSet(&m1)
	DumpMethodSet(&m2)
	DumpMethodSet((*PMall)(nil))

	fmt.Println(pmall)
}

// ---------------------------------------------------------------------------------------------------------------------
// 测试编译是否通过
// ...interface{} 是可以看作是 []interface{}
func dump(i ...interface{}) {

}

func Test9(t *testing.T) {
	// 以下代码编译不通过
	//arr := []string{"1", "2"}
	//dump(arr...)
}

// ---------------------------------------------------------------------------------------------------------------------
// 等待多个协程完成工作 方式1
func wait1(n int) chan int {
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			fmt.Println("i:", i)
			ch <- i
		}(i)
	}
	return ch
}

func Test10(t *testing.T) {
	ch := wait1(10)
	for i := 0; i < 10; i++ {
		<-ch
	}
	fmt.Println("over..")
}

// ---------------------------------------------------------------------------------------------------------------------
// 等待协程退出，方式2
func wait2(n int) chan int {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			time.Sleep(6 * time.Second)
			fmt.Println("i:", i)
			wg.Done()
		}(i)
	}
	ch := make(chan int)
	go func() {
		wg.Wait()
		ch <- 1
	}()
	return ch
}

func Test11(t *testing.T) {
	ch := wait2(10)
	<-ch
	fmt.Println("over..")
}

// ---------------------------------------------------------------------------------------------------------------------
//等待协程退出，如果在规定的时间内没退出，则超时
func Test12(t *testing.T) {
	ch := wait2(10)
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-ch:
		fmt.Println("over")
	case <-timer.C:
		fmt.Println("timeout")
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//等待协程退出，方式3
func Test13(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Println("i:", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("over...")
}

// ---------------------------------------------------------------------------------------------------------------------
// 协程退出模式的应用-并发退出
// 这里使用到了把一个函数赋值给接口的应用，通过把一个函数赋值给一个接口
type GracefullyShutdowner interface {
	Shutdown() error
}

type ShutdownerFunc func() error

func (f ShutdownerFunc) Shutdown() error {
	return f()
}

func shutdownMaker(processTm int) func() error {
	return func() error {
		time.Sleep(time.Second * time.Duration(processTm))
		fmt.Println("执行：", processTm)
		return nil
	}
}

func ConcurrentShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	c := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			go func(shutdowner GracefullyShutdowner) {
				defer wg.Done()
				shutdowner.Shutdown()
			}(g)
		}
		wg.Wait()
		c <- struct{}{}
	}()

	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()
	select {
	case <-c:
		return nil
	case <-timer.C:
		return errors.New("wait timeout")
	}
}

func TestConcurrentShutdown(t *testing.T) {
	f1 := shutdownMaker(2)
	f2 := shutdownMaker(6)
	err := ConcurrentShutdown(10*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))

	if err != nil {
		t.Errorf("want nil, actual: %s", err)
		return
	}

	err = ConcurrentShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err == nil {
		t.Error("want timeout, actual nil")
		return
	}
}

// 所有的任务执行时间总和不得超过waitTimeout时间
func SequentialShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	start := time.Now()
	var left time.Duration
	// 设置等待时间
	timer := time.NewTimer(waitTimeout)
	// 遍历所有的退出器
	for i, g := range shutdowners {
		// 计算前面任务消耗了多少时间，需要超时时间减掉这个时间
		elapsed := time.Since(start)
		left = waitTimeout - elapsed
		c := make(chan struct{})
		go func(i int, shutdowner GracefullyShutdowner) {
			shutdowner.Shutdown()
			fmt.Println(i, " 执行成功")
			c <- struct{}{}
		}(i, g)
		// 重置timer
		timer.Reset(left)

		select {
		case <-c:
			// 继续执行
		case <-timer.C:
			return errors.New("wait timeout")
		}
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// 协程退出模式的应用-串行退出
func TestSequentialShutdown(t *testing.T) {
	f1 := shutdownMaker(6)
	f2 := shutdownMaker(4)

	err := SequentialShutdown(10*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))

	if err != nil {
		t.Errorf("want nil, actual: %s", err)
		return
	}

	err = SequentialShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err == nil {
		t.Error("want timeout, actual nil")
		return
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// channel的类似linux的管道使用方式
func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

// 获取偶数
func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

// 把数据进行平方
func square(in int) (int, bool) {
	return in * in, true
}

func spawn(f func(int) (int, bool), in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			r, ok := f(v)
			if ok {
				out <- r
			}
		}
		close(out)
	}()
	return out
}

func TestChannel(t *testing.T) {
	in := newNumGenerator(1, 20)
	// 先获取偶数在平方
	out := spawn(square, spawn(filterOdd, in))
	for v := range out {
		println(v)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 管道模式的扇入和扇出模式
// chapter6/sources/go-concurrency-pattern-9.go
/*func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

func square(in int) (int, bool) {
	return in * in, true
}*/

func spawnGroup(name string, num int, f func(int) (int, bool), in <-chan int) <-chan int {
	// 开启多个协程读取通道数据，然后处理到outSlice切片中
	// 扇出模式
	var outSlice []chan int
	for i := 0; i < num; i++ {
		out := make(chan int)
		go func(i int) {
			name := fmt.Sprintf("%s-%d:", name, i)
			fmt.Printf("%s begin to work...\n", name)
			for v := range in {
				r, ok := f(v)
				if ok {
					out <- r
				}
			}
			close(out)
			fmt.Printf("%s work done\n", name)
		}(i)
		outSlice = append(outSlice, out)
	}

	// 扇入模式
	//
	// out --\
	//        \
	// out ---- --> groupOut
	//        /
	// out --/
	groupOut := make(chan int)
	go func() {
		var wg sync.WaitGroup
		for _, out := range outSlice {
			wg.Add(1)
			go func(out <-chan int) {
				for v := range out {
					groupOut <- v
				}
				wg.Done()
			}(out)
		}
		wg.Wait()
		close(groupOut)
	}()

	return groupOut
}

func TestChannel2(t *testing.T) {
	in := newNumGenerator(1, 20)
	out := spawnGroup("square", 2, square, spawnGroup("filterOdd", 3, filterOdd, in))

	time.Sleep(3 * time.Second) //为了输出更直观的结果，这里等上面的goroutine都就绪
	for v := range out {
		fmt.Println(v)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 实现精细超时控制
type result struct {
	value string
}

// 指获取第一个结果
func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result)
	// 这里设置取消http请求的取消操作
	ctx, cancel := context.WithCancel(context.Background())
	// 一旦方法被执行完成了，直接取消其他的请求
	// 这个非常秒呀！
	defer cancel()
	queryFunc := func(i int, server *httptest.Server) {
		url := server.URL
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("query goroutine-%d: http NewRequest error: %s\n", i, err)
			return
		}
		req = req.WithContext(ctx)

		log.Printf("query goroutine-%d: send request...\n", i)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("query goroutine-%d: get return error: %s\n", i, err)
			return
		}
		log.Printf("query goroutine-%d: get response\n", i)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		c <- result{
			value: string(body),
		}
		return
	}

	for i, serv := range servers {
		go queryFunc(i, serv)
	}

	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New("timeout")
	}
}

func fakeWeatherServer(name string, interval int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		log.Printf("%s receive a http request\n", name)
		time.Sleep(time.Duration(interval) * time.Millisecond)
		w.Write([]byte(name + ":ok"))
	}))
}

func TestWait(t *testing.T) {
	result, err := first(fakeWeatherServer("open-weather-1", 200),
		fakeWeatherServer("open-weather-2", 1000),
		fakeWeatherServer("open-weather-3", 600))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
}

// ---------------------------------------------------------------------------------------------------------------------
// channel 一对一通信
func TestOneTwoOne(t *testing.T) {
	var ch = make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("111")
		ch <- 1
	}()
	<-ch
	fmt.Println("over")
}

// ---------------------------------------------------------------------------------------------------------------------
// 一对多通信
func TestOneToMany1(t *testing.T) {

	var wg sync.WaitGroup
	var ch = make(chan int)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			time.Sleep(1 * time.Second)
			fmt.Println("sun func", j)
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		ch <- 1
	}()
	<-ch
	fmt.Println("over")
}

// ---------------------------------------------------------------------------------------------------------------------
// 广播机制,通过关闭channel来实现通知所有的等待协程
func TestOneToWaitMany(t *testing.T) {
	var ch = make(chan int)
	for i := 0; i < 10; i++ {
		go func(j int) {
		loop:
			for {
				select {
				default:
					fmt.Println(j, ":睡眠1s")
					time.Sleep(1 * time.Second)
				case <-ch:
					fmt.Println(j, ":quit")
					break loop
				}
			}
		}(i)
	}

	time.Sleep(3 * time.Second)
	close(ch)
	fmt.Println("over")
	time.Sleep(3 * time.Second)
}

// ---------------------------------------------------------------------------------------------------------------------
// 通过锁的方式来实现数字的累加
var counter Counter

type Counter struct {
	mu    sync.Mutex
	count int
}

func add() {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	counter.count++
}
func TestCounter1(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go func() {
			add()
		}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println(counter.count)
}

// ---------------------------------------------------------------------------------------------------------------------
// 通过channel的方式来实现累加
// 通过一个协程中启动一个协程来不断的循环累增加数字来实现累加
var counter2 = Counter2{ch: make(chan int)}

type Counter2 struct {
	ch    chan int
	count int
}

func add2() int {
	return <-counter2.ch
}
func TestCounter2(t *testing.T) {
	go func() {
		for true {
			counter2.ch <- counter2.count
			counter2.count++
		}
	}()
	for i := 0; i < 1000; i++ {
		go func() {
			add2()
		}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println(counter2.count)
}

// ---------------------------------------------------------------------------------------------------------------------
// 通过有缓存通道进行信号量
// 工作量为10
var jobs = make(chan int, 10)

// 每次只能有3个协程工作
var signal = make(chan struct{}, 3)

func TestSignal(t *testing.T) {
	for i := 0; i < 10; i++ {
		jobs <- i + 1
	}
	var wg sync.WaitGroup
	for job := range jobs {
		wg.Add(1)
		signal <- struct{}{}
		go func(j int) {
			defer wg.Done()
			fmt.Println("job", j)
			time.Sleep(1 * time.Second)
			<-signal
		}(job)
	}
	wg.Wait()
	fmt.Println("over")
}

// ---------------------------------------------------------------------------------------------------------------------
// 尝试通过select来判断有缓存是否有数据

func trySend(ch chan int, i int) bool {
	select {
	case ch <- i:
		return true
	default:
		return false
	}
}

func tryRec(ch chan int) (int, bool) {
	select {
	case i := <-ch:
		return i, true
	default:
		return 0, false
	}
}

func producer(ch chan int) {
	var count = 0
	for {
		send := trySend(ch, count)
		if send {
			count++
			fmt.Println("发送成功")
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func consumer(ch chan int) {
	for {
		rec, ok := tryRec(ch)
		if ok {
			fmt.Println("获取成功：", rec)
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

func TestIfHasData(t *testing.T) {
	var ch = make(chan int, 3)
	go producer(ch)
	go consumer(ch)
	select {}
}

// ---------------------------------------------------------------------------------------------------------------------
// 测试无缓存chan和有缓存channel长度为1的区别
func TestCompare(t *testing.T) {
	/*	var ch = make(chan int)
		go func(ch chan int) {
			ch <- 1
			// 永远不会执行
			fmt.Println("11111111111")
		}(ch)
		time.Sleep(5 * time.Second)
		fmt.Println("over")*/

	var ch = make(chan int, 1)
	go func(ch chan int) {
		ch <- 1
		// 执行成功
		fmt.Println("11111111111")
	}(ch)
	time.Sleep(5 * time.Second)
	fmt.Println("over")
}

// ---------------------------------------------------------------------------------------------------------------------
