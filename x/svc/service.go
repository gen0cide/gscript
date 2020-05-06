package svc

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/gen0cide/privcheck"
	"github.com/kardianos/service"
	"github.com/phayes/permbits"
	"github.com/pkg/errors"
)

var (
	actions = map[string]bool{
		"start":     true,
		"stop":      true,
		"restart":   true,
		"install":   true,
		"uninstall": true,
	}

	osxOptions = map[string]interface{}{
		"KeepAlive": true,
		"RunAtLoad": false,
	}

	supportedSystems = map[string]configLookupFunc{
		"windows-service": windowsLookup,
		"linux-systemd":   systemdLookup,
		"linux-upstart":   upstartLookup,
		"unix-systemv":    systemvLookup,
		"darwin-launchd":  launchdLookup,
	}

	// ErrSystemUnsupported is returned when a function call cannot identify the
	// current system's supervisor
	ErrSystemUnsupported = errors.New("no supported system supervisor located")

	// ErrExeNotExist is returned when the service cannot find it's set executable
	ErrExeNotExist = errors.New("no file located at executable path")

	// ErrExeNotExecutable is returned when an executable path contains a file without valid execute perms
	ErrExeNotExecutable = errors.New("file at executable path is not executable")

	// ErrNoName is returned when a service is attempted to be enumerated without a name
	ErrNoName = errors.New("service does not have a name")

	// ErrNoDescription is returned when a service does not have a description set
	ErrNoDescription = errors.New("service does not have a description")

	// ErrNoDisplayName is returned when a service does not have a display name set
	ErrNoDisplayName = errors.New("service does not have a display name")

	// ErrNotEnoughPerms is returned when the current process does not have sufficient privileges to install a service
	ErrNotEnoughPerms = errors.New("current user does not have privileges to manipulate services")

	// ErrWorkingDirNotExist is returned when the working directory path does not exist
	ErrWorkingDirNotExist = errors.New("working directory does not exist")
)

func windowsLookup(s *Service) (string, error) {
	return "", errors.New("windows does not have config files for services")
}

func systemdLookup(s *Service) (string, error) {
	name := s.Name
	if s.Name == "" {
		return "", ErrNoName
	}
	return filepath.Join("/etc/systemd/system", fmt.Sprintf("%s.service", name)), nil
}

func upstartLookup(s *Service) (string, error) {
	name := s.Name
	if s.Name == "" {
		return "", ErrNoName
	}
	return filepath.Join("/etc/init", fmt.Sprintf("%s.conf", name)), nil
}

func systemvLookup(s *Service) (string, error) {
	name := s.Name
	if s.Name == "" {
		return "", ErrNoName
	}
	return filepath.Join("/etc/init.d", name), nil
}

func launchdLookup(s *Service) (string, error) {
	name := s.Name
	if s.Name == "" {
		return "", ErrNoName
	}
	if s.Options.GetBool("UserService") {
		hd := ""
		u, err := user.Current()
		if err == nil {
			hd = u.HomeDir
		}
		if hd == "" {
			hd := os.Getenv("HOME")
			if hd == "" {
				return "", errors.New("unable to determine home directory of current user")
			}
		}
		return filepath.Join(hd, "Library", "LaunchAgents", fmt.Sprintf("%s.plist", name)), nil
	}
	return filepath.Join("/Library/LaunchDaemons", fmt.Sprintf("%s.plist", name)), nil
}

type configLookupFunc func(s *Service) (string, error)

// Service represents a definition for a system service (cross platform)
type Service struct {
	Name             string   `json:"name,omitempty"`
	DisplayName      string   `json:"display_name,omitempty"`
	Description      string   `json:"description,omitempty"`
	UserName         string   `json:"user_name,omitempty"`
	Arguments        []string `json:"arguments,omitempty"`
	ExecutablePath   string   `json:"executable_path,omitempty"`
	WorkingDirectory string   `json:"working_directory,omitempty"`
	Options          KeyValue `json:"options"`
}

// KeyValue is a helper type alias for platform specific options
type KeyValue map[string]interface{}

// GetBool will return a true when a value exists and is anything other a boolean "false".
// Otherwise it will return false.
func (k KeyValue) GetBool(key string) bool {
	val, ok := k[key]
	if !ok {
		return false
	}
	bv, ok := val.(bool)
	if !ok {
		return true
	}
	return bv
}

// NewAsZero instantiates a zero value Service object
func NewAsZero() *Service {
	return &Service{}
}

// NewFromArgs instantiates a Service object given the parameters
func NewFromArgs(name, displayName, description, exePath string, args []interface{}, opts map[string]interface{}) (*Service, error) {
	svc := &Service{
		Name:           name,
		DisplayName:    displayName,
		Description:    description,
		ExecutablePath: exePath,
		Arguments:      []string{},
		Options:        KeyValue(opts),
	}
	svc.AppendArguments(args)
	return svc, nil
}

