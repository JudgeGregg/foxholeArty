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

type Weapon struct {
	minRange            float64
	maxRange            float64
	minDispersionRadius float64
	maxDispersionRadius float64
	name                string
}

type Shell struct {
	fullDamageRadius    float64
	partialDamageRadius float64
	name                string
}

// Weapons

var Cremari = Weapon{45.0, 80.0, 5.5, 12, "Cremari Mortar"}
var GunboatMortar = Weapon{75.0, 100.0, 2.5, 14.5, "Gunboat Mortar"}
var HuberLariat = Weapon{100.0, 300.0, 25.0, 35.0, "Huber Lariat"}
var Koronides = Weapon{100.0, 250.0, 22.5, 30.0, "Koronides"}
var HuberExalt = Weapon{100.0, 300.0, 25.0, 35.0, "Huber Exalt"}
var StormCannon = Weapon{400.0, 1000.0, 50.0, 50.0, "Storm Cannon"}
var TempestCannon = Weapon{350.0, 500.0, 50.0, 50.0, "Tempest Cannon"}
var Conqueror = Weapon{100.0, 200.0, 2.5, 8.5, "Conqueror DD"}
var Blacksteele = Weapon{100.0, 200.0, 2.5, 8.5, "Blacksteele FF"}
var Squire = Weapon{375.0, 500.0, 39, 51, "Squire RAC"}
var Skycaller = Weapon{275.0, 350.0, 37.5, 60, "Skycaller"}

var Weapons = map[string]Weapon{
	"cremari":       Cremari,
	"gunboat":       GunboatMortar,
	"huberlariat":   HuberLariat,
	"koronides":     Koronides,
	"huberexalt":    HuberExalt,
	"stormcannon":   StormCannon,
	"tempestcannon": TempestCannon,
	"conqueror":     Conqueror,
	"blacksteele":   Blacksteele,
	"squire":        Squire,
	"skycaller":     Skycaller,
}

// Shells

var HEMortarShell = Shell{2.5, 5.0, "HE Mortar Shell"}
var ShrapnelMortarShell = Shell{4.5, 7.5, "Shrapnel Mortar Shell"}
var IncendiaryMortarShell = Shell{2.25, 5.5, "Incendiary Mortar Shell"}
var _120mmShell = Shell{4.0, 11.25, "120mm Shell"}
var _150mmShell = Shell{7.0, 11.25, "150mm Shell"}
var _300mmShell = Shell{4.5, 15.0, "300mm Shell"}
var ExplosiveRocket = Shell{2.25, 5.5, "High Explosive Rocket"}
var FireRocket = Shell{2.25, 5.5, "High Explosive Rocket"}

var Shells = map[string]Shell{
	"he":         HEMortarShell,
	"shrapnel":   ShrapnelMortarShell,
	"incendiary": IncendiaryMortarShell,
	"120mm":      _120mmShell,
	"150mm":      _150mmShell,
	"300mm":      _300mmShell,
	"explosive":  ExplosiveRocket,
	"fire":       FireRocket,
}

const MIN_X_AXIS = 0
const MIN_Y_AXIS = 0
const MAX_X_AXIS = 30000
const MAX_Y_AXIS = 20000
const RANGE_Y_AXIS = MAX_Y_AXIS - MIN_Y_AXIS

const SVG_TEMPLATE = `<?xml version="1.0" encoding="UTF-8"?>
 <!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
 <svg width="350mm" height="220mm" viewBox="0 0 35000 22000" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" stroke-width="28.222" stroke-linejoin="round" xml:space="preserve">
          <style>
    line {
      stroke: black;
    }
  </style>
 <path fill="rgb(255,255,255)" stroke="none" d="M 0,0 L 35000,0 35000,25000 0,25000 0,0"/>

 <path fill="rgb(255,134,134)" stroke="none" d="M 0,0 L 30000,0 30000,20000 0,20000 0,0 Z "/>

 <path fill="rgb(254,254,123)" stroke="none" d="M 0,%d L 7500,%d 15000,%d 22500,%d 30000,%d 30000,20000 0,20000 0,%d Z "/>

 <path fill="rgb(19,174,0)" stroke="none" d="M 0,%d L 7500,%d 15000,%d 22500,%d, 30000,%d 30000,20000 0,20000 0,%d Z "/>


 <!-- <path fill="rgb(255,134,134)" stroke="none" d="M 1751,606 L 29094,606 29094,13393 1751,1186 1751,606 Z "/>
 <path fill="rgb(254,254,123)" stroke="none" d="M 1751,1186 L 29094,13393 29094,17268 1751,12207 1751,1186 Z "/>
 <path fill="rgb(19,174,0)" stroke="none" d="M 1751,12207 L 29094,17268 29094,19981 1751,19981 1751,12207 Z "/> -->

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

const getTemplate = `
<html>
<form action="" method="post">
  <div>
    <label for="name">Target Radius (meters)</label>
    <input type="number" name="targetRadius" id="name" value=%.1f step=0.1 required />
  </div>
  <div>
  <label for="weapon-select">Choose a Gun and a Shell:</label>
