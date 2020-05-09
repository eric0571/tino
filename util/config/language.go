package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type LangResources struct {
	XMLName   xml.Name   `xml:"resource"`
	Version   string     `xml:"version,attr"`
	Propertys []Property `xml:"property"`
	//	Description string     `xml:",innerxml"`
}

type Property struct {
	Key  string `xml:"key,attr"`
	Code string `xml:"code"`
	Zh   string `xml:"zh"`
	En   string `xml:"en"`
}

var (
	LangRes = LangResources{}
)

func LoadResource() {
	file, err := os.Open("language.xml") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	err = xml.Unmarshal(data, &LangRes)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	if false {
		v := LangRes
		for i := 0; i < len(v.Propertys); i++ {
			fmt.Println("Key:[", v.Propertys[i].Key, "]\tEn:", v.Propertys[i].En, "\t\tCn-zh:", v.Propertys[i].Zh)
		}
	}
}

func GetKeyLang(key string) *Property {
	for i := 0; i < len(LangRes.Propertys); i++ {
		if key == LangRes.Propertys[i].Key {
			return &(LangRes.Propertys[i])
		}
	}
	return nil
}
