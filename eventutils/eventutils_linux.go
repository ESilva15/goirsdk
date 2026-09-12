//go:build linux && cgo

package mmaputils

/*
#include <semaphore.h>
#include <fcntl.h>
#include <time.h>
#include <errno.h>
#include <stdlib.h>

static inline void* open_posix_semaphore(const char* name) {
    sem_t* sem = sem_open(name, O_CREAT, 0666, 0);
    if (sem == SEM_FAILED) {
        return NULL;
    }
    return (void*)sem;
}

static inline void close_posix_semaphore(void* sem_ptr) {
    if (sem_ptr != NULL) {
        sem_close((sem_t*)sem_ptr);
    }
}

static inline int timed_wait_posix_semaphore(void* sem_ptr, long timeout_ms) {
		if (sem_ptr == NULL) {
        return -1;
    }
    sem_t* sem = (sem_t*)sem_ptr;

    struct timespec ts;
    clock_gettime(CLOCK_REALTIME, &ts);

    ts.tv_sec += timeout_ms / 1000;
    ts.tv_nsec += (timeout_ms % 1000) * 1000000;

    if (ts.tv_nsec >= 1000000000) {
        ts.tv_sec += ts.tv_nsec / 1000000000;
        ts.tv_nsec %= 1000000000;
    }

    int res;
    do {
        res = sem_timedwait(sem, &ts);
    } while (res == -1 && errno == EINTR); // Retry if interrupted by Go GC/scheduler

    return res;
}

static inline int signal_posix_semaphore(const char* name) {
    if (name == NULL) {
        return -1;
    }
    sem_t* sem = sem_open(name, 0);
    if (sem == SEM_FAILED) {
        return -1;
    }
    int res = sem_post(sem);
    sem_close(sem);
    return res;
}

static inline int post_posix_semaphore(void* sem_ptr) {
	if (sem_ptr != NULL) {
		sem_post((sem_t*)sem_ptr);
	}
	return -1;
}

static inline int unlink_posix_semaphore(const char* name) {
    if (name == NULL) {
        return -1;
    }
    return sem_unlink(name);
}
*/
import "C"

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/ESilva15/goirsdk/sharedMem"
)

type utils struct {
	semName string
	sem     unsafe.Pointer
}

func newUtils() (*utils, error) {
	return &utils{}, nil
}

func (u *utils) Close() {
	if u.sem != nil {
		C.close_posix_semaphore(u.sem)
		u.sem = nil
	}
}

// OpenMemMap returns a Reader interface that can be used to read the data
// No need to encapsulate it
func OpenMemMap(name string, size uint32) (Reader, error) {
	file, err := sharedMem.Open(name, size)
	if err != nil {
		return nil, fmt.Errorf("failed to open `%s` with err: %+v", name, err)
	}

	return file, nil
}

// OpenEvent creates or connects to a Unix socket for event signaling on Linux
func (u *utils) OpenEvent(eventName string) error {
	u.semName = "/" + strings.TrimPrefix(eventName, "/")

	cName := C.CString(u.semName)
	defer C.free(unsafe.Pointer(cName))

	sem := C.open_posix_semaphore(cName)
	if sem == nil {
		return fmt.Errorf("failed to open POSIX semaphore: %s", u.semName)
	}

	u.sem = sem
	return nil
}

func (u *utils) OpenBroadcastChannel(name string) error {
	// No-op or log stub on Linux
	return nil
}

// CheckValidDataEvent waits for a pulse byte sent over the socket
func (u *utils) CheckValidDataEvent(timeout time.Duration) bool {
	if u.sem == nil {
		return false
	}

	ms := C.long(timeout.Milliseconds())
	res := C.timed_wait_posix_semaphore(u.sem, ms)

	return res == 0
}

func (u *utils) SendBroadcastMessage(id, p1, p2 uintptr) error {
	return nil
}

func (u *utils) SignalEvent() error {
	if u.sem == nil {
		return fmt.Errorf("u.sem is not set")
	}

	C.post_posix_semaphore(u.sem)
	return nil
}

func signalEvent(name string) {
	semName := "/" + strings.TrimPrefix(name, "/")
	cName := C.CString(semName)
	defer C.free(unsafe.Pointer(cName))
	C.signal_posix_semaphore(cName)
}

func cleanupEvent(name string) {
	semName := "/" + strings.TrimPrefix(name, "/")
	cName := C.CString(semName)
	defer C.free(unsafe.Pointer(cName))
	C.unlink_posix_semaphore(cName)
}
