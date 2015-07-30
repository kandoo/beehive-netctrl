# Beehive Network Controller ![Travis Build Status](https://api.travis-ci.org/kandoo/beehive-netctrl.svg?branch=master)

This is a distributed SDN controller built on top of
[Beehive](http://github.com/kandoo/beehive). It supports
OpenFlow but can be easily extended for other southbound protocols.

Beehive network controller is high throughput, fault-tolerant and,
more importantly, can automatically optimize itself after a fault:

![Beehive Demo](http://raw.github.com/kandoo/beehive-netctrl/master/Docs/assets/beehive-optimization.gif)

Beehive network controller supports different forwarding and routing
methods, has automated discovery, end-to-end paths, and isolation.

Kandoo
------
You can find an implementation of Kandoo in
the [kandoo](https://github.com/kandoo/beehive-netctrl/tree/master/kandoo)
package.

