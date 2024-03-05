package parser

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/antchfx/xmlquery"
)

type Poi struct {
	dob    string
	gender string
	name   string
	email  string
	mobile string
}

type Poa struct {
	careof  string
	country string
	dist    string
	house   string
	loc     string
	pc      string
	po      string
	state   string
	street  string
	subdist string
	vtc     string
}

func Parser() {
	// var poa Poa
	//open our xmlFile
	xmlFile, err := os.Open("/home/arunkhattri/github/arunkhattri/ak-yojna/uidai_go/data/arun_eKYC.xml")
	if err != nil {
		fmt.Println("Error in opening xml file", err)
	}
	doc, err := xmlquery.Parse(xmlFile)
	if err != nil {
		fmt.Println("Error in parsing xml file", err)
	}

	defer xmlFile.Close()

	root := xmlquery.FindOne(doc, "//OfflinePaperlessKyc")

	uidData := root.SelectElement("//UidData")

	// Poi
	poiData := uidData.SelectElement("//Poi")
	poaData := uidData.SelectElement("//Poa")
	dateOfBirth := poiData.SelectAttr("dob")
	e := poiData.SelectAttr("e")
	sex := poiData.SelectAttr("gender")
	m := poiData.SelectAttr("m")
	name := poiData.SelectAttr("name")
	poi := Poi{
		dob:    dateOfBirth,
		gender: sex,
		name:   name,
		email:  e,
		mobile: m,
	}

	fmt.Printf("%+v\n", poi)

	// Poa
	careof := poaData.SelectAttr("careof")
	country := poaData.SelectAttr("country")
	dist := poaData.SelectAttr("dist")
	house := poaData.SelectAttr("house")
	loc := poaData.SelectAttr("loc")
	pc := poaData.SelectAttr("pc")
	po := poaData.SelectAttr("po")
	state := poaData.SelectAttr("state")
	street := poaData.SelectAttr("street")
	subdist := poaData.SelectAttr("subdist")
	vtc := poaData.SelectAttr("vtc")

	poa := Poa{
		careof:  careof,
		country: country,
		dist:    dist,
		house:   house,
		loc:     loc,
		pc:      pc,
		po:      po,
		state:   state,
		street:  street,
		subdist: subdist,
		vtc:     vtc,
	}
	fmt.Printf("%+v\n", poa)

	// Photograph
	photo := uidData.SelectElement("//Pht").InnerText()

	decodedPhoto, err := base64.StdEncoding.DecodeString(photo)

	if err != nil {
		log.Fatalf("Failed to decode photo string: %v", err)
	}
	fmt.Println(decodedPhoto)
	r := bytes.NewReader(decodedPhoto)
	im, err := jpeg.Decode(r)
	if err != nil {
		panic("Bad jpg")
	}

	f, err := os.OpenFile("ak_uid_img.png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}
	jpeg.Encode(f, im, &jpeg.Options{Quality: 75})
	fmt.Println("jpeg created...")

}
