EESchema Schematic File Version 4
EELAYER 30 0
EELAYER END
$Descr A4 11693 8268
encoding utf-8
Sheet 1 1
Title "NOT Gate"
Date ""
Rev ""
Comp ""
Comment1 ""
Comment2 ""
Comment3 ""
Comment4 ""
$EndDescr
$Comp
L pspice:R R1
U 1 1 5E93D37E
P 5150 3650
F 0 "R1" V 4945 3650 50  0000 C CNN
F 1 "10K" V 5036 3650 50  0000 C CNN
F 2 "" H 5150 3650 50  0001 C CNN
F 3 "~" H 5150 3650 50  0001 C CNN
	1    5150 3650
	0    1    1    0   
$EndComp
$Comp
L pspice:R R2
U 1 1 5E93E135
P 6400 2950
F 0 "R2" H 6332 2904 50  0000 R CNN
F 1 "1K" H 6332 2995 50  0000 R CNN
F 2 "" H 6400 2950 50  0001 C CNN
F 3 "~" H 6400 2950 50  0001 C CNN
	1    6400 2950
	-1   0    0    1   
$EndComp
$Comp
L pspice:VSOURCE Vcc
U 1 1 5E93FC3B
P 4050 2900
F 0 "Vcc" H 4278 2946 50  0000 L CNN
F 1 "5" H 4278 2855 50  0000 L CNN
F 2 "" H 4050 2900 50  0001 C CNN
F 3 "~" H 4050 2900 50  0001 C CNN
	1    4050 2900
	1    0    0    -1  
$EndComp
$Comp
L pspice:0 #GND?
U 1 1 5E941594
P 4050 3400
F 0 "#GND?" H 4050 3300 50  0001 C CNN
F 1 "0" H 4050 3489 50  0000 C CNN
F 2 "" H 4050 3400 50  0001 C CNN
F 3 "~" H 4050 3400 50  0001 C CNN
	1    4050 3400
	1    0    0    -1  
$EndComp
Wire Wire Line
	4050 3200 4050 3400
Text GLabel 4050 2550 1    50   Input ~ 0
Vcc
Text GLabel 6400 2550 1    50   Input ~ 0
Vcc
Wire Wire Line
	4050 2550 4050 2600
Wire Wire Line
	6400 2550 6400 2700
Wire Wire Line
	6400 3200 6400 3250
Text GLabel 4700 3650 0    50   Input ~ 0
in
Text GLabel 6700 3250 2    50   Output ~ 0
out
Wire Wire Line
	4700 3650 4900 3650
Wire Wire Line
	6400 3250 6700 3250
Connection ~ 6400 3250
$Comp
L pspice:0 #GND?
U 1 1 5E947C4E
P 6400 4400
F 0 "#GND?" H 6400 4300 50  0001 C CNN
F 1 "0" H 6400 4489 50  0000 C CNN
F 2 "" H 6400 4400 50  0001 C CNN
F 3 "~" H 6400 4400 50  0001 C CNN
	1    6400 4400
	1    0    0    -1  
$EndComp
$Comp
L pspice:VSOURCE Vin
U 1 1 5E94AEED
P 4050 4150
F 0 "Vin" H 4278 4196 50  0000 L CNN
F 1 "PULSE(0 5 2MS 2MS 2MS 50MS 100MS) " H 4278 4105 50  0000 L CNN
F 2 "" H 4050 4150 50  0001 C CNN
F 3 "~" H 4050 4150 50  0001 C CNN
	1    4050 4150
	1    0    0    -1  
$EndComp
Text GLabel 4050 3750 1    50   Input ~ 0
in
Wire Wire Line
	4050 3850 4050 3750
$Comp
L pspice:0 #GND?
U 1 1 5E94C724
P 4050 4650
F 0 "#GND?" H 4050 4550 50  0001 C CNN
F 1 "0" H 4050 4739 50  0000 C CNN
F 2 "" H 4050 4650 50  0001 C CNN
F 3 "~" H 4050 4650 50  0001 C CNN
	1    4050 4650
	1    0    0    -1  
$EndComp
Wire Wire Line
	4050 4450 4050 4650
Text Notes 4650 4700 0    50   ~ 0
.tran 0.1s 10s 0
$Comp
L Transistor_BJT:BC546 Q?
U 1 1 5E93E404
P 6300 3650
F 0 "Q?" H 6491 3696 50  0000 L CNN
F 1 "BC546" H 6491 3605 50  0000 L CNN
F 2 "Package_TO_SOT_THT:TO-92_Inline" H 6500 3575 50  0001 L CIN
F 3 "http://www.fairchildsemi.com/ds/BC/BC547.pdf" H 6300 3650 50  0001 L CNN
F 4 "Q" H 6300 3650 50  0001 C CNN "Spice_Primitive"
F 5 "BC546B" H 6300 3650 50  0001 C CNN "Spice_Model"
F 6 "Y" H 6300 3650 50  0001 C CNN "Spice_Netlist_Enabled"
F 7 "BC546.lib" H 6300 3650 50  0001 C CNN "Spice_Lib_File"
	1    6300 3650
	1    0    0    -1  
$EndComp
Wire Wire Line
	6400 3250 6400 3450
Wire Wire Line
	6400 3850 6400 4400
Wire Wire Line
	5400 3650 6100 3650
$EndSCHEMATC
