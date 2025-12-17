package kira

import (
	"fmt"
	"runtime"
	"time"
)

// Middleware handler.
func defaultPanic(ctx *Context, err any) {
	ctx.Log().Errorf("%s", err)

	if ctx.WantsJSON() {
		ctx.Response().Header().Set("Content-Type", "application/json")
	} else {
		ctx.Response().Header().Set("Content-Type", "text/html")
	}

	if ctx.WantsJSON() {
		frames := getFrames(200)
		kiraErr, ok := err.(*Error)
		if ok {
			ctx.Status(int(kiraErr.Status))
			kiraErr.Timestamp = time.Now().UTC()
			if ctx.Config().GetBool("app.debug", false) {
				kiraErr.Frames = frames
			}
			ctx.JSON(err)
			return
		}

		var generalError struct {
			Error  string          `json:"error"`
			Frames []runtime.Frame `json:"frames"`
		}
		if ctx.Config().GetBool("app.debug", false) {
			generalError.Frames = frames
		}
		generalError.Error = fmt.Sprintf("%s", err)

		ctx.JSON(generalError)
	} else {
		if ctx.ViewExists("errors/500") {
			ctx.View("errors/500")
		} else {
			ctx.WriteString(`<html><head><title>Internal Server Error</title></head><body>We're sorry, but something went wrong.</body></html>`)
		}
	}
}

func getFrames(limit int) (framesSlice []runtime.Frame) {
	// Ask runtime.Callers for up to 10 pcs, including runtime.Callers itself.
	pc := make([]uintptr, limit)
	// TODO: later we need to hide unnecessary callers.
	n := runtime.Callers(0, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()
		framesSlice = append(framesSlice, frame)
		if !more {
			break
		}
	}

	return framesSlice
}
