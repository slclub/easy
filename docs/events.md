# Event Asynchronization

## Install

```go get github.com/slclub/go-tips/events/module_event```

## Feature

- Asynchrionization;
- Two selectable models; Synchronous and Asynchronous;
- Customed pool of goroutine; implement module_event.AsyncSubmiter;

## Simple

- New
- Register
- Submit
- Emit
- Close

```go
func TestEventCommon(t *testing.T) {
	// new;  event master
    eventMonitor := NewEvent(&Option{})
    
    // register; message ID , message Handle binding them to the object
    eventMonitor.Register(EventHandle(handleLogin), &TEST_EVENT_ID_LOGIN)                        // register fail
    eventMonitor.Register(EventHandle(handleLogout), &TEST_EVENT_ID_LOGOUT)                      // register fail
    eventMonitor.Register(EventHandle(handleTrace), &TEST_EVENT_ID_LOGIN, &TEST_EVENT_ID_LOGOUT) // regsiter fail
    eventMonitor.Register(EventHandle(handleLogout), nil)                                        // register fail
    eventMonitor.Register(HandleConvert(anotherHandle), &TEST_EVENT_ID_ANOTHER)                  // register ok
    
    // trigger; submit event to queue of running.
    eventMonitor.Submit(&EventOper{&TEST_EVENT_ID_LOGIN, "", []any{t, 1, 2, "event"}}) // ok
    eventMonitor.Trigger(&TEST_EVENT_ID_ANOTHER, t, 3, "gold")
    
    // emit ; running all of the event handle
    eventMonitor.Emit()
    time.Sleep(10 * time.Millisecond)
    
    // release;
    eventMonitor.Close()
}
```

- Event Muster has more options; example:

```go
eventMonitor := NewEvent(&Option{
    InOrder:        false,  // 是否让所有事件 顺序同步执行
    TimeTickPeriod: time.Duration(10 * time.Millisecond), //不启用轮训机制 设置成0 即可
    Submiter:       antsPool(), // 植入自己的携程池
})
```

## Event ID type

We can define any type as the event id. But you need to implement the EventValue interface.

```go
type EventValue interface {
    Value() string
}
```

Default Event ID types. They implement the EventValue interface.

```go
// 事件ID 分两种，都继承了 EventValue
type EVENT_ID_INT int
type EVENT_ID_STRING string
```


## Event Handle

```go
type EventHandle func(oper *EventOper)
```

## HandleConvert

For your convenience, you can use it to converting the func(args ...any) type to EventHandle

```go
HandleConvert(fn func(args ...any)) EventHandle
```

## Event ID and Handle

Their relationship is many to many.

## Customize

- Implement AsyncSubmiter

```go
// routine pool interface
// please remenber use Release method to free your memory of pool
type AsyncSubmiter interface {
    Submit(func()) error
    Release()
}
```





