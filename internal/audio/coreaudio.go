//go:build darwin

package audio

/*
#cgo darwin LDFLAGS: -framework AudioToolbox -framework CoreAudio
#include <AudioToolbox/AudioToolbox.h>
#include <CoreAudio/CoreAudio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

extern void goAQInputCallback(uintptr_t handle, void* audio, uint32_t numFrames);

// AudioQueue input callback: forwards float32 mono samples to Go and re-enqueues the buffer.
static void inputCallback(void*               inUserData,
                          AudioQueueRef       inAQ,
                          AudioQueueBufferRef inBuffer,
                          const AudioTimeStamp* inStartTime,
                          UInt32              inNumPackets,
                          const AudioStreamPacketDescription* inPacketDesc)
{
	(void)inStartTime; (void)inPacketDesc; (void)inNumPackets;
	uint32_t frames = (uint32_t)(inBuffer->mAudioDataByteSize / sizeof(float));
	if (frames > 0 && inBuffer->mAudioData != NULL) {
		goAQInputCallback((uintptr_t)inUserData, inBuffer->mAudioData, frames);
	}
	AudioQueueEnqueueBuffer(inAQ, inBuffer, 0, NULL);
}

// Starts an AudioQueue capturing mono float32 at the requested sample rate.
// Returns non-zero OSStatus on error.
static OSStatus startInput(uintptr_t handle, Float64 sampleRate, int bufferFrames, int numBuffers, AudioQueueRef* outAQ) {
	AudioStreamBasicDescription fmt;
	memset(&fmt, 0, sizeof(fmt));
	fmt.mSampleRate       = sampleRate;
	fmt.mFormatID         = kAudioFormatLinearPCM;
	fmt.mFormatFlags      = kAudioFormatFlagIsFloat | kAudioFormatFlagIsPacked;
	fmt.mFramesPerPacket  = 1;
	fmt.mChannelsPerFrame = 1;
	fmt.mBitsPerChannel   = 32;
	fmt.mBytesPerPacket   = 4;
	fmt.mBytesPerFrame    = 4;

	OSStatus err = AudioQueueNewInput(&fmt, inputCallback, (void*)handle, NULL, NULL, 0, outAQ);
	if (err) return err;

	UInt32 bufSize = (UInt32)(bufferFrames * fmt.mBytesPerFrame);
	for (int i = 0; i < numBuffers; i++) {
		AudioQueueBufferRef buf = NULL;
		err = AudioQueueAllocateBuffer(*outAQ, bufSize, &buf);
		if (err) return err;
		buf->mAudioDataByteSize = bufSize;
		memset(buf->mAudioData, 0, bufSize);
		err = AudioQueueEnqueueBuffer(*outAQ, buf, 0, NULL);
		if (err) return err;
	}

	err = AudioQueueStart(*outAQ, NULL);
	if (err) return err;

	return 0;
}

static void stopInput(AudioQueueRef aq) {
	if (!aq) return;
	AudioQueueStop(aq, true);
	AudioQueueDispose(aq, true);
}
*/
import "C"

import (
	"context"
	"log"
	"runtime/cgo"
	"sync/atomic"
	"unsafe"
)

type streamState struct {
	out    chan<- []float32
	closed atomic.Bool
	ctx    context.Context
}

//export goAQInputCallback
func goAQInputCallback(h C.uintptr_t, audio unsafe.Pointer, numFrames C.uint32_t) {
	handle := cgo.Handle(h)
	v := handle.Value()
	st, ok := v.(*streamState)
	if !ok || st == nil {
		return
	}

	// If stopped or context canceled, drop immediately.
	if st.closed.Load() || st.ctx.Err() != nil {
		return
	}

	n := int(numFrames)
	if n <= 0 || audio == nil {
		return
	}

	src := unsafe.Slice((*C.float)(audio), n)
	buf := make([]float32, n)
	for i := 0; i < n; i++ {
		buf[i] = float32(src[i])
	}

	// Non-blocking send. Channel MUST NOT be closed elsewhere.
	select {
	case st.out <- buf:
	default:
		// drop to avoid stalling the realtime thread
	}
}

func Capture(ctx context.Context, out chan<- []float32) {
	log.Println("audio: starting input")
	const (
		targetSampleRate = 16000.0
		bufferFrames     = 1600 // ~100ms @ 16kHz
		numBuffers       = 3
	)

	st := &streamState{out: out, ctx: ctx}
	handle := cgo.NewHandle(st)
	defer handle.Delete()

	var aq C.AudioQueueRef
	if err := C.startInput(C.uintptr_t(handle), C.Float64(targetSampleRate), C.int(bufferFrames), C.int(numBuffers), &aq); err != 0 {
		log.Printf("audio: failed to start input (OSStatus=%d)", int32(err))
		close(out)
		return
	}

	// Wait for cancellation.
	<-ctx.Done()

	// Stop input synchronously, then close channel.
	C.stopInput(aq)
	st.closed.Store(true)
	close(out)
	log.Println("audio: stopped input")
}