// NewFromJSON instantiates a Service object from a given map of parameters
func NewFromJSON(jsonData map[string]interface{}) (*Service, error) {
	var svc Service
	err := json.Unmarshal(gabs.Wrap(jsonData).Bytes(), &svc)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal JSON argument")
	}
	return &svc, nil
}

// ServicePlatform returns a string of the current system's service platform
func ServicePlatform() string {
	if len(service.AvailableSystems()) == 0 {
		return ""
	}
	return service.ChosenSystem().String()
}

// SetName sets the service name
func (s *Service) SetName(name string) {
	s.Name = name
}

// SetDisplayName sets the service longer display name
func (s *Service) SetDisplayName(displayName string) {
	s.DisplayName = displayName
}

// SetDescription sets the service description
func (s *Service) SetDescription(description string) {
	s.Description = description
}

// SetUserName sets the user to run the service as
// Note: not supported on darwin-launchd, unix-systemv, windows-service requires a 'Password' option value.
func (s *Service) SetUserName(username string) {
	s.UserName = username
}

// SetArguments clears the service's defined arguments and re-defines them given the passed args.
func (s *Service) SetArguments(args []interface{}) {
	s.Arguments = []string{}
	s.AppendArguments(args)
}

// AppendArguments does not clear the current argument list, simply adds to it.
func (s *Service) AppendArguments(args []interface{}) {
	for _, a := range args {
		s.Arguments = append(s.Arguments, fmt.Sprintf("%v", a))
	}
}

// SetExecutablePath sets the location where the executable binary will be on the filesystem.
// Note - should ensure executable has valid execute permissions on the filesystem.
func (s *Service) SetExecutablePath(exepath string) {
	s.ExecutablePath = exepath
}

// SetWorkingDirectory allows you to tell the platform (where supported) what directory
// the program should be located in while running.
func (s *Service) SetWorkingDirectory(dir string) {
	s.WorkingDirectory = dir
}

// SetOptions allows you to overwrite the current options and it's values.
func (s *Service) SetOptions(opts map[string]interface{}) {
	s.Options = KeyValue(opts)
}

