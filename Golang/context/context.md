###Context (上下文)
#### context 可以用来设置截止日期、同步信号，传递请求相关值。



    type Context interface {
        Deadline() (deadline time.Time, ok bool)
        Done() <-chan struct{}
        Err() error
        Value(key interface{}) interface{}
    }
#### 1. Deadline 
##### - 返回这个context任务结束该结束的时间
#### 2. Done
##### - 当工作完成代表context应该取消的时候返回一个通道
##### - Done可能返回nil, 如果context永远不会被取消
##### - 连续调用Done会返回同样的值
##### - 在cancel函数返回之后, 完成关闭通道是可能异步发生的
##### - WithCancel安排在调用cancel时关闭Done
##### - WithDeadline安排在截止日期到期时关闭Done
##### - WithTimeout安排在超时结束时关闭Done
#### 3. Err        
##### - 如果Done还没有关闭, Err会返回nil
##### - 如果Done已经关闭, Err会返回不是nil的错误来解释原因:
##### - i. Canceled 如果context已经被取消的情况返回
##### - ii. DeadlineExceeded 如果context的截止时间已经超过
##### - Err返回了一个不是nil的值, 那么连续调用都会返回相同的错误
#### 4. Value
##### - 返回与context键相对应的值, 如果没有与key相关联的值返回nil
##### - 连续的调用会返回相同的值

##### 不要在结构体中存储context, 相反需要显示的传递给每个需要它的函数, 同时context应该是一个参数
##### 给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO()
##### 即使函数允许,也不要传递空的context
