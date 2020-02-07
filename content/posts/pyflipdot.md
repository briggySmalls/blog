---
title: "pyflipdot driver"
date: 2020-02-06T21:35:01Z
showDate: true
draft: true
series: ["flipdot"]
tags: ["blog","story"]
---

Before purchasing the signs I’d scoped out some ‘prior art’ and found that some existing drivers were out there.

I purchased a USB to RS232 adapter and initially I tried out the more ‘complete’ repo. After fixing a couple of bugs[^1] I had the signs displaying text!

[^1]: PRs: [use address argument](https://github.com/tuna-f1sh/node-flipdot/pull/3) and [remove blank fill](https://github.com/tuna-f1sh/node-flipdot/pull/4)

Having proved the signs worked I assessed the situation. Knowing my natural tendency to rewrite everything myself, I wanted to force myself to consider using this driver. Some considerations were:

- it was written in node
- It was a monolith
- It only seemed to have one or two usable fonts
- It didn’t appear to support drawing an arbitrary picture matrix, instead it’s api only accepted text
- It didn’t support multiple signs

Ok who am I kidding, I wanted to reinvent the wheel. Nicely. Beautifully. Simply. Readable. And unit tested so that future people could use it whatever their purpose. I was going to build the pyflipdot driver.

I chose python because at the time it was my go to, get results quick, language. I used my trusty cookie cutter template to initialise all the good housekeeping: unit testing, formatting, linting, CI pipeline, etc. I like to keep things ship shape.

I can’t take any credit for decoding the protocol. That was fairly well documented at this repo. All I did was give it a sensible API: create a controller that holds the serial connection and add to it as many signs as you like. Send images to the songs in the form of Boolean bumpy arrays. Tidy.

You can find abs install the package using pip in the normal way.
