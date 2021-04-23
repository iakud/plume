## network

#### TCPServer

通过`network.NewTCPServer`创建`TCPServer`。

`TCPServer`使用`ListenAndServe`开启监听，接受新建立的连接并处理。

```go
server := network.NewTCPServer("localhost:8000")
err := server.ListenAndServe(handler, codec)
if err != network.ErrServerClosed {
    // handle error
}
```

`handler`和`codec`实现连接事件处理和数据读写处理的接口，可以通过接口来实现自定义连接事件和数据读写处理。

`TCPServer`会对每个`TCPConnection`使用单独的协程来处理连接事件。

**Error Handling**

当`ListenAndServe`遇到不可恢复的错误时会将错误返回。

调用`Close`后，`ListenAndServe`会返回`network.ErrClientClosed`。

#### TCPClient

通过`network.NewTCPClient`创建`TCPClient`。

`TCPClient`使用`DialAndServe`主动建立连接并处理。

```go
client := network.NewTCPClient("localhost:8000")
err := client.DialAndServe(handler, codec)
if err != network.ErrClientClosed {
    // handle error
}
```

`handler`和`codec`实现连接事件处理和数据读写处理的接口，可以通过接口来实现自定义连接事件和数据读写处理。

`EnableRetry`和`DisableRetry`用于启用和禁用当连接建立遇到错误或者失败时重试。

**Error Handling**

当重试未启用，`DialAndServe`遇到不可恢复的错误时会将错误返回。

调用`Close`后，`DialAndServe`会返回`network.ErrClientClosed`。

#### Codec接口

通过实现`Read`和`Write`接口来实现`conn`的数据读写处理方法。

```go
type Codec interface {
    Read(r io.Reader) ([]byte, error)
    Write(w io.Writer, b []byte) error
}
```

数据处理接口：

`Read`方法实现如何从`conn`读取数据，返回读取的数据和是否遇到错误
`Write`方法实现如何将原始数据`b`写入`conn`，返回写入过程是否遇到错误

如果`TCPServer`或`TCPClient`的`codec`参数被设置为`nil`，则会调用`network.DefaultCodec`。

#### TCPHandler接口

通过实现`TCPHandler`来实现`conn`连接事件处理方法。

```go
type TCPHandler interface {
    Connect(conn *TCPConnection, connected bool)
    Receive(conn *TCPConnection, buf []byte)
}
```

连接处理接口：

- `Connect`为连接事件，`connected`表示事件是建立或断开
- `Receive`为接收事件，`buf`表示接收到的内容（由Codec定义）
- `conn`是与事件相关的连接

如果`TCPServer`或`TCPClient`的`handler`参数被设置为`nil`，则会使用`network.DefaultTCPHandler`。

`network.DefaultTCPHandler`会忽略所有的连接事件。

#### Close & Graceful Shutdown

`Close`可以主动关闭连接，同时会直接丢弃队列中未发送的数据和丢弃接收缓冲区未读取的数据。

`Shutdown()`将设置关闭状态，并在给定的超时时间后（默认3秒）完全关闭`socket`，当发送协程发现关闭状态被设置，则会在缓冲队列所有数据都写入`socket`后调用`conn.CloseWrite()`通知对端已经写入完毕，此时等待对端主动关闭

或者，直接调用`ShutdownIn(d time.Duration)`例如：

```go
conn.ShutdownIn(time.Second * 3)
```

#### 吞吐量测试

乒乓测试（单机）

```
9827028992 total bytes read
2399177 total messages read
4096 average message size
937 MiB/s throughput
```

详细代码见`pingpong_test.go`