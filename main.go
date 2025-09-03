package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Function int64

const (
	Linear Function = 0
	Exp    Function = 1
	Log    Function = 2
)

type Platform struct {
	minRange            float64
	maxRange            float64
	minDispersionRadius float64
	maxDispersionRadius float64
	name                string
}

type Ammo struct {
	fullDamageRadius    float64
	partialDamageRadius float64
	name                string
}

// Platforms

var Cremari = Platform{45.0, 80.0, 5.5, 12, "Cremari Mortar"}
var GunboatMortar = Platform{75.0, 100.0, 2.5, 14.5, "Gunboat Mortar"}
var HuberLariat = Platform{100.0, 300.0, 25.0, 35.0, "Huber Lariat"}
var Koronides = Platform{100.0, 250.0, 22.5, 30.0, "Koronides"}
var HuberExalt = Platform{100.0, 300.0, 25.0, 35.0, "Huber Exalt"}
var Thunderbolt = Platform{200.0, 350.0, 32.5, 40.0, "Thunderbolt"}
var StormCannon = Platform{400.0, 1000.0, 50.0, 50.0, "Storm Cannon"}
var TempestCannon = Platform{350.0, 500.0, 50.0, 50.0, "Tempest Cannon"}
var Conqueror = Platform{100.0, 200.0, 2.5, 8.5, "Conqueror"}
var Blacksteele = Platform{100.0, 200.0, 2.5, 8.5, "Blacksteele"}
var DevittCaine = Platform{45.0, 80.0, 2.5, 9.45, "Devitt-Caine"}
var Peltast = Platform{45.0, 80.0, 2.5, 9.45, "Peltast"}
var Skycaller = Platform{275.0, 350.0, 37.5, 60.0, "Skycaller"}
var Deioneus = Platform{350.0, 400.0, 41.5, 57.5, "Deioneus"}
var WaspNest = Platform{375.0, 450.0, 37.5, 60.0, "WaspNest"}
var Squire = Platform{375.0, 500.0, 39, 51, "Squire"}
var Hades = Platform{300.0, 575.0, 35, 52, "Hades"}
var Retiarius = Platform{375.0, 500.0, 37.5, 51, "Retiarius"}
var Trident = Platform{100.0, 225.0, 2.5, 8.5, "Trident"}
var Titan120 = Platform{100.0, 200.0, 2.5, 8.5, "Titan"}
var Callahan120 = Platform{100.0, 200.0, 2.5, 8.5, "Callahan"}
var Titan150 = Platform{100.0, 225.0, 2.5, 8.5, "Titan"}
var Callahan150 = Platform{100.0, 225.0, 2.5, 8.5, "Callahan"}
var Flood = Platform{120.0, 250.0, 25, 35, "Flood"}
var Sarissa = Platform{120.0, 250.0, 25, 35, "Sarissa"}

var Platforms = map[string]Platform{
	"cremari":       Cremari,
	"gunboat":       GunboatMortar,
	"huberlariat":   HuberLariat,
	"koronides":     Koronides,
	"huberexalt":    HuberExalt,
	"stormcannon":   StormCannon,
	"tempestcannon": TempestCannon,
	"conqueror":     Conqueror,
	"blacksteele":   Blacksteele,
	"skycaller":     Skycaller,
	"devittcaine":   DevittCaine,
	"peltast":       Peltast,
	"deioneus":      Deioneus,
	"waspnest":      WaspNest,
	"squire":        Squire,
	"hades":         Hades,
	"retiarius":     Retiarius,
	"thunderbolt":   Thunderbolt,
	"callahan120":   Callahan120,
	"callahan150":   Callahan150,
	"titan120":      Titan120,
	"titan150":      Titan150,
	"trident":       Trident,
	"flood":         Flood,
	"sarissa":       Sarissa,
}

// Ammos

var HEMortarShell = Ammo{2.5, 5.0, "HE Mortar Shell"}
var ShrapnelMortarShell = Ammo{4.5, 7.5, "Shrapnel Mortar Shell"}
var IncendiaryMortarShell = Ammo{2.25, 5.5, "Incendiary Mortar Shell"}
var _120mmShell = Ammo{4.0, 11.25, "120mm Shell"}
var _150mmShell = Ammo{7.0, 11.25, "150mm Shell"}
var _300mmShell = Ammo{4.5, 15.0, "300mm Shell"}
var ExplosiveRocket = Ammo{2.25, 5.5, "High Explosive Rocket"}
var FireRocket = Ammo{2.25, 5.5, "Fire Rocket"}