<select name="weapon" id="weapon-select">
  <option value="cremari-he-0">Cremari with HE shell</option>
  <option value="cremari-shrapnel-1">Cremari with Shrapnel shell</option>
  <option value="cremari-incendiary-2">Cremari with Incendiary shell</option>
  <option value="gunboat-he-3">Gunboat with HE shell</option>
  <option value="gunboat-shrapnel-4">Gunboat with Shrapnel shell</option>
  <option value="gunboat-incendiary-5">Gunboat with Incendiary shell</option>
  <option value="huberlariat-120mm-6">Huber Lariat with 120mm shell</option>
  <option value="koronides-120mm-7">Koronides with 120mm shell</option>
  <option value="huberexalt-150mm-8">Huber Exalt with 150mm shell</option>
  <option value="conqueror-120mm-9">Conqueror with 120mm shell</option>
  <option value="blacksteele-120mm-10">Blacksteele with 120mm shell</option>
  <option value="stormcannon-300mm-11">Storm Cannon with 300mm shell</option>
  <option value="tempestcannon-300mm-12">Tempest Cannon with 300mm shell</option>
  <option value="squire-explosive-13">Squire RAC with Explosive rocket</option>
  <option value="skycaller-fire-14">Skycaller with Fire rocket</option>
</select>
  </div>
  <div>
    <input type="submit" value="Generate chart" />
  </div>
</form>
%s
<script type="text/javascript" >
document.getElementById("weapon-select").selectedIndex = %s
</script>
</html>
`

func computeHits(targetRadius float64, weapon Weapon, shell Shell) string {
	title := fmt.Sprintf("%s with %s on %.1fm target radius", weapon.name, shell.name, targetRadius)
	fmt.Println(title)

	maxDispersionRadius := weapon.maxDispersionRadius
	minDispersionRadius := weapon.minDispersionRadius
	midDispersionRadius := minDispersionRadius + (maxDispersionRadius-minDispersionRadius)*0.5
	_75percentDispersionRadius := minDispersionRadius + (maxDispersionRadius-minDispersionRadius)*0.75
	_25percentDispersionRadius := minDispersionRadius + (maxDispersionRadius-minDispersionRadius)*0.25
	fmt.Println("Min Dispertion Radius", minDispersionRadius)
	fmt.Println("Mid Dispertion Radius", midDispersionRadius)
	fmt.Println("Max Dispertion Radius", maxDispersionRadius)
	minRange := weapon.minRange
	maxRange := weapon.maxRange
	midRange := minRange + (maxRange-minRange)*0.5
	_25percentRange := minRange + (maxRange-minRange)*0.25
	_75percentRange := minRange + (maxRange-minRange)*0.75

	fullDamageRadius := shell.fullDamageRadius
	partialDamageRadius := shell.partialDamageRadius

	maxDispersionArea := math.Pi * maxDispersionRadius * maxDispersionRadius
	_25percentDispersionArea := math.Pi * _25percentDispersionRadius * _25percentDispersionRadius
	midDispersionArea := math.Pi * midDispersionRadius * midDispersionRadius
	_75percentDispersionArea := math.Pi * _75percentDispersionRadius * _75percentDispersionRadius
	minDispersionArea := math.Pi * minDispersionRadius * minDispersionRadius

	fullArea := math.Pi * (targetRadius + fullDamageRadius) * (targetRadius + fullDamageRadius)
	partialArea := (math.Pi * (targetRadius + partialDamageRadius) * (targetRadius + partialDamageRadius))

	y_min_partial := MAX_Y_AXIS - RANGE_Y_AXIS*partialArea/minDispersionArea
	y_min_full := MAX_Y_AXIS - RANGE_Y_AXIS*fullArea/minDispersionArea

	y_25pc_partial := MAX_Y_AXIS - RANGE_Y_AXIS*partialArea/_25percentDispersionArea
	y_25pc_full := MAX_Y_AXIS - RANGE_Y_AXIS*fullArea/_25percentDispersionArea

	y_75pc_partial := MAX_Y_AXIS - RANGE_Y_AXIS*partialArea/_75percentDispersionArea
	y_75pc_full := MAX_Y_AXIS - RANGE_Y_AXIS*fullArea/_75percentDispersionArea

	y_mid_partial := MAX_Y_AXIS - RANGE_Y_AXIS*partialArea/midDispersionArea
	y_mid_full := MAX_Y_AXIS - RANGE_Y_AXIS*fullArea/midDispersionArea

	y_max_partial := MAX_Y_AXIS - RANGE_Y_AXIS*partialArea/maxDispersionArea
	y_max_full := MAX_Y_AXIS - RANGE_Y_AXIS*fullArea/maxDispersionArea

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

	out := fmt.Sprintf(SVG_TEMPLATE, int(y_min_partial), int(y_25pc_partial), int(y_mid_partial), int(y_75pc_partial), int(y_max_partial), int(y_min_partial), int(y_min_full), int(y_25pc_full), int(y_mid_full), int(y_75pc_full), int(y_max_full), int(y_min_full), minRangeS, _25percentRangeS, midRangeS, _75percentRangeS, maxRangeS, title)
	fileHandler, _ := os.Create("out.svg")
	fileHandler.Write([]byte(out))
	fileHandler.Close()
	return out
}

func getArty(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Printf("got /arty request\n")
		io.WriteString(w, fmt.Sprintf(getTemplate, 0.0, "", "0"))
	} else if r.Method == "POST" {
		fmt.Printf("got POST /arty request\n")
		weapon := r.PostFormValue("weapon")
		weapons := strings.Split(weapon, "-")
		targetRadius := r.PostFormValue("targetRadius")
		radius, _ := strconv.ParseFloat(targetRadius, 64)
		graph := computeHits(radius, Weapons[weapons[0]], Shells[weapons[1]])
		fmt.Println(weapon)
		io.WriteString(w, fmt.Sprintf(getTemplate, radius, graph, weapons[2]))
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
	return
}
