package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"strings"

	"bitbucket.bit.admin.ch/dipdvs/validatexml_xsd/xsdvalidate"
)

// cat crs_payload.xml | ./validatexml_xsd -s CrsXML_v2.0.xsd -v
func main() {

	doValidate := flag.Bool("v", false, "validate the XML document. Use -s switch to specify a xsd file")
	xsdPath := flag.String("s", "", "the path to the xsd file")
	xmlPath := flag.String("f", "", "the path to the input (XML) file. If omitted, standard input will be used")

	flag.Parse()

	if !*doValidate {
		fmt.Println("Please, specify -v to validate")
		os.Exit(1)
	}

	var buff []byte
	var err error
	if *xmlPath != "" {
		buff, err = ioutil.ReadFile(*xmlPath)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	} else {
		// Otherwise, try loading from standard input
		buff, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error occured while reading from standard input")
			os.Exit(1)
		}
	}

	if *doValidate {

		if *xsdPath == "" {
			fmt.Println("Please, specify a xsd schema for validate")
			os.Exit(1)
		}
		result, err := ValidateXMLWithXSD(buff, *xsdPath)
		if err != nil {
			fmt.Printf("An error occurred when trying to validate xml with xsd, error: %s\n", err.Error())
			os.Exit(1)
		}
		if len(result) > 0 {
			os.Stdout.Write([]byte(strings.Join(result, ";")))
			os.Exit(0)
		}
	}
	retResult := "Xml is Valid!"
	os.Stdout.Write([]byte(retResult))
	os.Exit(0)
}

// ValidateXMLWithXSD validates an xml file against the xml schema
func ValidateXMLWithXSD(xmlPayload []byte, xsdURL string) ([]string, error) {

	xsdvalidate.Init()
	defer xsdvalidate.Cleanup()
	xsdhandler, err := xsdvalidate.NewXsdHandlerUrl(xsdURL, xsdvalidate.ParsErrDefault)
	if err != nil {
		return nil, err
	}
	defer xsdhandler.Free()
	xmlhandler, err := xsdvalidate.NewXmlHandlerMem(xmlPayload, xsdvalidate.ParsErrDefault)
	if err != nil {
		return nil, err
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, xsdvalidate.ValidErrDefault)
	if err != nil {
		switch err.(type) {
		case xsdvalidate.ValidationError:
			errorMsgs := []string{}
			for i := range err.(xsdvalidate.ValidationError).Errors {
				msg := strings.Replace(err.(xsdvalidate.ValidationError).Errors[i].Message, "'", "", -1)
				errorMsgs = append(errorMsgs, fmt.Sprintf("Error Line %d: %s", err.(xsdvalidate.ValidationError).Errors[i].Line, msg))
			}
			return errorMsgs, nil
		default:
			return nil, err
		}
	}
	return nil, nil
}
