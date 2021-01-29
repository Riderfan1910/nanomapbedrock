import csv
import os

from PIL import Image

currentDir = os.path.dirname(__file__)
textureDir = os.path.join(currentDir, 'textures/blocks')

with open(os.path.join(currentDir, 'blocknames.csv'), 'r') as csvfile:
    blocknames = [b.strip() for b in csvfile.readlines()]

blocks = {}

def get_texturecolors():
    texturecolors = {}

    for imgfile in sorted(os.listdir(textureDir)):
        if imgfile[-4:] != '.png':
            continue

        name = imgfile.split('.')[0]

        if name not in blocknames:
            texturecolors[imgfile[:-4]] = (78, 118, 42, 255)
            continue
    
        img = Image.open(os.path.join(textureDir, imgfile)).convert('RGBA')
        px = img.load()
        pixels = []
        for x in range(img.width):
            for y in range(img.height):
                if px[x, y][3] > 0:
                    pixels.append(px[x, y])

        if len(pixels) == 0:
            color = (0, 0, 0, 0)
        else:
            color = (
                int(sum(p[0] for p in pixels) / len(pixels)),
                int(sum(p[1] for p in pixels) / len(pixels)),
                int(sum(p[2] for p in pixels) / len(pixels)),
                int(sum(p[3] for p in pixels) / len(pixels)),
            )
        
        texturecolors[imgfile[:-4]] = color
    
    return texturecolors

texturecolors = get_texturecolors()
with open(os.path.join(currentDir, 'texturecolors.csv'), 'w') as csvfile:
    writer = csv.writer(csvfile)

    for texture, color in texturecolors.items():
        writer.writerow([texture, color])