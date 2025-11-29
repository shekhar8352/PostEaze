# Layout System Documentation

## Overview

The PostEaze layout system uses Mantine's AppShell component to provide a consistent, responsive layout with sidebar navigation and nested channel navigation.

## Structure

```
src/features/layout/
├── components/
│   ├── MainLayout.tsx       # AppShell wrapper
│   ├── MainSidebar.tsx      # Main navigation sidebar
│   ├── ChannelsNav.tsx      # Collapsible channels navigation
│   └── AppHeader.tsx        # Top header with user menu
├── types.ts                 # TypeScript interfaces
├── constants.ts             # Navigation configuration
└── index.tsx                # Feature exports
```

## Adding New Navigation Items

### 1. Main Navigation Items

Edit `src/features/layout/constants.ts`:

```typescript
import { IconNewFeature } from '@tabler/icons-react';

export const MAIN_NAV_ITEMS: NavItem[] = [
  {
    label: 'Dashboard',
    icon: IconHome,
    path: '/dashboard',
  },
  {
    label: 'New Feature',  // Add your new item
    icon: IconNewFeature,
    path: '/new-feature',
  },
  // ... other items
];
```

### 2. Adding New Channels

To add a new social media channel (e.g., Twitter):

**Step 1:** Add icon import and navigation item in `constants.ts`:

```typescript
import { IconBrandTwitter } from '@tabler/icons-react';

export const CHANNEL_NAV_ITEMS: NavItem[] = [
  // ... existing channels
  {
    label: 'Twitter',
    icon: IconBrandTwitter,
    path: '/channels/twitter',
  },
];
```

**Step 2:** Create the channel page:

```typescript
// src/features/channels/pages/TwitterChannel.tsx
import { Container, Title, Text, Paper, Stack, Button } from '@mantine/core';
import { IconBrandTwitter, IconPlus } from '@tabler/icons-react';

const TwitterChannel = () => {
  return (
    <Container size="xl">
      <Stack gap="xl">
        <Title order={1}>Twitter</Title>
        {/* Your channel content */}
      </Stack>
    </Container>
  );
};

export default TwitterChannel;
```

**Step 3:** Add route in `channelRoutes.tsx`:

```typescript
const TwitterChannel = lazy(() => import('./pages/TwitterChannel'));

export const channelRoutes = [
  // ... existing routes
  {
    path: '/channels/twitter',
    element: <TwitterChannel />,
  },
];
```

### 3. Bottom Navigation Items

For settings, profile, etc.:

```typescript
export const BOTTOM_NAV_ITEMS: NavItem[] = [
  {
    label: 'Settings',
    icon: IconSettings,
    path: '/settings',
  },
  {
    label: 'Profile',
    icon: IconUser,
    path: '/profile',
  },
];
```

## Customizing the Layout

### Sidebar Width

Edit `constants.ts`:

```typescript
export const SIDEBAR_WIDTH = 280; // Change width in pixels
```

### Header Height

Edit `MainLayout.tsx`:

```typescript
<AppShell
  header={{ height: 60 }} // Change height here
  // ...
/>
```

### Sidebar Behavior

**Desktop Sidebar (always visible by default):**

```typescript
const [desktopOpened, { toggle: toggleDesktop }] = useDisclosure(true); // true = open by default
```

**Mobile Sidebar (hidden by default):**

```typescript
const [mobileOpened, { toggle: toggleMobile }] = useDisclosure(); // false = closed by default
```

## Responsive Breakpoints

The layout uses Mantine's responsive system:

- **Mobile**: `< 768px` (sm breakpoint)
- **Desktop**: `>= 768px`

To change the breakpoint, edit `MainLayout.tsx`:

```typescript
<AppShell
  navbar={{
    width: 280,
    breakpoint: 'md', // Change to 'md', 'lg', etc.
    collapsed: { mobile: !mobileOpened, desktop: !desktopOpened },
  }}
/>
```

## Styling

### Active Link Highlighting

Active links are automatically highlighted using React Router's `useLocation` hook. The logic is in `MainSidebar.tsx`:

```typescript
const isActive = location.pathname === item.path;
```

### Custom Styling

For custom styles, use Mantine's styling system or CSS modules:

```typescript
// Using inline styles
<NavLink
  style={{ backgroundColor: 'custom-color' }}
/>

// Using CSS modules
import classes from './MainSidebar.module.css';
<NavLink className={classes.customClass} />
```

## User Menu

The user menu in the header shows:
- User avatar and name
- Profile link
- Settings link
- Logout button

To add new menu items, edit `AppHeader.tsx`:

```typescript
<Menu.Dropdown>
  <Menu.Label>Account</Menu.Label>
  <Menu.Item
    leftSection={<IconNewItem />}
    onClick={() => navigate('/new-item')}
  >
    New Item
  </Menu.Item>
  {/* ... existing items */}
</Menu.Dropdown>
```

## Protected Routes

All routes using MainLayout are automatically protected by `ProtectedLayout`. The routing structure is:

```
ProtectedLayout (checks auth)
  └── MainLayout (shows sidebar/header)
      └── Your Pages
```

This ensures users must be logged in to access any page with the sidebar layout.

## Best Practices

1. **Keep navigation items organized** - Group related items together
2. **Use descriptive labels** - Make navigation clear
3. **Choose appropriate icons** - Use Tabler Icons for consistency
4. **Test responsive behavior** - Always check mobile view
5. **Update types** - Keep TypeScript interfaces in sync
6. **Follow the pattern** - Use existing components as templates

## Troubleshooting

### Sidebar not showing
- Check if route is wrapped with MainLayout
- Verify authentication is working

### Navigation not highlighting
- Ensure `path` in constants matches route exactly
- Check `useLocation` hook is working

### Mobile menu not working
- Verify burger button is visible on mobile
- Check `mobileOpened` state is toggling

### Icons not displaying
- Ensure icon is imported from `@tabler/icons-react`
- Check icon name is correct
