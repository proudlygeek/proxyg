Proxyg
======

A Very, very simple HTTP Proxy implementation.

Why?
----

This is just a *simple* HTTP Proxy implementation mainly used to study *Go*, it has **no auth** and it's probably buggy: i don't recommend using it for *serious stuff*.
I coded it using [another implementation][1] found on GitHub.

Installation
------------

```bash
go get github.com/proudlygeek/proxyg
```

Make sure you've exported **GOPATH**/bin to your **PATH** env var.

Usage
-----

To run an HTTP Proxy simply run:

```bash
proxyg
```

By default, *proxyg* listens on *localhost* port *8080*; if you need to change the default behavior just pass the **-host** and **-port** flags:

```bash
proxyg -host="0.0.0.0" -port=9081
```

License
-------

MIT


[1]: https://github.com/rmt/httpconnectproxy

