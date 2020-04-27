---
title: "Python Config"
date: 2020-04-16T22:20:05+01:00
showDate: true
tags: ["python","environs", "configuration"]
---

## My desire

Application configuration, one of the tenets of the [12 factor app][12-factor-config][^1]. Such a
reoccurring and important need that it deserves a reusable elegant solution. This blog post assumes
you've already drunk the Kool-Aid. If not perhaps follow the link before proceeding, they'll do a better
job than I in convincing you.

[12-factor-config]: https://12factor.net/config

[^1]: It is with great embarrassment that I confess that despite oft-quoting this methodology, I've
never read all of it, nor even do I know what all 12 factors _are_. It's on the to-do list.

From my recent development in Go, I've become familiar with [cobra] and its close relative [viper].
I think these tools readjusted what I have come to be done 'for me', and I have less patience when it
comes to rolling my own configuration from environment variables.

[cobra]: https://github.com/spf13/cobra
[viper]: https://github.com/spf13/viper

What do I mean? Well you define a command line application with cobra, including options and arguments.
Viper allows you to specify these options and arguments in other common ways: via a configuration file,
_or pulling values from environment variables_.

## What about Python

So enter Python. Let's say you've been prototyping some sort of client that makes requests:

```python
## some_sub_module.py
class Client:
    def configure():
        # Some preamble ...
        self._aws.configureEndpoint(
            "some-account.aws.com",
            8883)
```

You're a good citizen, you recognise that you should use environment variables to configure those parameters.
So that if you got a new account or, God forbid, anyone else actually wants to use your code.
You probably reach for:

```python
## some_sub_module.py
class Client:
    def configure():
        # Some preamble ...
        self._aws.configureEndpoint(
            os.getenv(APP_ENDPOINT, "some-account.aws.com"),
            os.getenv(APP_PORT, 8883))
```

OK, but remember, all environment variables are parsed as _strings_, urgh:

```diff
## some_sub_module.py
class Client:
    def configure():
        # Some preamble ...
        self._aws.configureEndpoint(
            os.getenv(APP_ENDPOINT, "some-account.aws.com"),
-           os.getenv(APP_PORT, 8883))
+           int(os.getenv(APP_PORT, "8883")))
```

Fine, but are you not worried now that parsing `APP_PORT` as an `int` might throw an exception if it's
been set to something silly? Your users may well have done. Sure we love exception flow in Python, but
are you really going to wrap `configure()` in a try catch for the generic `ValueError`?

And further to that, if you depend on an environment variable to be set who wants to wait for `configure()`
to be called to find out that your user hasn't configured the program correctly?

Frankly the more I think about the problem, the less time I have for it. I want to solve problems, not
write boilerplate.

## Environs, the 'hero'

I'll usually create a config class to encapsulate the application configuration:

```python
## config.py
from dataclasses import dataclass

@dataclass
class Config:
    endpoint: str
    port: int
```

Now let's add a way of creating a populated `Config` object from the environment. Enter [environs](https://github.com/sloria/environs),
a Python package that __isn't__ as good as Viper, but hey, anything that allows me to get on with my
day:

```python
## config.py
from dataclasses import dataclass
from environs import Env

@dataclass
class Config:
    endpoint: str
    port: int

    @classmethod
    def from_env(cls):
        """Parse configuration from environment variables
        Returns:
            Config: Application configuration
        """
        env = Env()
        # Read environment variables from .env file (if present)
        env.read_env()
        # Create a new Config from environment variables
        return Config(
            endpoint=env.str("APP_ENDPOINT"),
            port=env.int("APP_PORT"))
```

Do you see the magic? I don't have to worry about parsing strings to integers, environs has my back.
Plus I can use a [`.env` file][dotenv] to specify my environment variables on disk.

[dotenv]: https://www.freecodecamp.org/news/nodejs-custom-env-files-in-your-apps-fa7b3e67abe1/

Finally, I just need to create my config as soon as my program starts, and if my users have failed to
correctly configure the program it will exit immediately. So much better than finding out at the time
of use!

```python
from myapp.config import Config

def main():
    # Create the application config
    config = Config.from_env()
    # Do the rest
    # ...

if __name__ == '__main__':
    main()
```

It's worth pointing out, that this sort of approach can make code more testable as well. It's easy to
create a `Config` object with mocked values and inject it into whatever unit is under test.

Right what was I doing again?...
