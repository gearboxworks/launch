package cmd


var JsonFileSchema = `
{
	"schema": "gear-2",
	"meta": {
		"state": "production",
		"organization": "gearboxworks",
		"name": "tinygo",
		"label": "TinyGo Container",
		"info": "TinyGo is a project to bring the Go programming language to microcontrollers and modern web browsers by creating a new compiler based on LLVM.",
		"description": [
			""
		],
		"maintainer": "Gearbox Team <team@gearbox.works>",
		"class": "development",
		"refurl": "https://tinygo.org/"
	},
	"build": {
		"fixed_ports": {
			"http": "80",
			"https": "443"
		},
		"ports": {
			"http": "80",
			"https": "443"
		},
		"run": "",
		"args": "",
		"workdir": "",
		"env": {
			"TINYGO_USER": "gearbox"
		},
		"network": "--network gearboxnet",
		"volumes": "",
		"restart": "--restart no"
	},
	"run": {
		"commands": {
			"default": "/usr/local/tinygo/bin/tinygo",
			"tinygo": "/usr/local/tinygo/bin/tinygo"
		}
	},
	"project": {
	},
	"extensions": {
	},
	"versions": {
		"0.17.0": {
			"majorversion": "0.17",
			"latest": true,
			"ref": "tinygo/tinygo:0.17.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.16.0": {
			"majorversion": "0.16",
			"latest": false,
			"ref": "tinygo/tinygo:0.16.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.15.0": {
			"majorversion": "0.15",
			"latest": false,
			"ref": "tinygo/tinygo:0.15.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.14.1": {
			"majorversion": "0.14",
			"latest": false,
			"ref": "tinygo/tinygo:0.14.1",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.13.1": {
			"majorversion": "0.13",
			"latest": false,
			"ref": "tinygo/tinygo:0.13.1",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.12.0": {
			"majorversion": "0.12",
			"latest": false,
			"ref": "tinygo/tinygo:0.12.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.11.0": {
			"majorversion": "0.11",
			"latest": false,
			"ref": "tinygo/tinygo:0.11.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.10.0": {
			"majorversion": "0.10",
			"latest": false,
			"ref": "tinygo/tinygo:0.10.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.9.0": {
			"majorversion": "0.9",
			"latest": false,
			"ref": "tinygo/tinygo:0.9.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.8.0": {
			"majorversion": "0.8",
			"latest": false,
			"ref": "tinygo/tinygo:0.8.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.7.1": {
			"majorversion": "0.7",
			"latest": false,
			"ref": "tinygo/tinygo:0.7.1",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.6.1": {
			"majorversion": "0.6",
			"latest": false,
			"ref": "tinygo/tinygo:0.6.1",
			"base": "gearboxworks/gearbox-base:debian-buster"
		},
		"0.5.0": {
			"majorversion": "0.5",
			"latest": false,
			"ref": "tinygo/tinygo:0.5.0",
			"base": "gearboxworks/gearbox-base:debian-buster"
		}
	}
}`
