# Based on https://www.bidouille.org/prog/plasma
from math import *
from time import *

W = 320
H = 200

cc = 0
pc = 0

def setup():
    size(W, H)

def draw():
    for y in range(H):
        for x in range(W):
            col = gencolor(x, y)
            stroke(col)
            strokeWeight(5)
            point(x, y)

def gencolor(x, y):
    global cc
    global pc
    
    v = 0
    x = map(x, 0, W, -0.5, 0.5)
    y = map(y, 0, H, -0.5, 0.5)
    t = clock()
    col = color(0, 0, 0)

    if pc == 0:
        v = sin(x*10 + t)
    elif pc == 1:
        v = sin(10*(x*sin(t/2) + y*cos(t/3)) + t)
    elif pc == 2:
        cx = x + .5*sin(t/5)
        cy = y + .5*cos(t/3)
        v = sin(sqrt(100*(cx**2 + cy**2) + 1) + t)
    elif pc == 3:
        v += sin(x*10 + t)
        
        v += sin(10*(x*sin(t/2) + y*cos(t/3)) + t)
        
        cx = x + .5*sin(t/5)
        cy = y + .5*cos(t/3)
        v += sin(sqrt(100*(cx**2 + cy**2) + 1) + t)
        

    if cc == 0:
        r = sin(v*PI)
        g = cos(v*PI)
        b = 0
    elif cc == 1:
        r = 1
        g = cos(v*PI)
        b = sin(v*PI)
    elif cc == 2:
        r = sin(v*PI)
        g = sin(v*PI + 2*PI/3)
        b = sin(v*PI + 4*PI/3)
    elif cc == 3:
        r = sin(v*5*PI)
        g = r
        b = r
    
    r = map(r, 0, 1, 0, 255)
    g = map(g, 0, 1, 0, 255)
    b = map(b, 0, 1, 0, 255)
    col = color(r, g, b)
    
    return col
        
def keyPressed():
    global pc
    global cc
    key = chr(keyCode)
    if key == 'A':
        pc = (pc - 1) % 4
    elif key == 'S':
        pc = (pc + 1) % 4

    if key == 'Q':
        cc = (cc - 1) % 4
    elif key == 'W':
        cc = (cc + 1) % 4
        
