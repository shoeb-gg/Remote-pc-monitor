import sharp from 'sharp';
import { readFileSync } from 'fs';

const svg = readFileSync('./static/icon.svg');

// Generate 192x192 icon
await sharp(svg)
  .resize(192, 192)
  .png()
  .toFile('./static/icon-192.png');

console.log('✓ Generated icon-192.png');

// Generate 512x512 icon
await sharp(svg)
  .resize(512, 512)
  .png()
  .toFile('./static/icon-512.png');

console.log('✓ Generated icon-512.png');

// Generate favicon
await sharp(svg)
  .resize(32, 32)
  .png()
  .toFile('./static/favicon.png');

console.log('✓ Generated favicon.png');

console.log('\n✅ All icons generated successfully!');
