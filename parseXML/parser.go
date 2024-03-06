package parser

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
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
type Param struct {
	Em     string `default:"somemail@mail.com"`
	Pc     string `default:"5771"`
	Ld     int    `default:"1"`
	Poi    Poi
	Email  bool `default:"true"`
	Mobile bool `defualt:"false"`
}

func Parser() (poi Poi, poa Poa) {
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
	poi = Poi{
		dob:    dateOfBirth,
		gender: sex,
		name:   name,
		email:  e,
		mobile: m,
	}

	// fmt.Printf("%+v\n", poi)

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

	poa = Poa{
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
	// fmt.Printf("%+v\n", poa)

	// Photograph
	photo := uidData.SelectElement("//Pht").InnerText()

	decodedPhoto, err := base64.StdEncoding.DecodeString(photo)

	if err != nil {
		log.Fatalf("Failed to decode photo string: %v", err)
	}
	// fmt.Println(decodedPhoto)
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
	return
}

func generateSHA(em string, passCode string, lastDigit int) string {
	// Generate hash from email and mobile number

	res := em + passCode
	if lastDigit < 2 {
		lastDigit = 1
	}
	for i := 1; i <= lastDigit; i++ {
		h := sha256.Sum256([]byte(res))
		res = hex.EncodeToString(h[:])
	}
	return res
}

func VerifyEM(prm Param) (ok bool, err error) {
	// check if aadhar data has mobile and email
	switch {
	case prm.Email:
		if prm.Poi.email == "" {
			ok = false
			err = errors.New("email: Not Provided during enrollment.")
		} else if prm.Poi.email == generateSHA(prm.Em, prm.Pc, prm.Ld) {
			ok = true
			err = nil
		} else {
			ok = false
			err = errors.New("email: not verified")
		}
	case prm.Mobile:
		if prm.Poi.mobile == "" {
			ok = false
			err = errors.New("mobile: Not Provided during enrollment.")
		} else if prm.Poi.mobile == generateSHA(prm.Em, prm.Pc, prm.Ld) {
			ok = true
			err = nil
		} else {
			ok = false
			err = errors.New("mobile: not verified")
		}
	}
	return
}
