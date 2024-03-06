package main

import (
	"fmt"
	"io"
	"os"

	parser "github.com/arunkhattri/uidai_go/parseXML"
	"github.com/moov-io/signedxml"
)

func main() {
	// poi, poa := parser.Parser()
	poi, _ := parser.Parser()

	prm1 := parser.Param{
		Em:     "arun.kr.khattri@gmail.com",
		Poi:    poi,
		Pc:     "5771",
		Ld:     1,
		Email:  true,
		Mobile: false,
	}
	prm2 := parser.Param{
		Em:     "9312829394",
		Poi:    poi,
		Pc:     "5771",
		Ld:     1,
		Email:  false,
		Mobile: true,
	}

	// fmt.Printf("[POI]\n%+v\n", poi)
	// fmt.Printf("[POA]\n%+v\n", poa)

	// email
	ok, err := parser.VerifyEM(prm1)
	if err != nil {
		fmt.Println(err)
	}
	if ok {
		fmt.Println("email is verified.")
	}
	// mobile
	ok2, err2 := parser.VerifyEM(prm2)
	if err2 != nil {
		fmt.Println(err2)
	}
	if ok2 {
		fmt.Println("mobile is verified.")
	}

	// verify xml file
	xmlFile, err := os.Open("/home/arunkhattri/github/arunkhattri/ak-yojna/uidai_go/data/arun_eKYC.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()

	bxml, _ := io.ReadAll(xmlFile)
	sxml := fmt.Sprintf("%s", bxml)

	validator, err := signedxml.NewValidator(sxml)
	_, err = validator.ValidateReferences()
	if err != nil {
		fmt.Println("Doc verified")
	}
}
