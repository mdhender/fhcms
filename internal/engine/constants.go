/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
******************************************************************************/

package engine

const TRUE = 1
const FALSE = 0

/* Minimum and maximum values for a galaxy. */
const MIN_SPECIES = 1
const MAX_SPECIES = 100
const MIN_STARS = 12
const MAX_STARS = 1000
const MIN_RADIUS = 6
const MAX_RADIUS = 50
const MAX_DIAMETER = 2 * MAX_RADIUS
const MAX_PLANETS = 9 * MAX_STARS

const HP_AVAILABLE_POP = 1500

/* Assume at least 32 bits per long word. */
const NUM_CONTACT_WORDS = MAX_SPECIES + 1

/* Star types. */
const DWARF = 1
const DEGENERATE = 2
const MAIN_SEQUENCE = 3
const GIANT = 4

/* Star Colors. */
const BLUE = 1
const BLUE_WHITE = 2
const WHITE = 3
const YELLOW_WHITE = 4
const YELLOW = 5
const ORANGE = 6
const RED = 7

/* Gases in planetary atmospheres. */
const H2 = 1   /* Hydrogen */
const CH4 = 2  /* Methane */
const HE = 3   /* Helium */
const NH3 = 4  /* Ammonia */
const N2 = 5   /* Nitrogen */
const CO2 = 6  /* Carbon Dioxide */
const O2 = 7   /* Oxygen */
const HCL = 8  /* Hydrogen Chloride */
const CL2 = 9  /* Chlorine */
const F2 = 10  /* Fluorine */
const H2O = 11 /* Steam */
const SO2 = 12 /* Sulfur Dioxide */
const H2S = 13 /* Hydrogen Sulfide */

/* Tech level ids. */
const MI = 0 /* Mining tech level. */
const MA = 1 /* Manufacturing tech level. */
const ML = 2 /* Military tech level. */
const GV = 3 /* Gravitics tech level. */
const LS = 4 /* Life Support tech level. */
const BI = 5 /* Biology tech level. */

/* Item IDs. */
const RM = 0         /* Raw Material Units. */
const PD = 1         /* Planetary Defense Units. */
const SU = 2         /* Starbase Units. */
const DR = 3         /* Damage Repair Units. */
const CU = 4         /* Colonist Units. */
const IU = 5         /* Colonial Mining Units. */
const AU = 6         /* Colonial Manufacturing Units. */
const FS = 7         /* Fail-Safe Jump Units. */
const JP = 8         /* Jump Portal Units. */
const FM = 9         /* Forced Misjump Units. */
const FJ = 10        /* Forced Jump Units. */
const GT = 11        /* Gravitic Telescope Units. */
const FD = 12        /* Field Distortion Units. */
const TP = 13        /* Terraforming Plants. */
const GW = 14        /* Germ Warfare Bombs. */
const SG1 = 15       /* Mark-1 Auxiliary Shield Generators. */
const SG2 = 16       /* Mark-2. */
const SG3 = 17       /* Mark-3. */
const SG4 = 18       /* Mark-4. */
const SG5 = 19       /* Mark-5. */
const SG6 = 20       /* Mark-6. */
const SG7 = 21       /* Mark-7. */
const SG8 = 22       /* Mark-8. */
const SG9 = 23       /* Mark-9. */
const GU1 = 24       /* Mark-1 Auxiliary Gun Units. */
const GU2 = 25       /* Mark-2. */
const GU3 = 26       /* Mark-3. */
const GU4 = 27       /* Mark-4. */
const GU5 = 28       /* Mark-5. */
const GU6 = 29       /* Mark-6. */
const GU7 = 30       /* Mark-7. */
const GU8 = 31       /* Mark-8. */
const GU9 = 32       /* Mark-9. */
const X1 = 33        /* Unassigned. */
const X2 = 34        /* Unassigned. */
const X3 = 35        /* Unassigned. */
const X4 = 36        /* Unassigned. */
const X5 = 37        /* Unassigned. */
const MAX_ITEMS = 38 /* Always bump this up to a multiple of two. Don't forget to make room for zeroth element! */

