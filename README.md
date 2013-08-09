# cw-stats-server

Cube World doesn't have a query interface.. so we're blazing ahead with a JSON one for server listers.

This works by transforming stats tracked through our monitoring scripts into the
*de facto* json data structure defined by [cuwo](https://github.com/matpow2/cuwo).

## Example Response:

```json
{
    "ip": "cw.gcg.io", 
    "location": "US", 
    "max": 144, 
    "mode": "default", 
    "name": "GC Gaming cw.gcg.io connects to best of 12 servers", 
    "players": 72
}
```

## Getting started

At the moment this is highly tailored to talk with http://status.gcg.io. But
assuming you're endpoint exposes a series of servers in the following format it
should work for you too:

	{
	    "servers": [
	        {
	            "ip": "cw1.gcg.io",
	            "max": 12,
	            "current": "12"
	        },
	        {
	            "ip": "cw2.gcg.io",
	            "max": 12,
	            "current": "2"
	        }
	    ]
	}

It will accumulate all "backend" stats and present them in the *de facto*
"standard".

### Installation

See the Releases tab on GitHub.

### Usage

Just run it. It supports the following flags:

`-http=<address>` The HTTP interface cw-stats will bind to.
`-location="US"` The location of the servers we'll report.
`-mode="default"` The game mode
`-name="CW Server"` Name of Server
`-poll=10s` How often we'll query the `-query` option below
`-query="http://path/to/backend"` Backend query URL
`-serverip="cw.gcg.io"` The Cube Wold public server address
