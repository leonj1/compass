package routes

import "github.com/leonj1/compass/services"

type App struct {
	Compass services.Compass
	Version string
}

const ContentType = "Content-Type"
const JSON = "application/json"
const TEXT = "text/plain; charset=us-ascii"
const POST = "POST"
const PUT = "PUT"
const GET = "GET"
const DELETE = "DELETE"
const HEAD = "HEAD"
const OPTIONS = "OPTIONS"
