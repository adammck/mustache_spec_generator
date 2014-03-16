## Mustache Spec Generator

This is a little tool to generate an extensive set of integration tests for
[mustache.go] [mgo], straight from [the Mustache spec] [spec]. I'm working on
getting this merged upstream [in #43] [pr].

[mgo]: https://github.com/hoisie/mustache
[spec]: https://github.com/mustache/spec
[pr]: https://github.com/hoisie/mustache/pull/43


## Usage

    git clone --recursive https://github.com/adammck/mustache_spec_generator.git
    cd mustache_spec_generator
    go build
    ./mustache_spec_generator.go


## License

MIT.
