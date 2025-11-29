# Theme System

A modular, scalable theme configuration for the PostEaze application.

## Structure

```
theme/
├── index.ts          # Main theme export
├── typography.ts     # Font families, sizes, weights, headings
├── colors.ts         # Brand, semantic, and channel colors
├── spacing.ts        # Spacing scale, layout sizes, borders, shadows
└── components.ts     # Mantine component overrides
```

## Usage

### Importing the Theme

The theme is automatically applied via `MantineProvider`. No action needed.

### Using Theme Constants in Components

Import specific constants from the theme modules:

```tsx
import { FONT_WEIGHTS, FONT_SIZES } from '@/app/theme';
import { CHANNEL_COLORS, GRADIENTS } from '@/app/theme';
import { SPACING, LAYOUT_SIZES } from '@/app/theme';

// Use in components
<Text fw={FONT_WEIGHTS.bold} size={FONT_SIZES.lg}>
  Bold large text
</Text>

<Box p={SPACING.lg} style={{ maxWidth: LAYOUT_SIZES.maxContentWidth }}>
  Content
</Box>
```

### Typography

**Font Families:**
- `FONT_FAMILIES.body` - Inter (body text)
- `FONT_FAMILIES.heading` - Plus Jakarta Sans (headings)
- `FONT_FAMILIES.mono` - Fira Code (code)

**Font Weights:**
- `light` (300), `regular` (400), `medium` (500)
- `semibold` (600), `bold` (700), `extrabold` (800)

**Font Sizes:**
- `xs` (12px), `sm` (14px), `md` (16px), `lg` (18px), `xl` (20px)

**Heading Sizes:**
- `h1` through `h6` with predefined sizes and weights

### Colors

**Brand Colors:**
- `BRAND_COLORS.primary` - indigo
- `BRAND_COLORS.secondary` - cyan
- `BRAND_COLORS.accent` - blue

**Semantic Colors:**
- `SEMANTIC_COLORS.success` - green
- `SEMANTIC_COLORS.warning` - yellow
- `SEMANTIC_COLORS.error` - red
- `SEMANTIC_COLORS.info` - blue

**Channel Colors:**
- `CHANNEL_COLORS.instagram` - gradient and solid
- `CHANNEL_COLORS.facebook` - solid
- `CHANNEL_COLORS.youtube` - solid
- And more...

**Gradients:**
- `GRADIENTS.primary`, `success`, `warning`, `error`

### Spacing & Layout

**Spacing Scale:**
- `xs` (8px), `sm` (12px), `md` (16px)
- `lg` (24px), `xl` (32px), `xxl` (48px)

**Layout Sizes:**
- `LAYOUT_SIZES.headerHeight` - 60px
- `LAYOUT_SIZES.sidebarWidth` - 280px
- `LAYOUT_SIZES.sidebarCollapsedWidth` - 80px
- `LAYOUT_SIZES.maxContentWidth` - 1200px

**Border Radius:**
- `xs` (4px) through `xl` (24px), `full` (9999px)

**Shadows:**
- `xs`, `sm`, `md`, `lg`, `xl`

**Z-Index:**
- Predefined z-index values for consistent layering

### Icons

**Centralized Icon Management:**
All icons are imported from `@tabler/icons-react` and re-exported through the `Icons` object.

**Usage:**
```tsx
import { Icons } from '@/app/theme';

<Icons.Home size={20} />
<Icons.Instagram size={40} />
<Icons.Plus size={16} />
```

**Available Icon Categories:**
- **Navigation**: Home, ChartBar, Calendar, Settings
- **Social Media**: Instagram, Facebook, YouTube, Twitter, LinkedIn
- **Actions**: Plus, Logout, Refresh, Check, X, Edit, Trash, Download, Upload, Search
- **UI**: ChevronDown, ChevronRight, ChevronLeft, ChevronUp, ArrowLeft, ArrowRight, Menu
- **Status**: AlertCircle, AlertTriangle, InfoCircle, CheckCircle, XCircle
- **User**: User, Mail, Lock, Eye, EyeOff
- **Content**: Photo, Video, File, FileText

**Icon Sizes:**
```tsx
import { ICON_SIZES } from '@/app/theme';

// Predefined sizes
ICON_SIZES.xs  // 14px
ICON_SIZES.sm  // 16px
ICON_SIZES.md  // 20px (default)
ICON_SIZES.lg  // 24px
ICON_SIZES.xl  // 28px
```

## Making Changes

### Updating Typography

Edit `typography.ts`:
```ts
export const FONT_FAMILIES = {
  body: "'Your Font', sans-serif",
  // ...
};
```

Changes will reflect across the entire app.

### Updating Colors

Edit `colors.ts`:
```ts
export const BRAND_COLORS = {
  primary: 'violet', // Change primary color
  // ...
};
```

### Adding New Icons

1. Import the icon from `@tabler/icons-react` in `icons.ts`
2. Add to appropriate category in the imports
3. Export in the `Icons` object
4. Use across the app via `Icons.YourIcon`

### Adding New Design Tokens

1. Add to appropriate module (e.g., `colors.ts`)
2. Export the constant
3. Use in components via import

### Component Defaults

Edit `components.ts` to change default props for Mantine components:
```ts
export const COMPONENT_OVERRIDES = {
  Button: {
    defaultProps: {
      radius: 'lg', // All buttons now have large radius
    },
  },
};
```

## Benefits

✅ **Single Source of Truth** - Change once, update everywhere  
✅ **Type Safety** - TypeScript autocomplete for all values  
✅ **Consistency** - Prevents hardcoded values  
✅ **Scalability** - Easy to extend with new tokens  
✅ **Maintainability** - Clear organization  
✅ **Future-Proof** - Ready for dark mode, themes, etc.

## Best Practices

1. **Always import from theme** - Never hardcode colors, sizes, or spacing
2. **Use semantic names** - Prefer `SEMANTIC_COLORS.error` over `'red'`
3. **Extend, don't modify** - Add new constants rather than changing existing ones
4. **Document changes** - Update this README when adding new tokens
