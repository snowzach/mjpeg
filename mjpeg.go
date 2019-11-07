// Package mjpeg implements a simple MJPEG streamer.
//
// Stream objects implement the http.Handler interface, allowing to use them with the net/http package like so:
//	stream = mjpeg.NewStream()
//	http.Handle("/camera", stream)
// Then push new JPEG frames to the connected clients using stream.UpdateJPEG().
package mjpeg

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Stream represents a single video feed.
type Stream struct {
	clients       map[chan []byte]struct{}
	FrameInterval time.Duration
	sync.Mutex
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

// NewStream initializes and returns a new Stream.
func NewStream(interval time.Duration) *Stream {
	if interval == 0 {
		interval = 50 * time.Millisecond
	}
	return &Stream{
		clients:       make(map[chan []byte]struct{}),
		FrameInterval: interval,
	}
}

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	c := make(chan []byte)
	s.Lock()
	s.clients[c] = struct{}{}
	s.Unlock()

	for {
		time.Sleep(s.FrameInterval)
		b := <-c
		_, err := w.Write(b)
		if err != nil {
			break
		}
	}

	s.Lock()
	delete(s.clients, c)
	s.Unlock()
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func (s *Stream) UpdateJPEG(jpeg []byte) {
	header := fmt.Sprintf(headerf, len(jpeg))
	frame := make([]byte, (len(jpeg) + len(header)))

	copy(frame, header)
	copy(frame[len(header):], jpeg)

	s.Lock()
	for c := range s.clients {
		select {
		case c <- frame:
		default:
		}
	}
	s.Unlock()
}
