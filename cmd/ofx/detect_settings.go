package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yury-dymov/ofxgo"
)

var detectSettingsCommand = command{
	Name:        "detect-settings",
	Description: "Attempt to guess client settings needed for a particular financial institution",
	Flags:       flag.NewFlagSet("detect-settings", flag.ExitOnError),
	CheckFlags:  checkServerFlags,
	Do:          detectSettings,
}

var delay uint64

func init() {
	detectSettingsCommand.Flags.StringVar(&serverURL, "url", "", "Financial institution's OFX Server URL (see ofxhome.com if you don't know it)")
	detectSettingsCommand.Flags.StringVar(&username, "username", "", "Your username at financial institution")
	detectSettingsCommand.Flags.StringVar(&password, "password", "", "Your password at financial institution")
	detectSettingsCommand.Flags.StringVar(&org, "org", "", "'ORG' for your financial institution")
	detectSettingsCommand.Flags.StringVar(&fid, "fid", "", "'FID' for your financial institution")
	detectSettingsCommand.Flags.Uint64Var(&delay, "delay", 500, "How long to delay between two subsequent requests, in milliseconds")
}

// We keep a separate list of APPIDs to preserve the ordering (ordering isn't
// guaranteed in maps). We want to try them in order from 'best' and most
// likely to work to 'worse' and least likely to work
var appIDs = []string{
	"OFXGO", // ofxgo (this library)
	"QWIN",  // Intuit Quicken Windows
	"QMOFX", // Intuit Quicken Mac
	"QB",    // Intuit QuickBooks Windows
	"Money", // Microsoft Money 2007
}

var appVersions = map[string][]string{
	"OFXGO": { // ofxgo (this library)
		"0001",
	},
	"QWIN": { // Intuit Quicken Windows
		"2600", // 2017
		"2500", // 2016
		"2400", // 2015
		"2300", // 2014
		"2200", // 2013
		"2100", // 2012
		"2000", // 2011
		"1900", // 2010
		"1800", // 2009
		"1700", // 2008
		"1600", // 2007
		"1500", // 2006
		"1400", // 2005
	},
	"QMOFX": { // Intuit Quicken Mac
		"1700", // 2008
		"1600", // 2007
		"1500", // 2006
		"1400", // 2005
	},
	"QB": { // Intuit QuickBooks Windows
		"1800", // 2008
		"1700", // 2007
		"1600", // 2006
		"1500", // 2005
	},
	"Money": { // Microsoft Money 2007
		"1600", // 2007
		"1500", // 2006
		"1400", // 2005
		"1200", // 2004
		"1100", // 2003
	},
}

var versions = []string{
	"203",
	"103",
	"200",
	"201",
	"202",
	"210",
	"211",
	"102",
	"151",
	"160",
	"220",
}

func detectSettings() {
	var attempts uint
	for _, appID := range appIDs {
		for _, appVer := range appVersions[appID] {
			for _, version := range versions {
				for _, noIndent := range []bool{false, true} {
					if tryProfile(appID, appVer, version, noIndent) {
						fmt.Println("The following settings were found to work:")
						fmt.Printf("AppID: %s\n", appID)
						fmt.Printf("AppVer: %s\n", appVer)
						fmt.Printf("OFX Version: %s\n", version)
						fmt.Printf("noindent: %t\n", noIndent)
						os.Exit(0)
					} else {
						attempts++
						var noIndentString string
						if noIndent {
							noIndentString = " noindent"
						}
						fmt.Printf("Attempt %d failed (%s %s %s%s), trying again after %dms...\n", attempts, appID, appVer, version, noIndentString, delay)
						time.Sleep(time.Duration(delay) * time.Millisecond)
					}
				}
			}
		}
	}
}

const anonymous = "anonymous00000000000000000000000"

func tryProfile(appID, appVer, version string, noindent bool) bool {
	ver, err := ofxgo.NewOfxVersion(version)
	if err != nil {
		fmt.Println("Error creating new OfxVersion enum:", err)
		os.Exit(1)
	}
	var client = ofxgo.GetClient(serverURL,
		&ofxgo.BasicClient{
			AppID:       appID,
			AppVer:      appVer,
			SpecVersion: ver,
			NoIndent:    noindent,
		})

	var query ofxgo.Request
	query.URL = serverURL
	query.Signon.ClientUID = ofxgo.UID(clientUID)
	query.Signon.UserID = ofxgo.String(username)
	query.Signon.UserPass = ofxgo.String(password)
	query.Signon.Org = ofxgo.String(org)
	query.Signon.Fid = ofxgo.String(fid)

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	profileRequest := ofxgo.ProfileRequest{
		TrnUID:   *uid,
		DtProfUp: ofxgo.Date{Time: time.Unix(0, 0)},
	}
	query.Prof = append(query.Prof, &profileRequest)

	_, err = client.Request(&query)
	if err == nil {
		return true
	}

	// try again with anonymous logins
	query.Signon.UserID = ofxgo.String(anonymous)
	query.Signon.UserPass = ofxgo.String(anonymous)

	_, err = client.Request(&query)
	return err == nil
}
