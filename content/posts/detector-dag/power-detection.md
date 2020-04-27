---
title: "Power detection"
date: 2020-04-26T12:48:14+01:00
showDate: true
series: ["detector-dag"]
series_weight: 3
tags: ["raspberry-pi"]
---

I connected a Raspberry Pi to the [UPS-Lite](https://www.aliexpress.com/item/32954180664.html) that I'd
identified in the [concepts]({{< relref "concepts.md" >}}) blog post. Right out of the box I could switch
off power and the Pi would stay alive!

![UPS-Lite schematic](/upslite-and-pi.jpg)

I moved on to _detecting_ the power state changes with the Pi: I wanted to send MMQT messages
with my [AWS IoT client]({{< relref "./aws-iot.md" >}}) when these events occurred.

Documentation on the UPS-Lite was thin on the ground, but I'd found a tweet
from the manufacturer that indicated I could detect this on the Pi's GPIO7:

{{< tweet 1069568475683647488 >}}

I selected the higher-level [gpiozero] python library over [RPi.GPIO][rpi-gpio]. It appears to be the
[recommended library][gpiozero-recommended] these days, and critically for me, it has a simple edge detection
API. This allowed me to declaratively assign callbacks to be executed when an input pin changes
value - tidy.

```python
from gpiozero import DigitalInputDevice  # noqa: E501, pylint: disable=import-error,import-outside-toplevel
_POWER_PIN = 4  # gpiozero uses Broadcom pin numbering

power_status_input = DigitalInputDevice(_POWER_PIN, bounce_time=0.2)
power_status_input.when_activated = _publish_update
power_status_input.when_deactivated = _publish_update

def _publish_update(self, device: DigitalInputDevice) -> None:
    print("Status changed!")
```

[gpiozero]: https://gpiozero.readthedocs.io/en/stable/#
[rpi-gpio]: https://sourceforge.net/projects/raspberry-gpio-python/
[gpiozero-recommended]: https://www.raspberrypi.org/documentation/usage/gpio/python/README.md

I tore my hair out trying to get that pin detection to work. I found the manufacturer's [GitHub repo][ups-lite-repo],
but the example code demonstrated reading the voltage/capacity of the battery, not the power status.

[ups-lite-repo]: https://github.com/linshuqin329/UPS-Lite

I double-checked the [pin numbering system], and confirmed it _had_ to be pin BCM4/WiringPi7, not least
because I could _see_ BCM7 wasn't connected to by a pogo pin.

[pin numbering system]: https://pinout.xyz/#

I manually read the pin a few times and whilst knocking it about found the voltage seemed to be floating:
sometimes it was high, and sometimes low. Going _really_ spare I got my multi-meter out and confirmed
it.

I studied the owners manual, written in Chinese, pouring over the same wiring schematic I'd found in
the tweet: everything looked kosher.

Finally, I found a [translated manual][translated-manual] that revealed my error:

[translated-manual]: https://github.com/linshuqin329/UPS-Lite/issues/1

> Also UPS-Lite insertion detecting function with a power adaptor, the insertion of the power pi io4
(BCM number) detects the high level, when pulled low, enabling the weld shorting function requires two
back of the UPS disc, as shown below in detail.

![UPS-Lite solder pads](/upslite-pads.png)

I needed to short the two solder pads on the back of the board to have the pogo pin carry the power
status signal! If my electronics knowledge had been better this is actually obvious from the schematic:

![UPS-Lite schematic](/upslite-schematic.png)

...I'm telling you, I felt like a right fool.
