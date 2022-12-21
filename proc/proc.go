package proc

import (
	"errors"
	"fmt"
	"github.com/dewey/tbm/config"
	"github.com/dewey/tbm/log"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ServicesService is holding all the information to manage an in-memory process register
type ServicesService struct {
	// Configuration is the configuration read from the yaml file
	Configuration config.Configuration
	// maxProcNameLength is the longest name of a proc. This is used to align the console output properly.
	maxProcNameLength int
	// procs is the in-memory representation of all currently running processes
	procs []*ProcInfo
	mu    sync.Mutex
}

// NewServicesService returns a new services service
func NewServicesService(cfg config.Configuration) *ServicesService {
	return &ServicesService{
		Configuration:     cfg,
		maxProcNameLength: 0,
		procs:             []*ProcInfo{},
		mu:                sync.Mutex{},
	}
}

var colors = []int{
	32, // green
	36, // cyan
	35, // magenta
	33, // yellow
	34, // blue
	31, // red
}

// ProcInfo defines the structure of a single process
type ProcInfo struct {
	name        string
	environment string
	cmdline     string
	cmd         *exec.Cmd
	port        uint
	setPort     bool
	colorIndex  int

	// True if we called stopProc to kill the process, in which case an
	// *os.ExitError is not the fault of the subprocess
	stoppedBySupervisor bool

	mu      sync.Mutex
	cond    *sync.Cond
	waitErr error
}

// Procs returns all initialized procs
func (svc *ServicesService) Procs() []*ProcInfo {
	return svc.procs
}

// FindProc returns a single proc object from the in-memory object, selected by name
func (svc *ServicesService) FindProc(name string) *ProcInfo {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	for _, proc := range svc.procs {
		if proc.name == name {
			return proc
		}
	}
	return nil
}

func (p *ProcInfo) ClearName() string {
	return ProcClearName(p.name, p.environment)
}

// ProcClearName returns the clear service name of a proc
func ProcClearName(name string, environment string) string {
	return strings.Replace(name, "-"+environment, "", -1)
}

// spawnProc starts the specified proc, and returns any error from running it.
func (svc *ServicesService) spawnProc(name string, errCh chan<- error) {
	cproc := svc.FindProc(name)
	logger := log.New(name, cproc.environment, cproc.colorIndex, svc.maxProcNameLength)

	cs := append(cmdStart, cproc.cmdline)
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = logger
	cmd.Stderr = logger
	cmd.SysProcAttr = procAttrs

	if cproc.setPort {
		fmt.Fprintf(logger, "Starting %s on port %d\n", cproc.ClearName(), cproc.port)
	}
	if err := cmd.Start(); err != nil {
		select {
		case errCh <- err:
		default:
		}
		fmt.Fprintf(logger, "Failed to start %s: %s\n", name, err)
		return
	}
	cproc.cmd = cmd
	cproc.stoppedBySupervisor = false
	cproc.mu.Unlock()
	err := cmd.Wait()
	cproc.mu.Lock()
	cproc.cond.Broadcast()
	if err != nil && !cproc.stoppedBySupervisor {
		select {
		case errCh <- err:
		default:
		}
	}
	cproc.waitErr = err
	cproc.cmd = nil
	fmt.Fprintf(logger, "Terminating %s\n", name)
}

// stopProc is stopping the specified process. Issuing os.Kill if it does not terminate within 10 seconds. If signal is
// nil, os.Interrupt is used.
func (svc *ServicesService) stopProc(name string, signal os.Signal) error {
	if signal == nil {
		signal = os.Interrupt
	}
	proc := svc.FindProc(name)
	if proc == nil {
		return errors.New("unknown proc: " + name)
	}

	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.cmd == nil {
		return nil
	}
	proc.stoppedBySupervisor = true

	err := terminateProc(proc, signal)
	if err != nil {
		return err
	}

	timeout := time.AfterFunc(10*time.Second, func() {
		proc.mu.Lock()
		defer proc.mu.Unlock()
		if proc.cmd != nil {
			err = killProc(proc.cmd.Process)
		}
	})
	proc.cond.Wait()
	timeout.Stop()
	return err
}

// ReadProcfile reads a configuration object and stores it in the global, in-memory object used to keep track of it.
func (svc *ServicesService) ReadProcfile(cfg config.Configuration) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.procs = []*ProcInfo{}
	index := 0
	for key, service := range cfg.Services {
		// Skip all the services that don't pass the validation (Not enabled, erroneous configuration etc.)
		if !service.Valid() {
			continue
		}
		// Create proc based on configuration
		cmd, err := service.InterpolatedCommand()
		if err != nil {
			return err
		}

		proc := &ProcInfo{
			name:        fmt.Sprintf("%s-%s", key, service.Environment),
			environment: service.Environment,
			cmdline:     cmd,
			colorIndex:  index,
		}
		exists, val := service.VariableValue("port")
		if exists {
			i, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			proc.port = uint(i)
			proc.setPort = true
		}
		proc.cond = sync.NewCond(&proc.mu)
		svc.procs = append(svc.procs, proc)
		index = (index + 1) % len(colors)
	}

	if len(svc.procs) > svc.maxProcNameLength {
		svc.maxProcNameLength = len(svc.procs)
	}
	if len(svc.procs) == 0 {
		return errors.New("no valid service entry in configuration file")
	}
	return nil
}

