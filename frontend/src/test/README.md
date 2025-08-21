# Test Setup

Testing configuration and setup for the PostEaze React frontend application using Vitest, React Testing Library, and Jest DOM.

## Purpose

This directory contains the core testing setup, configuration, and utilities needed to test React components in the PostEaze frontend. It provides a comprehensive testing environment with proper mocking for Mantine UI components and Redux store integration.

## Testing Framework

The frontend uses **Vitest** as the test runner with the following key libraries:

- **Vitest**: Fast unit test framework with native ES modules support
- **React Testing Library**: Component testing utilities focused on user interactions
- **Jest DOM**: Custom matchers for DOM assertions
- **JSDOM**: Browser environment simulation for Node.js testing

## Configuration

### Vitest Configuration

Testing is configured in `vite.config.ts` with the following settings:

```typescript
test: {
  globals: true,           // Enable global test functions (describe, it, expect)
  environment: 'jsdom',    // Browser-like environment for React components
  setupFiles: './src/test/setup.ts',  // Global test setup file
}
```

### Test Setup File

The `setup.ts` file provides essential mocking and configuration:

#### Jest DOM Integration
```typescript
import '@testing-library/jest-dom';
```
Enables custom DOM matchers like `toBeInTheDocument()`, `toHaveClass()`, etc.

#### Mantine UI Mocking
- **matchMedia**: Mocks `window.matchMedia` for responsive components
- **ResizeObserver**: Mocks ResizeObserver API used by Mantine components
- **scrollIntoView**: Mocks scroll behavior for component interactions

## Available Test Utilities

### Core Utilities

- **[renderWithStore](./utils/README.md)**: Primary utility for testing Redux-connected components with Mantine provider setup

## Testing Patterns

### Component Testing Structure

```typescript
import { renderWithStore } from '@/test/utils/renderWithStore';
import { screen, fireEvent } from '@testing-library/react';
import { MyComponent } from './MyComponent';

describe('MyComponent', () => {
  test('renders correctly', () => {
    renderWithStore(<MyComponent />);
    expect(screen.getByText('Expected Text')).toBeInTheDocument();
  });

  test('handles user interactions', async () => {
    const { store } = renderWithStore(<MyComponent />);
    
    const button = screen.getByRole('button');
    fireEvent.click(button);
    
    // Assert component behavior or store state changes
    expect(store.getState().someSlice.someValue).toBe(expectedValue);
  });
});
```

### Testing Connected Components

For components that use Redux state or dispatch actions:

```typescript
import { renderWithStore } from '@/test/utils/renderWithStore';
import { ConnectedComponent } from './ConnectedComponent';

test('renders with auth state', () => {
  const preloadedState = {
    auth: {
      isAuthenticated: true,
      user: { id: 1, email: 'test@example.com' }
    }
  };

  renderWithStore(<ConnectedComponent />, { preloadedState });
  expect(screen.getByText('Welcome back!')).toBeInTheDocument();
});
```

### Testing Mantine Components

Components using Mantine UI work seamlessly with the provided setup:

```typescript
import { renderWithStore } from '@/test/utils/renderWithStore';
import { Button } from '@mantine/core';

test('mantine button renders correctly', () => {
  renderWithStore(<Button>Click me</Button>);
  expect(screen.getByRole('button')).toHaveTextContent('Click me');
});
```

## Running Tests

### Available Commands

```bash
# Run all tests
npm run test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run specific test file
npm run test MyComponent.test.tsx
```

### Test File Conventions

- **Location**: Place test files next to the components they test
- **Naming**: Use `.test.tsx` or `.spec.tsx` extensions
- **Structure**: Group related tests using `describe` blocks

## Mocking Strategies

### API Calls

For components that make API calls, mock the service layer:

```typescript
import { vi } from 'vitest';
import * as authService from '@/services/authService';

vi.mock('@/services/authService', () => ({
  login: vi.fn(),
  logout: vi.fn(),
}));
```

### Router Navigation

For components using React Router:

```typescript
import { vi } from 'vitest';

const mockNavigate = vi.fn();
vi.mock('react-router-dom', () => ({
  ...vi.importActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
```

## Best Practices

### Test Organization
1. **Group Related Tests**: Use `describe` blocks to organize tests by component or feature
2. **Clear Test Names**: Write descriptive test names that explain the expected behavior
3. **Setup and Teardown**: Use `beforeEach` and `afterEach` for common setup/cleanup

### Component Testing
1. **Test User Interactions**: Focus on how users interact with components
2. **Test Props and State**: Verify components respond correctly to different props and state
3. **Test Error States**: Include tests for error conditions and edge cases

### Redux Testing
1. **Test State Changes**: Verify that actions properly update the store state
2. **Test Component Integration**: Ensure components correctly read from and dispatch to the store
3. **Use Preloaded State**: Set up specific state scenarios for comprehensive testing

### Performance
1. **Mock Heavy Dependencies**: Mock complex external libraries when possible
2. **Isolate Components**: Test components in isolation when feasible
3. **Avoid Over-Mocking**: Only mock what's necessary for the test

## Troubleshooting

### Common Issues

#### Mantine Components Not Rendering
- Ensure you're using `renderWithStore` which provides the Mantine provider
- Check that the component is properly wrapped in the test

#### Redux Store Errors
- Verify that required reducers are included in the test store
- Check that preloaded state matches the expected shape

#### Async Testing Issues
- Use `await` with async operations
- Consider using `waitFor` for elements that appear after async operations

### Debug Tips

```typescript
import { screen } from '@testing-library/react';

// Debug rendered output
screen.debug();

// Debug specific element
screen.debug(screen.getByRole('button'));
```

## Dependencies

### Core Testing Dependencies
- `vitest`: Test runner and framework
- `@testing-library/react`: React component testing utilities
- `@testing-library/jest-dom`: DOM assertion matchers
- `@testing-library/user-event`: User interaction simulation
- `jsdom`: Browser environment simulation

### Application Dependencies
- `@reduxjs/toolkit`: Redux store configuration for tests
- `react-redux`: Redux provider for component tests
- `@mantine/core`: UI components and theming

## Related Documentation

- [Test Utils](./utils/README.md) - Test utility functions and helpers
- [Features](../features/README.md) - Feature-based architecture and testing
- [App Store](../app/README.md) - Redux store configuration
- [Services](../services/README.md) - API service layer and mocking strategies