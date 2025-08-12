# Frontend Assets

This folder contains static assets that are imported and processed by Vite's build pipeline. These assets are optimized, bundled, and receive cache-busting hashes during the build process, making them ideal for component-level imports and dynamic asset loading.

## Contents

- **react.svg**: React logo asset used in components (currently available but not actively used)

## Asset Management

### Processed Asset Pipeline

Assets in the `src/assets` folder are:
- Imported using ES6 import statements in React components
- Processed by Vite's build pipeline (optimization, minification, hash generation)
- Bundled with the application code for optimal loading
- Automatically optimized for different formats and sizes when possible
- Given cache-busting hashes in production builds

### Usage Patterns

#### Importing Assets in Components
```tsx
// Import the asset
import logoImage from '../assets/logo.png';
import iconSvg from '../assets/icons/user-icon.svg';

// Use in JSX
const MyComponent = () => (
  <div>
    <img src={logoImage} alt="Company Logo" />
    <img src={iconSvg} alt="User Icon" />
  </div>
);
```

#### Dynamic Imports
```tsx
// For code-splitting or lazy loading
const loadAsset = async () => {
  const { default: dynamicImage } = await import('../assets/large-image.jpg');
  return dynamicImage;
};
```

#### CSS Imports
```css
/* In CSS files */
.hero-section {
  background-image: url('./assets/hero-background.jpg');
}
```

### When to Use src/assets

Use the `src/assets` folder for:
- **Component images**: Images directly used in React components
- **Icons and graphics**: SVG icons, logos, and graphics that need optimization
- **Background images**: Images referenced in CSS files
- **Dynamic assets**: Images that might be conditionally loaded
- **Optimizable content**: Assets that benefit from Vite's processing pipeline

### When NOT to Use src/assets

Avoid the `src/assets` folder for:
- **Static public files**: Use `public/` for files that don't need processing
- **Large media files**: Consider CDN hosting for very large assets
- **Third-party assets**: External resources should typically go in `public/`
- **SEO files**: robots.txt, sitemap.xml belong in `public/`

## Asset Organization

### Recommended Structure
```
src/assets/
├── images/             # General images
│   ├── logos/         # Company and brand logos
│   ├── backgrounds/   # Background images
│   └── illustrations/ # Custom illustrations
├── icons/             # Icon assets
│   ├── ui/           # UI icons (buttons, navigation)
│   ├── social/       # Social media icons
│   └── status/       # Status and state icons
├── fonts/            # Custom font files (if not using CDN)
└── data/             # Static data files (JSON, etc.)
```

### File Naming Conventions
- Use lowercase filenames with hyphens for multi-word names
- Group related assets in subfolders
- Use descriptive names that indicate the asset's purpose
- Include size or variant information when relevant (e.g., `logo-small.png`)

## Build Process

During the build process (`npm run build`):
1. Vite processes and optimizes all imported assets
2. Assets receive unique hash-based filenames for cache-busting
3. Unused assets are automatically excluded from the build
4. Images may be converted to more efficient formats
5. Assets are bundled and served from the optimized build directory

## Performance Considerations

- **Import only what you need**: Unused imports are tree-shaken out
- **Image optimization**: Consider using modern formats (WebP, AVIF) when supported
- **Lazy loading**: Use dynamic imports for large assets that aren't immediately needed
- **Asset size**: Keep individual assets reasonably sized; consider compression
- **Bundle splitting**: Large assets may be split into separate chunks automatically

## TypeScript Support

Vite provides built-in TypeScript support for asset imports:

```tsx
// TypeScript will recognize these imports
import imageUrl from './image.png';  // string
import svgUrl from './icon.svg';     // string

// For SVG as React components (requires additional setup)
import { ReactComponent as Icon } from './icon.svg';
```

## Related Documentation

- [Frontend Public Assets](../../public/README.md) - For static, unprocessed assets
- [Frontend Utils](../utils/README.md) - For asset-related utility functions
- [Vite Asset Handling Guide](https://vitejs.dev/guide/assets.html)