// startProc a specified proc by name. If proc is started already, return nil.
func (svc *ServicesService) startProc(name string, wg *sync.WaitGroup, errCh chan<- error) error {
	proc := svc.FindProc(name)
	if proc == nil {
		return errors.New("unknown name: " + name)
	}

	proc.mu.Lock()
	if proc.cmd != nil {
		proc.mu.Unlock()
		return nil
	}

	if wg != nil {
		wg.Add(1)
	}
	go func() {
		svc.spawnProc(name, errCh)
		if wg != nil {
			wg.Done()
		}
		proc.mu.Unlock()
	}()
	return nil
}

// stopProcs attempts to stop every running process and returns any non-nil
// error, if one exists. stopProcs will wait until all procs have had an
// opportunity to stop.
func (svc *ServicesService) stopProcs(sig os.Signal) error {
	var err error
	for _, proc := range svc.procs {
		stopErr := svc.stopProc(proc.name, sig)
		if stopErr != nil {
			err = stopErr
		}
	}
	return err
}

// StartProcs starts all procs in separate go routines
func (svc *ServicesService) StartProcs(sc <-chan os.Signal, exitOnError bool, exitOnStop bool) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, proc := range svc.procs {
		svc.startProc(proc.name, &wg, errCh)
	}

	allProcsDone := make(chan struct{}, 1)
	if exitOnStop {
		go func() {
			wg.Wait()
			allProcsDone <- struct{}{}
		}()
	}
	for {
		select {
		case err := <-errCh:
			if exitOnError {
				svc.stopProcs(os.Interrupt)
				return err
			}
		case <-allProcsDone:
			return svc.stopProcs(os.Interrupt)
		case sig := <-sc:
			return svc.stopProcs(sig)
		}
	}
}

const sigint = unix.SIGINT
const sigterm = unix.SIGTERM
const sighup = unix.SIGHUP

var cmdStart = []string{"/bin/sh", "-c"}
var procAttrs = &unix.SysProcAttr{Setpgid: true}

// terminateProc is killing a proc by pid
func terminateProc(proc *ProcInfo, signal os.Signal) error {
	p := proc.cmd.Process
	if p == nil {
		return nil
	}

	pgid, err := unix.Getpgid(p.Pid)
	if err != nil {
		return err
	}

	// use pgid, ref: http://unix.stackexchange.com/questions/14815/process-descendants
	pid := p.Pid
	if pgid == p.Pid {
		pid = -1 * pid
	}

	target, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return target.Signal(signal)
}

// killProc kills the proc with pid, as well as its children.
func killProc(process *os.Process) error {
	return unix.Kill(-1*process.Pid, unix.SIGKILL)
}

func NotifyCh() <-chan os.Signal {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, sigterm, sigint, sighup)
	return sc
}
