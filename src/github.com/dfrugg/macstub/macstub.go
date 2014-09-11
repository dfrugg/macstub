package main

import (
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "log"
    "log/syslog"
)

func checkForError(err error){
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}


func main(){
	syslogger, err := syslog.New(syslog.LOG_ERR, "[]")
	defer syslogger.Close()
	checkForError(err)
	
	syslogger.Info("Entering JavaApplicationStub.\n")
	for _, arg := range os.Args{
		syslogger.Info(arg)
		syslogger.Info("\n")
	}
	
	// Determine Location
	macOSPath := filepath.Dir(os.Args[0])
	syslogger.Info("MacOS Path: " + macOSPath + "\n")
	
	contentsPath := filepath.Dir(macOSPath)
	syslogger.Info("Contents Path: " + contentsPath + "\n")
	
	appPath := filepath.Dir(contentsPath)
	syslogger.Info("Application Path: " + appPath + "\n")
	
	plistPath := filepath.Clean(contentsPath + "/Info.plist")
	syslogger.Info("Info.plist Path: " + plistPath + "\n")
	
	// Read Info.plist
	plistFile, err := os.Open(plistPath)
	checkForError(err)
	
	props, err := plistToMap(plistFile)
	checkForError(err)
	
	// Check for Application Resources
	applicationName := props["CFBundleName"]
	applicationIcon := props["CFBundleIconFile"]
	javaFolder := appPath + "/Contents/Resources/Java"
	resourceFolder := appPath + "/Contents/Resources"
	file, err := os.Open(javaFolder)
	if err != nil {
		javaFolder = appPath + "/Contents/Java"
	}
	file.Close()
	
	// Replace $APP_PACKAGE
	for key, value := range props {
		if strings.Contains(value, "$APP_PACKAGE"){
			props[key] = strings.Replace(value, "$APP_PACKAGE", appPath, -1)
		}
		syslogger.Info("Key: " + key + ", Value: " + value + "\n")
	}
		
	// Build Java Command
	javaHome := os.Getenv("JAVA_HOME");
	if javaHome == "" {
		out, err := exec.Command("/usr/libexec/java_home").Output()
		checkForError(err)
		if out == nil || len(out) == 0 {
			syslogger.Err("Unable to obtain java home.")
			os.Exit(1)
		}
		javaHome = strings.Trim(string(out), " \n\r")
	}
	javaCommand := javaHome + "/bin/java "
	
	// Build Arguments
	mainClass := props["Java|MainClass"]
	workingDirectory := props["Java|WorkingDirectory"]
	classPath := props["Java|ClassPath"] + ":" + javaFolder + "/*"
	//arguments := props["Java|Arguments"]
	//vmoptions := props["Java|VMOptions"]
	
	if classPath != "" {javaCommand = javaCommand + `-cp "` + classPath + `" `}
	javaCommand = javaCommand + ` -Xdock:icon="` + resourceFolder + "/" + applicationIcon + `"`
	javaCommand = javaCommand + ` -Xdock:name="` + applicationName + `"`
			
	for key, value := range props {
		if strings.Contains(key, "Java|Properties|"){
			prop := " -D" + strings.Replace(key, "Java|Properties|", "", 1) + `="` + value + `" `
			javaCommand = javaCommand + prop
		}
	}
	
	if mainClass != "" {javaCommand = javaCommand + mainClass + " "}
	
	syslogger.Info("Command: " + javaCommand + "\n")
	syslogger.Info("\n")
	
	// Change Directory
	err = os.Chdir(workingDirectory)
	checkForError(err)
	
	command := exec.Command("sh", "-c", javaCommand);
	command.Dir = workingDirectory
	
	err = command.Start()
	if err != nil {
		syslogger.Info("Command Error\n")
		syslogger.Info(err.Error() + "\n")
	}
	syslogger.Info("\n")
	syslogger.Info("Exiting JavaApplicationStub.\n")
}