/* Status codes for named planets. These are logically ORed together. */
const HOME_PLANET = 1
const COLONY = 2
const POPULATED = 8
const MINING_COLONY = 16
const RESORT_COLONY = 32
const DISBANDED_COLONY = 64

/* Ship classes. */
const PB = 0  /* Picketboat. */
const CT = 1  /* Corvette. */
const ES = 2  /* Escort. */
const DD = 3  /* Destroyer. */
const FG = 4  /* Frigate. */
const CL = 5  /* Light Cruiser. */
const CS = 6  /* Strike Cruiser. */
const CA = 7  /* Heavy Cruiser. */
const CC = 8  /* Command Cruiser. */
const BC = 9  /* Battlecruiser. */
const BS = 10 /* Battleship. */
const DN = 11 /* Dreadnought. */
const SD = 12 /* Super Dreadnought. */
const BM = 13 /* Battlemoon. */
const BW = 14 /* Battleworld. */
const BR = 15 /* Battlestar. */
const BA = 16 /* Starbase. */
const TR = 17 /* Transport. */

const NUM_SHIP_CLASSES = 18

/* Ship types. */
const FTL = 0
const SUB_LIGHT = 1
const STARBASE = 2

/* Ship status codes. */
const UNDER_CONSTRUCTION = 0
const ON_SURFACE = 1
const IN_ORBIT = 2
const IN_DEEP_SPACE = 3
const JUMPED_IN_COMBAT = 4
const FORCED_JUMP = 5

/* Interspecies transactions. */

const MAX_TRANSACTIONS = 1000

const EU_TRANSFER = 1
const MESSAGE_TO_SPECIES = 2
const BESIEGE_PLANET = 3
const SIEGE_EU_TRANSFER = 4
const TECH_TRANSFER = 5
const DETECTION_DURING_SIEGE = 6
const SHIP_MISHAP = 7
const ASSIMILATION = 8
const INTERSPECIES_CONSTRUCTION = 9
const TELESCOPE_DETECTION = 10
const ALIEN_JUMP_PORTAL_USAGE = 11
const KNOWLEDGE_TRANSFER = 12
const LANDING_REQUEST = 13
const LOOTING_EU_TRANSFER = 14
const ALLIES_ORDER = 15

/* Command codes. */
const UNDEFINED = 0
const ALLY = 1
const AMBUSH = 2
const ATTACK = 3
const AUTO = 4
const BASE = 5
const BATTLE = 6
const BUILD = 7
const CONTINUE = 8
const DEEP = 9
const DESTROY = 10
const DEVELOP = 11
const DISBAND = 12
const END = 13
const ENEMY = 14
const ENGAGE = 15
const ESTIMATE = 16
const HAVEN = 17
const HIDE = 18
const HIJACK = 19
const IBUILD = 20
const ICONTINUE = 21
const INSTALL = 22
const INTERCEPT = 23
const JUMP = 24
const LAND = 25
const MESSAGE = 26
const MOVE = 27
const NAME = 28
const NEUTRAL = 29
const ORBIT = 30
const PJUMP = 31
const PRODUCTION = 32
const RECYCLE = 33
const RENAME = 34
const REPAIR = 35
const RESEARCH = 36
const SCAN = 37
const SEND = 38
const SHIPYARD = 39
const START = 40
const SUMMARY = 41
const SURRENDER = 42
const TARGET = 43
const TEACH = 44
const TECH = 45
const TELESCOPE = 46
const TERRAFORM = 47
const TRANSFER = 48
const UNLOAD = 49
const UPGRADE = 50
const VISITED = 51
const WITHDRAW = 52
const WORMHOLE = 53
const ZZZ = 54
const NUM_COMMANDS = ZZZ + 1