// CheckConfig looks to make sure the options and information you've defined
// in your Service object checks out for the current system.
func (s *Service) CheckConfig(autofix bool) (bool, error) {
	// verify that we are on a compatible system
	if ServicePlatform() == "" {
		return false, ErrSystemUnsupported
	}

	// verify that we are indeed an administrator
	if !privcheck.IsAdmin() {
		return false, ErrNotEnoughPerms
	}

	// verify that the executable exists and that it's not a directory
	fs, err := os.Stat(s.ExecutablePath)
	if err != nil && os.IsNotExist(err) {
		return false, ErrExeNotExist
	} else if err != nil {
		return false, err
	}

	if fs.IsDir() {
		return false, ErrExeNotExist
	}

	// verify the permissions on the executable and fix if they're wrong
	perms := permbits.FileMode(fs.Mode())
	permsChanged := false
	if !perms.UserExecute() {
		if !autofix {
			return false, ErrExeNotExecutable
		}
		permsChanged = true
		perms.SetUserExecute(true)
	}
	if !perms.GroupExecute() {
		if !autofix {
			return false, ErrExeNotExecutable
		}
		permsChanged = true
		perms.SetGroupExecute(true)
	}
	if !perms.OtherExecute() {
		if !autofix {
			return false, ErrExeNotExecutable
		}
		permsChanged = true
		perms.SetOtherExecute(true)
	}

	if permsChanged {
		err = permbits.Chmod(s.ExecutablePath, perms)
		if err != nil {
			return false, errors.Wrap(err, "could not set permissions of executable")
		}
	}

	// verify that the program's main attributes (name, description, display_name) are present
	if s.Name == "" {
		return false, ErrNoName
	}

	if s.Description == "" {
		return false, ErrNoDescription
	}

	if s.DisplayName == "" {
		return false, ErrNoDisplayName
	}

	if s.WorkingDirectory == "" && ServicePlatform() == "windows-service" {
		return true, nil
	}

	// verify that working directory actually exists and make it if it doesn't (autofix)
	err = nil
	_, err = os.Stat(s.WorkingDirectory)
	if err != nil && os.IsNotExist(err) {
		if !autofix {
			return false, ErrWorkingDirNotExist
		}
		mkdirErr := os.MkdirAll(s.WorkingDirectory, 0700)
		if mkdirErr != nil {
			return false, errors.WithMessage(err, ErrWorkingDirNotExist.Error())
		}
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// IsErrExeNotExist is used to check if a returned error is because the executable doesn't exist
func IsErrExeNotExist(e error) bool {
	return e == ErrExeNotExist
}

// IsErrExeNotExecutable is used to check if a returned error is because the executable has insufficient permissions
func IsErrExeNotExecutable(e error) bool {
	return e == ErrExeNotExecutable
}

// IsErrNoName is used to check if a returned error is because no name is set
func IsErrNoName(e error) bool {
	return e == ErrNoName
}

// IsErrNoDescription is used to check if a returned error is because no description is set
func IsErrNoDescription(e error) bool {
	return e == ErrNoDescription
}

// IsErrNoDisplayName is used to check if a returned error is because no display name is set
func IsErrNoDisplayName(e error) bool {
	return e == ErrNoDisplayName
}

// IsErrSystemUnsupported is used to check if a returned error is because the current system does not have a supported supervisor
func IsErrSystemUnsupported(e error) bool {
	return e == ErrSystemUnsupported
}

// IsErrNotEnoughPerms is used to check if a returned error is because the current process lacks privileges to manipulate system services
func IsErrNotEnoughPerms(e error) bool {
	return e == ErrNotEnoughPerms
}

// IsErrWorkingDirNotExist is used to check if a returned error is because the defined working directory does not exist
func IsErrWorkingDirNotExist(e error) bool {
	return e == ErrWorkingDirNotExist
}

// CheckExists will attempt to check to see if the service
// as defined already exists on the system.
// Note: this currently does not work on Windows since there is
// no file location to check.
func (s *Service) CheckExists() (bool, error) {
	fileloc, err := s.ServiceFilePath()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(fileloc); err != nil && os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

// ServiceFilePath returns the path to config file (where applicable) for the
// defined configuration.
func (s *Service) ServiceFilePath() (string, error) {
	fn, ok := supportedSystems[service.ChosenSystem().String()]
	if !ok {
		return "", ErrSystemUnsupported
	}
	return fn(s)
}

// Install attempts to install the defined system service
func (s *Service) Install(overwrite bool) error {
	svc, err := s.translate()
	if err != nil {
		return err
	}
	if overwrite {
		svc.Stop()
		time.Sleep(20 * time.Millisecond)
		svc.Uninstall()
		time.Sleep(20 * time.Millisecond)
		if ServicePlatform() == "unix-systemv" {
			s.removesystemvsyms()
		}
		err = s.removeConfigFile()
		if err != nil {
			return err
		}
	}
	return svc.Install()
}

// Stop attempts to stop the running service
func (s *Service) Stop() error {
	svc, err := s.translate()
	if err != nil {
		return err
	}
	return svc.Stop()
}

// Start attempts to start the running service
func (s *Service) Start() error {
	svc, err := s.translate()
	if err != nil {
		return err
	}
	return svc.Start()
}

// Restart attempts to restart a service with a delay
// of 50 milliseconds.
func (s *Service) Restart() error {
	svc, err := s.translate()
	if err != nil {
		return err
	}
	return svc.Restart()
}

// Uninstall attempts to remove a system service
// from the local system including:
// - stopping the service
// - uninstalling the service
// - removing symlinks to any scripts
// - removing config files
// - reloading platform service supervisors when applicable
func (s *Service) Uninstall() error {
	svc, err := s.translate()
	if err != nil {
		return err
	}
	err = svc.Uninstall()
	if err != nil {
		return err
	}
	s.removeConfigFile()
	return nil
}

// RestartWithDelay will stop a service, wait for a given period of time,
// then attempt to start the service after the delay period.
func (s *Service) RestartWithDelay(seconds int64) error {
	svc, err := s.translate()
	if err != nil {
		return err
	}
	svc.Stop()
	time.Sleep(time.Duration(seconds) * time.Second)
	return svc.Start()
}

func (s *Service) translate() (service.Service, error) {
	conf := service.Config{
		Name:             s.Name,
		DisplayName:      s.DisplayName,
		Description:      s.Description,
		UserName:         s.UserName,
		Arguments:        s.Arguments,
		Executable:       s.ExecutablePath,
		WorkingDirectory: s.WorkingDirectory,
		Option:           service.KeyValue(s.Options),
	}

	return service.ChosenSystem().New(nil, &conf)
}

func (s *Service) removeConfigFile() error {
	switch ServicePlatform() {
	case "windows-service":
		return nil
	case "unix-systemv":
		s.removesystemvsyms()
	}

	exists, err := s.CheckExists()
	if err != nil {
		return err
	}
	if exists {
		fp, _ := s.ServiceFilePath()
		os.Remove(fp)
	}
	return nil
}

func (s *Service) removesystemvsyms() {
	for _, i := range []int{2, 3, 4, 5} {
		os.Remove(filepath.Join("/etc", fmt.Sprintf("rc%d.d", i), fmt.Sprintf("S50%s", s.Name)))
	}
	for _, i := range []int{0, 1, 6} {
		os.Remove(filepath.Join("/etc", fmt.Sprintf("rc%d.d", i), fmt.Sprintf("K02%s", s.Name)))
	}
}
