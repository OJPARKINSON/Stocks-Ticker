import glob
import os
import time
import json
from sys import exit

from font_fredoka_one import FredokaOne
from inky import inky
from PIL import Image, ImageDraw, ImageFont

try:
    import requests
except ImportError:
    exit("This script requires the requests module\nInstall with: sudo pip install requests")

try:
    import geocoder
except ImportError:
    exit("This script requires the geocoder module\nInstall with: sudo pip install geocoder")

print("Inky pHAT: Ticker \n")

PATH = os.path.dirname(__file__)

inky_display = inky.Inky(resolution=(212, 104), colour='red', h_flip=True, v_flip=True)

inky_display.set_border(inky_display.BLACK)
font = ImageFont.truetype(FredokaOne, 22)

response = requests.get('https://rgy1uko64b.execute-api.eu-west-2.amazonaws.com/default/Stonks')
prices = json.loads(response.text)

# Placeholder variables
cmcsaPrice = 0
cmcsaProfit = 0

if prices:
    cmcsaPrice = prices["cmcsaPrice"]
    cmcsaProfit = prices["cmcsaProf"]
else:
    print("Warning, no stock information found!")

# Create a new canvas to draw on
img = Image.open(os.path.join(PATH, "phat/resources/backdrop.png")).resize(inky_display.resolution)
draw = ImageDraw.Draw(img)

# Draw lines to frame the weather data
draw.line((31, 35, 184, 35))      # Horizontal top line
draw.line((169, 58, 169, 58), 2)  # Red seaweed pixel :D

draw.text((70, 12), "CMCSA", inky_display.WHITE, font=font)

draw.text((75, 34), cmcsaPrice, inky_display.WHITE, font=font)
draw.text((55, 54), cmcsaProfit, inky_display.WHITE, font=font)

inky_display.set_image(img)
inky_display.show()