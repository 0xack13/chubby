package main

import (
	"cos518project/chubby/api"
	"cos518project/chubby/client"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var clientID1		string		// ID of client 1
var clientID2		string		// ID of client 2

func init() {
	flag.StringVar(&clientID1, "clientID1", "simple_client_1", "ID of client 1")
	flag.StringVar(&clientID2, "clientID2", "simple_client_2", "ID of client 2")
}

func main() {
	// Parse flags from command line.
	flag.Parse()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Establish two sessions
	sess1, err := client.InitSession(api.ClientID(clientID1))
	sess2, err := client.InitSession(api.ClientID(clientID2))

	// Test Open Locks
	errOpenLock1 := sess1.OpenLock("LOCK/Lock1")
	errOpenLock2 := sess2.OpenLock("LOCK/Lock2")

	if errOpenLock1 != nil {
		log.Printf("Session 1 has trouble opening lock ")
		log.Fatal(errOpenLock1)
	} else {
		log.Printf("Session 1 has opened lock successfully")
	}
	if errOpenLock2 != nil {
		log.Printf("Session 2 has trouble opening lock")
		log.Fatal(errOpenLock2)
	} else {
		log.Printf("Session 2 has opened lock successfully")
	}

	errOpenLock1 = sess1.OpenLock("LOCK/LockShared")
	if errOpenLock1 != nil {
		log.Printf("Session 1 has trouble opening lock")
		log.Fatal(errOpenLock1)
	} else {
		log.Printf("Session 1 has opened lock successfully")
	}

	// Test TryAcquire Lock
	isSuccessful, acquireErr := sess1.TryAcquireLock("LOCK/Lock1", api.EXCLUSIVE)
	if !isSuccessful {
		log.Printf("Try Acquire Lock failed when it should succeed")
	}
	if acquireErr != nil {
		log.Printf("Try Acquire Lock Unexpected Error")
		log.Fatal(acquireErr)
	}

	// Try Acquire a Shared Lock
	isSuccessful, acquireErr = sess1.TryAcquireLock("LOCK/LockShared", api.SHARED)
	if !isSuccessful {
		log.Printf("Try Acquire Shared Lock failed when it should succeed")
	}
	if acquireErr != nil {
		log.Printf("Try Acquire Lock Unexpected Error")
		log.Fatal(acquireErr)
	}

	isSuccessful, acquireErr = sess2.TryAcquireLock("LOCK/LockShared", api.SHARED)
	if !isSuccessful {
		log.Printf("Try Acquire Shared Lock failed when it should succeed")
	}
	if acquireErr != nil {
		log.Printf("Try Acquire Shared Lock Unexpected Error")
		log.Fatal(acquireErr)
	}

	// Should not be able to acquire a lock you already acquired
	isSuccessful, acquireErr = sess1.TryAcquireLock("LOCK/Lock1", api.EXCLUSIVE)
	if isSuccessful {
		log.Printf("Should fail because the lock we are trying to acquire is in exclusive mode")
	}
	if acquireErr == nil {
		log.Printf("Should fail because the lock we are trying to acquire is in exclusive mode")
	}

	// Should not be able to acquire a lock someone else acquired in exclusive mode
	isSuccessful, acquireErr = sess2.TryAcquireLock("LOCK/Lock1", api.EXCLUSIVE)
	if isSuccessful {
		log.Printf("Session 2 Should fail because the lock we are trying to acquire is in exclusive mode")
	}
	if acquireErr == nil {
		log.Printf("Session 2 Should fail because the lock we are trying to acquire is in exclusive mode")
	}

	// Should not be able to release a lock you don't own
	releaseErr := sess2.ReleaseLock("LOCK/Lock1")
	if releaseErr == nil {
		log.Printf("Should fail because the lock we are trying to release is a lock we don't own")
	}

	// Test release lock
	releaseErr = sess1.ReleaseLock("LOCK/Lock1")
	if releaseErr != nil {
		log.Printf("Unexpected Lock release failure")
		log.Fatal(releaseErr)
	}

	// Test Delete Lock


	if err != nil {
		log.Fatal(err)
	}

	// Exit on signal.
	<-quitCh
}
