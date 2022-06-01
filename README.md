# 笔记

## [Go语言学习](https://draveness.me/golang/)

### [defer](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)

Go语言的`defer`会在当前函数返回前执行传入的函数, 经常被用于关闭文件描述符, 关闭数据库连接以及解锁资源

使用`defer`的最常见场景是在函数调用后完成一些收尾工作, 例如在`defer`中回滚数据库的事务

```go
func createPost(db *gorm.DB) error {
    tx := db.Begin()
    defer tx.Rollback()
    
    if err := tx.Create(&Post{Author: "Draveness"}).Error; err != nil {
        return err
    }
    
    return tx.Commit().Error
}
```

在使用数据库事务时, 我们可以使用上面的代码, 在创建事务后就立刻调用`defer tx.Rollback()`保证事务一定会回滚

哪怕事务真的执行成功了, 那么调用`tx.Commit()`之后再执行`tx.Rollback()`也不会影响已提交的事务, 而如果执行失败了, 则可以保证`tx.Rollback()`最终会被执行, 从而保证事务被回滚

#### 1.现象

我们在Go语言中使用`defer`时会遇到两个常见问题, 这里会介绍具体的场景并分析这两个现象背后的设计原理:

- `defer`关键字的调用时机以及多次调用`defer`时执行顺序是如何确定的
- `defer`关键字使用传值方式传递参数时会进行预计算, 导致不符合预期的结果

##### 作用域

向`defer`关键字传入的函数会在函数返回前运行. 假设我们在`for`循环中多次调用`defer`关键字

```go
func main() {
	for i := 0; i < 5; i++ {
		defer fmt.Println(i)
	}
}
```

运行结果如下

```shell
$ go run main.go
4
3
2
1
0
```

运行上述代码会倒叙执行传入`defer`关键字的所有表达式, 因为最后一次调用`defer`传入了`fmt.Println(4)`, 所以这段代码会优先打印4. 我们可以通过下面这个简单例子强化对`defer`执行时机的理解

```go
func main() {
    {
        defer fmt.Println("defer runs")
        fmt.Println("block ends")
    }
    
    fmt.Println("main ends")
}
```

运行结果如下

```shell
$ go run main.go
block ends
main ends
defer runs
```

从上述代码的输出我们会发现, `defer`传入的函数不是在退出代码块的作用域时执行的, 它只会在当前函数和方法返回之前被调用

#### 预计算参数

Go语言中所有的函数调用都是传值的, 虽然`defer`是关键字, 但是也继承了这个特性. 假设我们想要计算`main`函数运行的时间, 可能会写出以下的代码

```go
func main(){
    startedAt := time.Now()
    defer fmt.Println(time.Since(startedAt))

    time.Sleep(time.Second)
}
```

运行结果如下

```shell
$ go run main.go
0s
```

然而上述代码的运行结果并不符合我们的预期, 这个现象背后的原因是什么呢? 经过分析, 我们会发现调用`defer`关键字会立刻拷贝函数中引用的外部参数, 所以`time.Since(startedAt)`的结果不是在`main`函数退出前计算的, 而是在`defer`关键字调用时计算的, 最终导致上述代码输出0s

想要解决这个问题的方法非常简单, 我们只需要向`defer`关键字传入匿名函数

```go
func main(){
    startedAt := time.Now()
    defer func() { fmt.Println(time.Since(startedAt)) }()

    time.Sleep(time.Second)
}
```

运行结果如下

```shell
$ go run main.go
1s
```

虽然调用`defer`关键字时也使用值传递, 但是因为拷贝的是函数指针, 所以`time.Since(startedAt)会在`main`函数返回前调用并打印出符合预期的结果

### recover

