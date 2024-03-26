<p align="center">
  <img src="https://github.com/njern/broadcast/blob/master/broadcast.png?raw=true" alt="Broadcast logo"/>
</p>

<br><br>

<div align="center">

[![License](https://img.shields.io/badge/license-MIT-blue)](#license)
[![Go Report Card](https://goreportcard.com/badge/github.com/njern/broadcast)](https://goreportcard.com/badge/github.com/njern/broadcast)
[![Issues - broadcast](https://img.shields.io/github/issues/njern/broadcast)](https://github.com/njern/broadcast/issues)
![GitHub Release](https://img.shields.io/github/release/njern/broadcast)

</div>

<h1 align="center">The #1 Go package for implementing the publish-subscribe pattern</h1>
<div align="center">
Publish any type of data to any number of subscribers.
</div>
</br>

<p align="center">
    <a href="https://github.com/njern/broadcast/issues/new?assignees=&labels=bug&projects=&template=bug_report.md&title=%F0%9F%90%9B+Bug+Report%3A+">Report Bugs</a>
    ·
    <a href="https://github.com/njern/broadcast/issues/new?assignees=&labels=enhancement&projects=&template=feature_request.md&title=%F0%9F%9A%80+Feature%3A+">Request Features</a>
    ·
    <a href="https://twitter.com/njern">Twitter</a>
  </p>

## About
This package provides a simple and efficient way to implement a publish/subscribe pattern in Go, allowing one sender to broadcast messages to multiple receivers. It supports dynamic subscription and unsubscription, ensuring flexibility and control over message delivery and resource management.

## Features

- **Generic Implementation:** Works with any data type.
- **Timeout Control:** Customize the time to wait for subscribers to receive messages.
- **Dynamic Subscription Management:** Subscribers can join and leave at any time.
- **Automatic Cleanup:** Automatically closes all subscriber channels when the broadcaster is closed.

## Getting Started

### Installation

To use the broadcast library in your project, first, ensure your project is initialized as a Go module, then add the library to your project with:

```bash
go get github.com/njern/broadcast
```

### Basic Usage
Start by importing the package in your Go file.

```go
import (
    "github.com/njern/broadcast"
)
```

Create a broadcaster with a specified buffer size and timeout for broadcast operations.

```go
b := broadcast.New[string](10, 0)
```

Subscribers can now subscribe to the broadcaster, receiving a channel to listen for messages. The caller provides the channel buffer size as input (or `0` if you prefer an unbuffered channel).

```go
ch, err := b.Subscribe(2)
if err != nil {
    log.Fatalf("Failed to subscribe: %v", err)
}

go func() {
    for msg := range ch {
        fmt.Println("Received:", msg)
    }
}()
```


Send messages to all subscribers by sending them to the broadcaster's channel.

```go
b.Chan() <- "Hello, Broadcasters!"
```

Subscribers can stop receiving messages by unsubscribing.

```
b.Unsubscribe(ch)
```

When the broadcaster is no longer needed, close it to release all resources.

```go
b.Close()
```


### Custom Timeouts
You can control how long the broadcaster waits for subscribers to receive messages. This is set when creating the  broadcaster and applies to all messages.

```go
// This broadcaster will wait up to 10 seconds for subscribers to  
// receive a message before continuing to the next subscriber.
b := broadcast.New[int](10, 10*time.Second)
```

### Handling Closed Broadcasters
Attempting to subscribe to a closed broadcaster will result in an `ErrBroadcasterClosed` error.

```go
_, err := b.Subscribe(0)
if err == broadcast.ErrBroadcasterClosed {
    fmt.Println("Broadcaster has been closed.")
}
```

## Contributing

Contributions to improve this library are welcome. Feel free to fork the repository, make your changes, and submit a pull request.

## License

This library is licensed under the [MIT License](LICENSE).

## Acknowledgements

- I used the tips from [this](https://www.daytona.io/dotfiles/how-to-write-4000-stars-github-readme-for-your-project) very neat article by the [Daytona](https://www.daytona.io) team to put together the README.