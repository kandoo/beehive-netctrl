Kandoo
======
This package implements
[Kandoo](http://dl.acm.org/citation.cfm?id=2342441.2342446) on top of 
[Beehive SDN/OpenFlow Controller](https://github.com/kandoo/beehive-netctrl).

To boot the first controller, run:

```
# go run main/main.go -addr ADDR1:PORT1
```

where `ADDR1` is the listening address and `PORT1` is the listening port.

To connect a new controller running on another machine to
your first controller, run:

```
# go run main/main.go -addr ADDRN:PORTN -paddrs ADDR1:PORT1
```

All controllers will be local controllers and one of them automatically
becomes the root controller.

