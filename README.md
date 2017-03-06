# venom [![Build Status](https://travis-ci.org/themalkolm/venom.svg?branch=master)](https://travis-ci.org/themalkolm/venom)

Add some venom to make [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper)
even more dangerous!

## Usage

See `_example/` folder for an example how to use venom. You can also use it as a skeletop for
any [12-factor](https://12factor.net) app you plan to use. It won't solve all problems but it will take care to solve
the [config](https://12factor.net/config) one i.e. it will allow you to store your configuration in the environment variables.

The twist is that it doesn't *require* you to store all configuration in the environment variables. It is up to you to define how exactly you want to configure the application. It is even OK to mix however you want:

* cli flags
* config file
* environment variables