var Ammos = map[string]Ammo{
	"he":         HEMortarShell,
	"shrapnel":   ShrapnelMortarShell,
	"incendiary": IncendiaryMortarShell,
	"120mm":      _120mmShell,
	"150mm":      _150mmShell,
	"300mm":      _300mmShell,
	"explosive":  ExplosiveRocket,
	"fire":       FireRocket,
}

const MinYAxis = 0
const MaxYAxis = 20000
const RangeYAxis = MaxYAxis - MinYAxis

const SVGTemplate = `<svg viewBox="0 0 32000 22000" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" stroke-width="28.222" stroke-linejoin="round" xml:space="preserve">
  <style>
    line {
      stroke: black;
    }
  </style>
 <path fill="rgb(255,255,255)" stroke="none" d="M 0,0 L 32000,0 32000,22000 0,22000 0,0"/>

 <path fill="rgb(255,134,134)" stroke="none" d="M 0,0 L 30000,0 30000,20000 0,20000 0,0 Z "/>

 <path fill="rgb(254,254,123)" stroke="none" d="M 0,%d L 7500,%d 15000,%d 22500,%d 30000,%d 30000,20000 0,20000 0,%d Z "/>

 <path fill="rgb(19,174,0)" stroke="none" d="M 0,%d L 7500,%d 15000,%d 22500,%d, 30000,%d 30000,20000 0,20000 0,%d Z "/>

 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="0" y="20549"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">%s</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="7500" y="20549"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">%s</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="15000" y="20549"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">%s</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="22500" y="20549"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">%s</tspan></tspan></tspan></text>
 <line x1="7500" y1="0" x2="7500" y2="20000" stroke-dasharray="40" />
 <line x1="15000" y1="0" x2="15000" y2="20000" stroke-dasharray="40" />
 <line x1="22500" y1="0" x2="22500" y2="20000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="30000" y="20549"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">%s</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="1004" y="20000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">0%%</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="18000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">10%%</tspan></tspan></tspan></text>
 <line x1="0" y1="18000" x2="30000" y2="18000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="16000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">20%%</tspan></tspan></tspan></text>
 <line x1="0" y1="16000" x2="30000" y2="16000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="14000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">30%%</tspan></tspan></tspan></text>
 <line x1="0" y1="14000" x2="30000" y2="14000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="12000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">40%%</tspan></tspan></tspan></text>
 <line x1="0" y1="12000" x2="30000" y2="12000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="10000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">50%%</tspan></tspan></tspan></text>
 <line x1="0" y1="10000" x2="30000" y2="10000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="8000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">60%%</tspan></tspan></tspan></text>
 <line x1="0" y1="8000" x2="30000" y2="8000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="6000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">70%%</tspan></tspan></tspan></text>
 <line x1="0" y1="6000" x2="30000" y2="6000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="4000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">80%%</tspan></tspan></tspan></text>
 <line x1="0" y1="4000" x2="30000" y2="4000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="813" y="2000"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">90%%</tspan></tspan></tspan></text>
 <line x1="0" y1="2000" x2="30000" y2="2000" stroke-dasharray="40" />
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="621" y="0"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">100%%</tspan></tspan></tspan></text>
 <path fill="rgb(19,174,0)" stroke="none" d="M 30642,11121 L 30536,11121 30536,10910 30747,10910 30747,11121 30642,11121 Z "/>
 <path fill="rgb(254,254,123)" stroke="none" d="M 30642,10623 L 30536,10623 30536,10413 30747,10413 30747,10623 30642,10623 Z "/>
 <path fill="rgb(255,134,134)" stroke="none" d="M 30642,10126 L 30536,10126 30536,9915 30747,9915 30747,10126 30642,10126 Z "/>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="30847" y="10140"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">Miss</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="30847" y="10638"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">Partial</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="30847" y="11135"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">Hit</tspan></tspan></tspan></text>
 <text class="SVGTextShape"><tspan class="TextParagraph"><tspan class="TextPosition" x="13000" y="21500"><tspan font-family="Liberation Sans, sans-serif" font-size="353px" font-weight="400" fill="rgb(0,0,0)" stroke="none" style="white-space: pre">%s</tspan></tspan></tspan></text>
 </svg>`