/* Constants needed for parsing. */
const UNKNOWN = 0
const TECH_ID = 1
const ITEM_CLASS = 2
const SHIP_CLASS = 3
const PLANET_ID = 4
const SPECIES_ID = 5

var type_char = []string{" ", "d", "D", " ", "g"}
var color_char = []string{" ", "O", "B", "A", "F", "G", "K", "M"}
var size_char = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var gas_string = [14]string{"   ", "H2", "CH4", "He", "NH3", "N2", "CO2", "O2", "HCl", "Cl2", "F2", "H2O", "SO2", "H2S"}

var tech_abbr = [6]string{"MI", "MA", "ML", "GV", "LS", "BI"}

var tech_name = [6]string{"Mining", "Manufacturing", "Military", "Gravitics", "Life Support", "Biology"}

var item_name = [MAX_ITEMS]string{
	"Raw Material Unit",
	"Planetary Defense Unit",
	"Starbase Unit",
	"Damage Repair Unit",
	"Colonist Unit",
	"Colonial Mining Unit",
	"Colonial Manufacturing Unit",
	"Fail-Safe Jump Unit",
	"Jump Portal Unit",
	"Forced Misjump Unit",
	"Forced Jump Unit",
	"Gravitic Telescope Unit",
	"Field Distortion Unit",
	"Terraforming Plant",
	"Germ Warfare Bomb",
	"Mark-1 Shield Generator",
	"Mark-2 Shield Generator",
	"Mark-3 Shield Generator",
	"Mark-4 Shield Generator",
	"Mark-5 Shield Generator",
	"Mark-6 Shield Generator",
	"Mark-7 Shield Generator",
	"Mark-8 Shield Generator",
	"Mark-9 Shield Generator",
	"Mark-1 Gun Unit",
	"Mark-2 Gun Unit",
	"Mark-3 Gun Unit",
	"Mark-4 Gun Unit",
	"Mark-5 Gun Unit",
	"Mark-6 Gun Unit",
	"Mark-7 Gun Unit",
	"Mark-8 Gun Unit",
	"Mark-9 Gun Unit",
	"X1 Unit",
	"X2 Unit",
	"X3 Unit",
	"X4 Unit",
	"X5 Unit"}

var item_abbr = [MAX_ITEMS]string{
	"RM", "PD", "SU", "DR", "CU", "IU", "AU", "FS",
	"JP", "FM", "FJ", "GT", "FD", "TP", "GW", "SG1",
	"SG2", "SG3", "SG4", "SG5", "SG6", "SG7", "SG8", "SG9",
	"GU1", "GU2", "GU3", "GU4", "GU5", "GU6", "GU7", "GU8",
	"GU9", "X1", "X2", "X3", "X4", "X5"}

var item_cost = [MAX_ITEMS]int{
	1, 1, 110, 50, 1, 1, 1, 25,
	100, 100, 125, 500, 50, 50000, 1000, 250,
	500, 750, 1000, 1250, 1500, 1750, 2000, 2250,
	250, 500, 750, 1000, 1250, 1500, 1750, 2000,
	2250, 9999, 9999, 9999, 9999, 9999}

var item_carry_capacity = [MAX_ITEMS]int{
	1, 3, 20, 1, 1, 1, 1, 1,
	10, 5, 5, 20, 1, 100, 100, 5,
	10, 15, 20, 25, 30, 35, 40, 45,
	5, 10, 15, 20, 25, 30, 35, 40,
	45, 9999, 9999, 9999, 9999, 9999}

var item_critical_tech = [MAX_ITEMS]int{
	MI, ML, MA, MA, LS, MI, MA, GV,
	GV, GV, GV, GV, LS, BI, BI, LS,
	LS, LS, LS, LS, LS, LS, LS, LS,
	ML, ML, ML, ML, ML, ML, ML, ML,
	ML, 99, 99, 99, 99, 99}

