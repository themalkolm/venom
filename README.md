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

* `$ example --foo 42`
* env as flags
 * `$ example -e EXAMPLE_FOO 42`
 * `$ example --env-file example.env # (assuming it has EXAMPLE_FOO=42 line)`
* env as env
 * `$ EXAMPLE_ENV=EXAMPLE_FOO=42 ./bin/example`
 * `$ EXAMPLE_ENV_FILE=example.env ./bin/example`
* `$ EXAMPLE_FOO=42 example`

You probably should not use env as env trick as it is very confusing for any user.

## Autoflags

It is possible to ask venom to define flags for you. You need to provide a struct or pointer to struct that has
special `flag` tag set e.g.

```
type Config struct {
	FooMoo int `flag:"foo-moo,m,Some mooness must be set"`
}
```

This will allow venom to find this tag and parse long flag name, short flag name and the description. It expect you to
define it as a comma separated triplet. It has some logic to deduce what you meant in case you have use only one or two
comma separated values.

To define flags you simply run `DefineFlags`. Note that for flags only the type of the passed object matters. The struct
field values are ignored:

```
flags := venom.DefineFlags(Config{})
RootCmd.PersistentFlags().AddFlagSet(flags)
```