const HTMLTemplate = `
<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width,initial-scale=1" />
<style>
    .container {
	display: grid;
	grid-template-columns: 1fr;
	gap: 20px;
	font: 1em "Liberation Sans", sans-serif;
    }
    .form-container {
	display: grid;
	justify-content: center;
	grid-template-columns: 1fr 1fr 1fr 1fr;
	font: 1.2em "Liberation Sans", sans-serif;
    }
    .form-center {
	justify-content:center;
	display: grid;
    }
    .svg-container {
	display: grid;
	position:relative;
	left:25%%;
    }
    button, input, select, textarea {
	font: inherit;
    }
    svg {
		max-height: 50%%;
    }
</style>
</head>
<body class="container">
<div class="form-center">
<form class="form-container" action="" method="post">
  <div class="form-center">
    <label for="name">Target Radius (meters):</label>
    <input type="number" name="targetRadius" id="name" value=%.1f step=0.1 required />
  </div>
  <div class="form-center">
  <label for="weapon-select">Platform and Ammo:</label>
  <select name="weapon" id="weapon-select">
  <option value="">--- Mortar ---</option>
  <option value="cremari-he-1">Cremari with HE shell</option>
  <option value="cremari-shrapnel-2">Cremari with Shrapnel shell</option>
  <option value="cremari-incendiary-3">Cremari with Incendiary shell</option>
  <option value="gunboat-he-4">Gunboat with HE shell</option>
  <option value="gunboat-shrapnel-5">Gunboat with Shrapnel shell</option>
  <option value="gunboat-incendiary-6">Gunboat with Incendiary shell</option>
  <option value="devittcaine-he-7">Devitt-Caine with HE shell</option>
  <option value="devittcaine-shrapnel-8">Devitt-Caine with Shrapnel shell</option>
  <option value="devittcaine-incendiary-9">Devitt-Caine with Incendiary shell</option>
  <option value="peltast-he-10">Peltast with HE shell</option>
  <option value="peltast-shrapnel-11">Peltast with Shrapnel shell</option>
  <option value="peltast-incendiary-12">Peltast with Incendiary shell</option>

  <option value="">--- 120mm ---</option>

  <option value="trident-120mm-14">Trident with 120mm shell</option>
  <option value="koronides-120mm-15">Koronides with 120mm shell</option>
  <option value="huberlariat-120mm-16">Huber Lariat with 120mm shell</option>
  <option value="conqueror-120mm-17">Conqueror with 120mm shell</option>
  <option value="titan120-120mm-18">Titan with 120mm shell</option>
  <option value="blacksteele-120mm-19">Blacksteele with 120mm shell</option>
  <option value="callahan120-120mm-20">Callahan with 120mm shell</option>

  <option value="">--- 150mm ---</option>

  <option value="huberexalt-150mm-22">Huber Exalt with 150mm shell</option>
  <option value="flood-150mm-23">Flood with 150mm shell</option>
  <option value="sarissa-150mm-24">Sarissa with 150mm shell</option>
  <option value="thunderbolt-150mm-25">Thunderbolt with 150mm shell</option>
  <option value="titan150-150mm-26">Titan with 150mm shell</option>
  <option value="callahan150-150mm-27">Callahan with 150mm shell</option>

  <option value="">--- 300mm ---</option>

  <option value="stormcannon-300mm-29">Storm Cannon with 300mm shell</option>
  <option value="tempestcannon-300mm-30">Tempest Cannon with 300mm shell</option>

  <option value="">--- HE Rocket ---</option>

  <option value="squire-explosive-32">Squire with Explosive rocket</option>
  <option value="hades-explosive-33">Hades with Explosive rocket</option>
  <option value="retiarius-explosive-34">Retiarius with Explosive rocket</option>

  <option value="">--- Fire Rocket ---</option>

  <option value="skycaller-fire-36">Skycaller with Fire rocket</option>
  <option value="deioneus-fire-37">Deioneus with Fire rocket</option>
  <option value="waspnest-fire-38">Wasp Nest with Fire rocket</option>
  </select>
  </div>
  <div class="form-center">
  <label for="function-select">Function for Dispersion Radius</label>
  <select name="function" id="function-select">
  <option value="0">Linear</option>
  <option value="1">Exponential</option>
  <option value="2">Logarithmic</option>
	<select/>
  </div>
  <div class="form-center">
    <input type="submit" value="Plot chart" />
  </div>
</form>
</div>
<div class=svg-container>
%s
</div>
<script type="text/javascript" >
document.getElementById("weapon-select").selectedIndex = %s
document.getElementById("function-select").selectedIndex = %s
</script>
</body>
</html>
`

