# Public Assets

This folder contains static assets that are served directly by Vite during development and copied to the build output during production builds. These assets are publicly accessible and can be referenced using absolute paths from the application root.

## Contents

- **vite.svg**: Default Vite logo used as the application favicon

## Asset Management

### Static File Serving

Assets in the `public` folder are:
- Served directly at the root path during development (`/vite.svg`)
- Copied to the build output directory during production builds
- Accessible via absolute URLs without import statements
- Not processed by Vite's build pipeline (no bundling, minification, or hash generation)

### Usage Patterns

#### Referencing Public Assets
```html
<!-- In HTML files -->
<link rel="icon" type="image/svg+xml" href="/vite.svg" />
<img src="/logo.png" alt="Logo" />
```

```tsx
// In React components
<img src="/vite.svg" alt="Vite logo" />
```

#### When to Use Public Assets
Use the `public` folder for:
- **Favicons and app icons**: Files referenced in HTML meta tags
- **Static images**: Images that don't need processing or optimization
- **Third-party assets**: External libraries or resources
- **SEO assets**: robots.txt, sitemap.xml, etc.
- **PWA assets**: Service worker files, manifest.json

#### When NOT to Use Public Assets
Avoid the `public` folder for:
- **Component assets**: Images imported in React components (use `src/assets` instead)
- **Processed assets**: Files that need optimization, resizing, or bundling
- **Dynamic imports**: Assets that should be code-split or lazy-loaded

## Asset Organization

### Recommended Structure
```
public/
├── favicon.ico          # Browser favicon
├── logo.png            # Application logo
├── icons/              # App icons for different platforms
│   ├── icon-192.png
│   └── icon-512.png
├── images/             # Static images
│   └── hero-bg.jpg
└── manifest.json       # PWA manifest (if applicable)
```

### File Naming Conventions
- Use lowercase filenames with hyphens for multi-word names
- Include dimensions in icon filenames (e.g., `icon-192.png`)
- Use descriptive names that indicate the asset's purpose

## Build Process

During the build process (`npm run build`):
1. Vite copies all files from `public/` to the build output directory
2. Files maintain their original names and paths
3. No processing, optimization, or cache-busting hashes are applied
4. Assets remain accessible at their original paths in production

## Performance Considerations

- **File Size**: Keep public assets optimized since they're not processed by Vite
- **Caching**: Consider manual cache-busting for frequently updated assets
- **CDN**: Large static assets may benefit from CDN hosting
- **Compression**: Ensure images are properly compressed before adding to public folder

## Related Documentation

- [Frontend Assets](../src/assets/README.md) - For processed/imported assets
- [Vite Static Assets Guide](https://vitejs.dev/guide/assets.html#the-public-directory)