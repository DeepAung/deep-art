from PIL import Image  
import random
import sys
import glob
import os

def to_hex(r: int, g: int, b: int) -> str:
    return '#%02x%02x%02x' % (r, g, b)


width = 400
height = 300

n = 1
if len(sys.argv) > 1:
    n = int(sys.argv[1])

files = glob.glob("./images/*")
for file in files:
    print(file)
    os.remove(file)

for i in range(n):
    r, g, b = random.randint(0, 255), random.randint(0, 255), random.randint(0, 255)
    hex_color = to_hex(r, g, b)

    img  = Image.new( mode = "RGB", size = (width, height), color = (r, g, b) )
    img.save(f"./images/{hex_color}.png")