func computeDispersionRadius(minDispersionRadius, maxDispersionRadius, minRange, maxRange float64, function Function) (float64, float64, float64) {
	var aCoeff float64
	var bCoeff float64
	var _25percentDispersionRadius float64
	var midDispersionRadius float64
	var _75percentDispersionRadius float64
	midRange := minRange + (maxRange-minRange)*0.5
	_25percentRange := minRange + (maxRange-minRange)*0.25
	_75percentRange := minRange + (maxRange-minRange)*0.75
	switch function {
	case Linear: //  linear function => y = a.x + b
		aCoeff = (maxDispersionRadius - minDispersionRadius) / (maxRange - minRange)
		bCoeff = minDispersionRadius - (aCoeff * minRange)
		_25percentDispersionRadius = aCoeff*_25percentRange + bCoeff
		midDispersionRadius = aCoeff*midRange + bCoeff
		_75percentDispersionRadius = aCoeff*_75percentRange + bCoeff
	case Exp: // exp function => y = a.exp(x.b)
		bCoeff = 1 / (minRange - maxRange) * math.Log(minDispersionRadius/maxDispersionRadius)
		aCoeff = maxDispersionRadius / math.Exp(bCoeff*maxRange)
		_25percentDispersionRadius = aCoeff * math.Exp(_25percentRange*bCoeff)
		midDispersionRadius = aCoeff * math.Exp(midRange*bCoeff)
		_75percentDispersionRadius = aCoeff * math.Exp(_75percentRange*bCoeff)
	case Log: // log function => y = a.ln(x.b)
		aCoeff = (minDispersionRadius - maxDispersionRadius) / (math.Log(minRange) - math.Log(maxRange))
		bCoeff = math.Exp(minDispersionRadius/aCoeff) / minRange
		_25percentDispersionRadius = aCoeff * math.Log(_25percentRange*bCoeff)
		midDispersionRadius = aCoeff * math.Log(midRange*bCoeff)
		_75percentDispersionRadius = aCoeff * math.Log(_75percentRange*bCoeff)

	}

	return _25percentDispersionRadius, midDispersionRadius, _75percentDispersionRadius
}

