package main

import (
	"errors"
	"strings"
	"io"
	"encoding/xml"
)

func readNextEnd(decoder *xml.Decoder) (*xml.EndElement, error) {
	for token, err := decoder.Token() ; token != nil; token, err = decoder.Token() {
		checkForError(err)
		switch t1 := token.(type) {
			case xml.EndElement:
				return &t1, nil
		}
	}
	return nil, errors.New("No EndElement has been found.")
}

func readNextCharData(decoder *xml.Decoder) (*xml.CharData, error) {
	for token, err := decoder.Token() ; token != nil; token, err = decoder.Token() {
		checkForError(err)
		switch t1 := token.(type) {
			case xml.CharData:
				return &t1, nil
		}
	}
	return nil, errors.New("No Character Data has been found.")
}

func readNextStart(decoder *xml.Decoder) (*xml.StartElement, error) {
	for token, err := decoder.Token() ; token != nil; token, err = decoder.Token() {
		checkForError(err)
		switch t1 := token.(type) {
			case xml.StartElement:
				return &t1, nil
		}
	}
	return nil, errors.New("No StartElement has been found.")
}

func readNextStartUntil(decoder *xml.Decoder, localName string) (*xml.StartElement, error)  {
	for token, err := readNextStart(decoder) ; token != nil; token, err = readNextStart(decoder) {
		checkForError(err)
		if token.Name.Local == localName {
			return token, nil
		}
	}
	return nil, errors.New("No StartElement of name " + localName + " has been encountered.");
}

func readNextCharDataOrEnd(decoder *xml.Decoder) (*xml.CharData, *xml.EndElement, error)  {
	for token, err := decoder.Token() ; token != nil; token, err = decoder.Token() {
		checkForError(err)
		switch t1 := token.(type) {
			case xml.CharData:
				return &t1, nil, nil
			case xml.EndElement:
				return nil, &t1, nil
		}
	}
	return nil, nil, errors.New("No Character Data or EndElement has been found.");
}

func readNextStartOrEnd(decoder *xml.Decoder, startName string, endName string) (*xml.StartElement, *xml.EndElement, error)  {
	for token, err := decoder.Token() ; token != nil; token, err = decoder.Token() {
		checkForError(err)
		switch t1 := token.(type) {
			case xml.StartElement:
				if t1.Name.Local == startName{
					return &t1, nil, nil
				}
			case xml.EndElement:
				if t1.Name.Local == endName{
					return nil, &t1, nil
				}
		}
	}
	return nil, nil, errors.New("No StartElement of name '" + startName + 
		"' or EndElement of name '" + endName + "' has been found.");
}

func readPastElement(decoder *xml.Decoder, element string){
	for start, end, err := readNextStartOrEnd(decoder, element, element) ; start != nil; 
			start, end, err = readNextStartOrEnd(decoder, element, element){
		checkForError(err)
		if(start != nil){
			readPastElement(decoder, element)
		}
		if(end != nil){
			return
		}
	}
}

func readKeyValue(decoder *xml.Decoder, path string, values map[string]string) {
	token, err := readNextStart(decoder)
	checkForError(err)
	switch token.Name.Local {
		case "true":
			values[path] = "true"
		case "false":
			values[path] = "false"
		case "string":
			data, _, err := readNextCharDataOrEnd(decoder)
			checkForError(err)
			if(data != nil) {
				values[path] = string(*data)
			}
		case "dict":
			readDictionary(decoder, path, values)
		case "array":
			readPastElement(decoder, "array")
	}
}

func makeKey(current string, key string) (string) {
	if len(current) == 0 {
		return key
	}
	return current + "|" + key
}

func readDictionary(decoder *xml.Decoder, path string, properties map[string]string)  {
	for start, _, err := readNextStartOrEnd(decoder, "key", "dict") ; start != nil; 
			start, _, err = readNextStartOrEnd(decoder, "key", "dict"){
		checkForError(err)
		if(start != nil){
			data, _, err := readNextCharDataOrEnd(decoder)
			checkForError(err)
			if(data != nil) {
				readKeyValue(decoder, makeKey(path, string(*data)), properties)
			}
		}
	}
}

func plistToMap(data io.Reader) (map[string]string, error){
	key := ""
	var properties map[string]string
	properties = make(map[string]string)
	
	decoder := xml.NewDecoder(data)
	readNextStartUntil(decoder, "dict")
	readDictionary(decoder, key, properties)
	return properties, nil
}

func data() (io.Reader){
	return strings.NewReader(`<plist version="0.9">
	<dict>
		<key>CFBundleName</key>
		<string>Visual Paradigm CE 11.1 Installer</string>
		<key>CFBundleExecutable</key>
		<string>JavaApplicationStub</string>
		<key>CFBundleIconFile</key>
		<string>app.icns</string>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>CFBundlePackageType</key>
		<string>APPL</string>
		<key>CFBundleSignature</key>
		<string>????</string>
		<key>CFBundleIdentifier</key>
		<string>com.install4j.1106-5897-7327-6550.37147</string>
		<key>CFBundleVersion</key>
		<string>11.1</string>
		<key>CFBundleShortVersionString</key>
		<string>11.1</string>
		<key>CFBundleGetInfoString</key>
		<string>11.1</string>
		<key>CFBundleDevelopmentRegion</key>
		<string>en</string>
		<key>CFBundleAllowMixedLocalizations</key>
		<true />
		<key>NSHighResolutionCapable</key>
		<true />
		<!-- I4J_INSERT_DOCTYPE -->
		<key>Java</key>
		<dict>
			<key>MainClass</key>
			<string>com.install4j.runtime.launcher.MacLauncher</string>
			<key>VMOptions</key>
			<string></string><!-- I4J_INSERT_VMOPTIONS -->
			<key>Arguments</key>
			<string></string>
			<key>Properties</key>
			<dict>
				<key>exe4j.moduleName</key>
				<string>$APP_PACKAGE</string>
				<key>sun.java2d.noddraw</key>
				<string>true</string>
				<!-- I4J_INSERT_PROPERTIES -->
			</dict>
			<key>JVMVersion</key>
			<string>1.4+</string>
			<key>ClassPath</key>
			<string>$APP_PACKAGE/Contents/Resources/app/.install4j/i4jruntime.jar:$APP_PACKAGE/Contents/Resources/app/user.jar</string><!-- 
				I4J_INSERT_CLASSPATH -->
			<key>WorkingDirectory</key>
			<string>$APP_PACKAGE/Contents/Resources/app/.</string>
		</dict>
	</dict>
</plist>`)
}