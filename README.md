# venom [![Build Status](https://travis-ci.org/themalkolm/venom.svg?branch=master)](https://travis-ci.org/themalkolm/venom)

Add some venom to make [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper)
even more dangerous!

## Usage

See [_example](https://github.com/themalkolm/venom/tree/master/_example) folder for an example how to use venom. You can
also use it as a skeleton for any [12-factor](https://12factor.net) app you plan to use. It won't solve all problems but
it will take care to solve the [config](https://12factor.net/config) one i.e. it will allow you to store your
configuration in the environment variables.

The twist is that it doesn't *require* you to store all configuration in the environment variables. It is up to you to
define how exactly you want to configure the application. It is even OK to mix however you want:

* [cli flags](https://github.com/spf13/cobra#working-with-flags)
* [config file/dir](https://github.com/spf13/viper#reading-config-files)
* [environment variables](https://github.com/spf13/viper#working-with-flags)
* [...](https://github.com/spf13/viper#what-is-viper)

## Priority

If you use `TwelveFactorCmd` then here is the priority of resolution (highest to lowest):

1. `$ example --foo 42`
2. `$ example -e EXAMPLE_FOO 42`
3. `$ example --env-file example.env # (assuming it has EXAMPLE_FOO=42 line)`
4. `$ EXAMPLE_ENV=EXAMPLE_FOO=42 ./bin/example`
5. `$ EXAMPLE_ENV_FILE=example.env ./bin/example`
6. `$ EXAMPLE_FOO=42 example`
