# venom [![Build Status](https://travis-ci.org/themalkolm/venom.svg?branch=master)](https://travis-ci.org/themalkolm/venom)

Add some venom to make [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper)
even more dangerous!

## Usage

See [_example](https://github.com/themalkolm/venom/tree/master/_example) folder for an example how to use venom. You can also use it as a skeletop for
any [12-factor](https://12factor.net) app you plan to use. It won't solve all problems but it will take care to solve
the [config](https://12factor.net/config) one i.e. it will allow you to store your configuration in the environment variables.

The twist is that it doesn't *require* you to store all configuration in the environment variables. It is up to you to define how exactly you want to configure the application. It is even OK to mix however you want:

* [cli flags](https://github.com/spf13/cobra#working-with-flags)
* [config file/dir](https://github.com/spf13/viper#reading-config-files)
* [environment variables](https://github.com/spf13/viper#working-with-flags)
* [...](https://github.com/spf13/viper#what-is-viper)

## Gotchas

This is very nice to allow you to configure how you want. Interesting that this makes it hard to make sure that every developer is
using the same configuraiton and they are easily translated between each other. It is matter of time when you have the following questions to answer:

* How to define environment variables in a config file?
* How to pass a configuration file content as an env variable?

There is no "right" and easy way to solve this problem. The more flexible is configuration - the more complex setups you get. Yep, here `venom` goes long way to break the very thing 12-factor apps is about.

✌.|•͡˘‿•͡˘|.✌