func computeHits(targetRadius float64, platform Platform, ammo Ammo, function Function) string {
	title := fmt.Sprintf("%s with %s on %.1fm radius target", platform.name, ammo.name, targetRadius)
	fmt.Println(title)

	minRange := platform.minRange
	maxRange := platform.maxRange
	midRange := minRange + (maxRange-minRange)*0.5
	_25percentRange := minRange + (maxRange-minRange)*0.25
	_75percentRange := minRange + (maxRange-minRange)*0.75

	maxDispersionRadius := platform.maxDispersionRadius
	minDispersionRadius := platform.minDispersionRadius

	_25percentDispersionRadius, midDispersionRadius, _75percentDispersionRadius := computeDispersionRadius(minDispersionRadius, maxDispersionRadius, minRange, maxRange, function)

	fmt.Println("Min Dispertion Radius", minDispersionRadius)
	fmt.Println("25pc Dispertion Radius", _25percentDispersionRadius)
	fmt.Println("Mid Dispertion Radius", midDispersionRadius)
	fmt.Println("75pc Dispertion Radius", _75percentDispersionRadius)
	fmt.Println("Max Dispertion Radius", maxDispersionRadius)

	fullDamageRadius := ammo.fullDamageRadius
	partialDamageRadius := ammo.partialDamageRadius

	maxDispersionArea := math.Pi * maxDispersionRadius * maxDispersionRadius
	_25percentDispersionArea := math.Pi * _25percentDispersionRadius * _25percentDispersionRadius
	midDispersionArea := math.Pi * midDispersionRadius * midDispersionRadius
	_75percentDispersionArea := math.Pi * _75percentDispersionRadius * _75percentDispersionRadius
	minDispersionArea := math.Pi * minDispersionRadius * minDispersionRadius

	fullArea := math.Pi * (targetRadius + fullDamageRadius) * (targetRadius + fullDamageRadius)
	partialArea := (math.Pi * (targetRadius + partialDamageRadius) * (targetRadius + partialDamageRadius))

	yMinPartial := MaxYAxis - RangeYAxis*partialArea/minDispersionArea
	yMinFull := MaxYAxis - RangeYAxis*fullArea/minDispersionArea

	y25pcPartial := MaxYAxis - RangeYAxis*partialArea/_25percentDispersionArea
	y25pcFull := MaxYAxis - RangeYAxis*fullArea/_25percentDispersionArea

	y75pcPartial := MaxYAxis - RangeYAxis*partialArea/_75percentDispersionArea
	y75pcFull := MaxYAxis - RangeYAxis*fullArea/_75percentDispersionArea

	yMidPartial := MaxYAxis - RangeYAxis*partialArea/midDispersionArea
	yMidFull := MaxYAxis - RangeYAxis*fullArea/midDispersionArea

	yMaxPartial := MaxYAxis - RangeYAxis*partialArea/maxDispersionArea
	yMaxFull := MaxYAxis - RangeYAxis*fullArea/maxDispersionArea

	fmt.Println("Hit % at minRange:", math.Round(fullArea/minDispersionArea*100))
	fmt.Println("Hit % at midRange:", math.Round(fullArea/midDispersionArea*100))
	fmt.Println("Hit % at maxRange:", math.Round(fullArea/maxDispersionArea*100))

	fmt.Println("Partial Hit % at minRange:", math.Round(partialArea/minDispersionArea*100))
	fmt.Println("Partial Hit % at midRange:", math.Round(partialArea/midDispersionArea*100))
	fmt.Println("Partial Hit % at maxRange:", math.Round(partialArea/maxDispersionArea*100))

	minRangeS := fmt.Sprintf("%.fm", minRange)
	_25percentRangeS := fmt.Sprintf("%.fm", _25percentRange)
	midRangeS := fmt.Sprintf("%.fm", midRange)
	_75percentRangeS := fmt.Sprintf("%.fm", _75percentRange)
	maxRangeS := fmt.Sprintf("%.fm", maxRange)

	out := fmt.Sprintf(SVGTemplate, int(yMinPartial), int(y25pcPartial), int(yMidPartial), int(y75pcPartial), int(yMaxPartial), int(yMinPartial), int(yMinFull), int(y25pcFull), int(yMidFull), int(y75pcFull), int(yMaxFull), int(yMinFull), minRangeS, _25percentRangeS, midRangeS, _75percentRangeS, maxRangeS, title)
	fileHandler, _ := os.Create("out.svg")
	fileHandler.Write([]byte(out))
	fileHandler.Close()
	return out
}

func getArty(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Printf("got /arty request\n")
		io.WriteString(w, fmt.Sprintf(HTMLTemplate, 0.0, "", "0"))
	} else if r.Method == "POST" {
		fmt.Printf("got POST /arty request\n")
		targetRadius := r.PostFormValue("targetRadius")
		radius, _ := strconv.ParseFloat(targetRadius, 64)
		function := r.PostFormValue("function")
		functionIndex, _ := strconv.ParseInt(function, 10, 64)

		weapon := r.PostFormValue("weapon")
		weapons := strings.Split(weapon, "-")
		var platform Platform
		var ammo Ammo
		var index string
		if len(weapons) == 1 {
			platform = Cremari
			ammo = HEMortarShell
			index = "1"
		} else {
			platform = Platforms[weapons[0]]
			ammo = Ammos[weapons[1]]
			index = weapons[2]
		}
		graph := computeHits(radius, platform, ammo, Function(functionIndex))
		fmt.Println(weapon)
		io.WriteString(w, fmt.Sprintf(HTMLTemplate, radius, graph, index, function))
	}
}

func main() {

	http.HandleFunc("/foxholeArty", getArty)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