var item_tech_requirment = [MAX_ITEMS]int{
	1, 1, 20, 30, 1, 1, 1, 20,
	25, 30, 40, 50, 20, 40, 50, 10,
	20, 30, 40, 50, 60, 70, 80, 90,
	10, 20, 30, 40, 50, 60, 70, 80,
	90, 999, 999, 999, 999, 999}

var ship_abbr = [NUM_SHIP_CLASSES]string{
	"PB", "CT", "ES", "FF", "DD", "CL", "CS",
	"CA", "CC", "BC", "BS", "DN", "SD", "BM",
	"BW", "BR", "BA", "TR"}

var ship_type = [3]string{"", "S", "S"}

var ship_tonnage = [NUM_SHIP_CLASSES]int{
	1, 2, 5, 10, 15, 20, 25,
	30, 35, 40, 45, 50, 55, 60,
	65, 70, 1, 1}

var ship_cost = [NUM_SHIP_CLASSES]int{
	100, 200, 500, 1000, 1500, 2000, 2500,
	3000, 3500, 4000, 4500, 5000, 5500, 6000,
	6500, 7000, 100, 100}

var command_abbr = [NUM_COMMANDS]string{
	"   ", "ALL", "AMB", "ATT", "AUT", "BAS", "BAT", "BUI", "CON",
	"DEE", "DES", "DEV", "DIS", "END", "ENE", "ENG", "EST", "HAV",
	"HID", "HIJ", "IBU", "ICO", "INS", "INT", "JUM", "LAN", "MES",
	"MOV", "NAM", "NEU", "ORB", "PJU", "PRO", "REC", "REN", "REP",
	"RES", "SCA", "SEN", "SHI", "STA", "SUM", "SUR", "TAR", "TEA",
	"TEC", "TEL", "TER", "TRA", "UNL", "UPG", "VIS", "WIT", "WOR",
	"ZZZ"}

var command_name = [NUM_COMMANDS]string{
	"Undefined", "Ally", "Ambush", "Attack", "Auto", "Base",
	"Battle", "Build", "Continue", "Deep", "Destroy", "Develop",
	"Disband", "End", "Enemy", "Engage", "Estimate", "Haven",
	"Hide", "Hijack", "Ibuild", "Icontinue", "Install", "Intercept",
	"Jump", "Land", "Message", "Move", "Name", "Neutral", "Orbit",
	"Pjump", "Production", "Recycle", "Rename", "Repair", "Research",
	"Scan", "Send", "Shipyard", "Start", "Summary", "Surrender", "Target",
	"Teach", "Tech", "Telescope", "Terraform", "Transfer", "Unload",
	"Upgrade", "Visited", "Withdraw", "Wormhole", "ZZZ"}

// constants from combat.h

/* Maximum number of battle locations for all players. */
const MAX_BATTLES = 50

/* Maximum number of ships at a single battle. */
const MAX_SHIPS = 200

/* Maximum number of engagement options that a player may specify for a single battle. */
const MAX_ENGAGE_OPTIONS = 20

/* Types of combatants. */
const SHIP = 1
const NAMPLA = 2
const GENOCIDE_NAMPLA = 3
const BESIEGED_NAMPLA = 4

/* Types of special targets. */
const TARGET_WARSHIPS = 1
const TARGET_TRANSPORTS = 2
const TARGET_STARBASES = 3
const TARGET_PDS = 4

/* Types of actions. */
const DEFENSE_IN_PLACE = 0
const DEEP_SPACE_DEFENSE = 1
const PLANET_DEFENSE = 2
const DEEP_SPACE_FIGHT = 3
const PLANET_ATTACK = 4
const PLANET_BOMBARDMENT = 5
const GERM_WARFARE = 6
const SIEGE = 7

/* Special types. */
const NON_COMBATANT = 1
