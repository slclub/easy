# Option Package


## Overview

It is suitable for service object initializing. Because of it has low performance.
we had used the packege of reflect to assigning between different Objects mapping assignment.
so, you could use it anywhere if does not worry about the performance of your program.

The main focus is convenience and powerful functionality.

## Install

```go get github.com/slclub/easy/vendors/option```

## Using


- Example OptionFunc

using option packge initialized the client of ETCD. used OptionFunc method to initializing.

```go
etcd.NewWithOption(option.OptionWith(nil).Default(
    option.OptionFunc(func() (string, any) {
        return "Endpoints", strings.Split(etcdAddr, ";")
    }),
))

	
```

- Example Struct

using any struct your defined to initialized the gprc.Client .

```go
client := cgrpc.NewClient(option.OptionWith(&struct {
    Name      string
    Namespace string
}{"server1", namespace}))
```

- Example 3

We can use both of them.

```go
func NewWithOption(assignment option.Assignment) {
    var err error
    
    v3config := clientv3.Config{}
    assignment.Target(&v3config)
    // set target default value
    // please running it before Apply() method.
    assignment.Default(option.OptionFunc(func() (string, any) {
        return "DialTimeout", CONNECTION_ETCD_TIMEOUT_DEFAULT * time.Second
    }))
    
    assignment.Apply()
    
    ecli, err = clientv3.New(v3config)
    if err != nil {
        log.Fatal("[ETCD] client created error")
        return
    }
}

```

When you need to invoke it, you can just use  fields or methods of the struct to finish your initializing.
```go
    // setting ETCD service
    eoption := &etcd.Option{}
    eoption.Conv(etcdAddr)
    // etcd.NewWithOption(option.OptionWith(eoption))
    etcd.NewWithOption(option.OptionWith(eoption).Default(option.DEFAULT_IGNORE_ZERO))
```

## Rules

- The target object ```option.Target(object)``` must be an struct pointer.
  if it is not a pointer you should use ```&``` (take address symbol).
- The name of fields or methods that come from source object sould be the same as target.
- The fields first letter of target can be in lowercase or uppercase.
- The first letter of fields of the source object must be in uppercase. otherwise  will be ignored.
- Both of fields or methods in source object are allowed to initializing the target object, If it can match.
- The source object can be nil.

example:
- Valid assignment

```go
target.name  = source.Name
target.Name  = source.Name
target.name  = source.Name() string
target.Name  = source.Name() string
```

- Unvalid

```go
target.name  != source.name
target.Name  != source.name
target.name  != source.name() string
target.Name  != source.name() string
target.Name  != source.NameX
```


## Configure Part of Option

- Default; (use one or more OptionFunc to initialize the target object)
- Source; Object (config object)
- Final; (use one or more OptionFunc to nitializing)

Final is the last execution. so the value of target field must be same with the value comed from Final. 
