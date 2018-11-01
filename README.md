# glitzz [![Build Status](https://travis-ci.org/lovelaced/glitzz.svg?branch=master)](https://travis-ci.org/lovelaced/glitzz) [![codecov](https://codecov.io/gh/lovelaced/glitzz/branch/master/graph/badge.svg)](https://codecov.io/gh/lovelaced/glitzz)
![glitzz logo](logo.png)

**glitzz** is an IRC bot which is optimized for shitposting and the basic features any IRC channel would likely enjoy.

## Getting Started

If you wanna install this trash:

    $ go get -u github.com/lovelaced/glitzz/cmd/glitzz

You'll need a config; luckily for your lazy ass there's a sample one provided. Change it to suit your liking.

    $ glitzz default_config > config.json
    $ glitzz run config.json

You can also run `glitzz --help` for help.
