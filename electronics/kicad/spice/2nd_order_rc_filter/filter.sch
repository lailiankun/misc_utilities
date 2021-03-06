EESchema Schematic File Version 4
EELAYER 30 0
EELAYER END
$Descr A4 11693 8268
encoding utf-8
Sheet 1 1
Title "2nd Order RC Filter"
Date ""
Rev ""
Comp ""
Comment1 ""
Comment2 ""
Comment3 ""
Comment4 ""
$EndDescr
$Comp
L pspice:VSOURCE V1
U 1 1 5E913B85
P 3250 2800
F 0 "V1" H 3478 2846 50  0000 L CNN
F 1 "dc 0 ac 1" H 3478 2755 50  0000 L CNN
F 2 "" H 3250 2800 50  0001 C CNN
F 3 "~" H 3250 2800 50  0001 C CNN
	1    3250 2800
	1    0    0    -1  
$EndComp
$Comp
L pspice:0 #GND?
U 1 1 5E91402F
P 3250 3350
F 0 "#GND?" H 3250 3250 50  0001 C CNN
F 1 "0" H 3250 3439 50  0000 C CNN
F 2 "" H 3250 3350 50  0001 C CNN
F 3 "~" H 3250 3350 50  0001 C CNN
	1    3250 3350
	1    0    0    -1  
$EndComp
$Comp
L pspice:0 #GND?
U 1 1 5E914375
P 4350 3400
F 0 "#GND?" H 4350 3300 50  0001 C CNN
F 1 "0" H 4350 3489 50  0000 C CNN
F 2 "" H 4350 3400 50  0001 C CNN
F 3 "~" H 4350 3400 50  0001 C CNN
	1    4350 3400
	1    0    0    -1  
$EndComp
$Comp
L pspice:0 #GND?
U 1 1 5E914972
P 5250 3400
F 0 "#GND?" H 5250 3300 50  0001 C CNN
F 1 "0" H 5250 3489 50  0000 C CNN
F 2 "" H 5250 3400 50  0001 C CNN
F 3 "~" H 5250 3400 50  0001 C CNN
	1    5250 3400
	1    0    0    -1  
$EndComp
$Comp
L pspice:R R1
U 1 1 5E915372
P 4100 2300
F 0 "R1" V 3895 2300 50  0000 C CNN
F 1 "10K" V 3986 2300 50  0000 C CNN
F 2 "" H 4100 2300 50  0001 C CNN
F 3 "~" H 4100 2300 50  0001 C CNN
	1    4100 2300
	0    1    1    0   
$EndComp
$Comp
L pspice:R R2
U 1 1 5E917135
P 5000 2300
F 0 "R2" V 4795 2300 50  0000 C CNN
F 1 "1K" V 4886 2300 50  0000 C CNN
F 2 "" H 5000 2300 50  0001 C CNN
F 3 "~" H 5000 2300 50  0001 C CNN
	1    5000 2300
	0    1    1    0   
$EndComp
$Comp
L pspice:C C2
U 1 1 5E917A24
P 5250 2850
F 0 "C2" H 5428 2896 50  0000 L CNN
F 1 "100n" H 5428 2805 50  0000 L CNN
F 2 "" H 5250 2850 50  0001 C CNN
F 3 "~" H 5250 2850 50  0001 C CNN
	1    5250 2850
	1    0    0    -1  
$EndComp
$Comp
L pspice:C C1
U 1 1 5E9180E8
P 4350 2850
F 0 "C1" H 4528 2896 50  0000 L CNN
F 1 "1u" H 4528 2805 50  0000 L CNN
F 2 "" H 4350 2850 50  0001 C CNN
F 3 "~" H 4350 2850 50  0001 C CNN
	1    4350 2850
	1    0    0    -1  
$EndComp
Wire Wire Line
	3250 3350 3250 3100
Wire Wire Line
	3250 2500 3250 2300
Wire Wire Line
	3250 2300 3850 2300
Wire Wire Line
	4350 2300 4350 2600
Wire Wire Line
	4350 3100 4350 3400
Wire Wire Line
	4350 2300 4750 2300
Connection ~ 4350 2300
Wire Wire Line
	5250 2300 5250 2600
Wire Wire Line
	5250 3100 5250 3400
Text GLabel 5650 2300 2    50   Output ~ 0
OUT
Text GLabel 3000 2300 0    50   Input ~ 0
IN
Wire Wire Line
	3000 2300 3250 2300
Connection ~ 3250 2300
Wire Wire Line
	5250 2300 5650 2300
Connection ~ 5250 2300
Text Notes 3050 3550 0    50   ~ 0
.ac dec 10 1 100k
$EndSCHEMATC
