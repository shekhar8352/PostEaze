# Utils

Utility functions and helper modules for the PostEaze React frontend application.

## Purpose

This folder contains reusable utility functions, helper modules, and common functionality that can be shared across different features and components in the application. These utilities help maintain DRY (Don't Repeat Yourself) principles and provide consistent behavior throughout the frontend.

## Current Status

The utils folder is currently empty but is structured to accommodate various types of utility functions as the application grows.

## Planned Utility Categories

Based on the current application architecture and dependencies, the following utility categories are recommended:

### Form Utilities
- **Validation helpers**: Common validation functions for use with Formik and Yup
- **Form field utilities**: Reusable form field configurations and transformations
- **Error handling**: Standardized form error processing and display

### API Utilities
- **Request helpers**: Common API request patterns and transformations
- **Response processing**: Standardized API response handling and error processing
- **Authentication helpers**: Token management and auth-related utilities

### Data Utilities
- **Formatters**: Date, currency, and text formatting functions
- **Transformers**: Data transformation and normalization utilities
- **Validators**: Client-side validation functions beyond form validation

### UI Utilities
- **Theme helpers**: Mantine theme-related utility functions
- **Responsive utilities**: Breakpoint and responsive design helpers
- **Animation helpers**: Common animation and transition utilities

### Storage Utilities
- **Local storage**: Safe localStorage access with error handling
- **Session storage**: Session management utilities
- **Cache utilities**: Client-side caching helpers

### String Utilities
- **Text processing**: String manipulation and formatting functions
- **URL utilities**: URL parsing, building, and validation
- **Slug generation**: URL-friendly string generation

## Usage Patterns

When creating utilities, follow these patterns:

### File Organization
```
utils/
├── README.md
├── api/
│   ├── index.ts
│   └── request-helpers.ts
├── forms/
│   ├── index.ts
│   └── validation.ts
├── formatters/
│   ├── index.ts
│   ├── date.ts
│   └── currency.ts
└── storage/
    ├── index.ts
    └── local-storage.ts
```

### Export Pattern
Each utility module should have a clear export structure:

```typescript
// utils/formatters/index.ts
export { formatDate, formatRelativeTime } from './date';
export { formatCurrency, formatNumber } from './currency';
```

### Function Naming
- Use descriptive, verb-based names: `formatDate`, `validateEmail`, `parseApiError`
- Prefix boolean utilities with `is`, `has`, or `can`: `isValidEmail`, `hasPermission`
- Use consistent naming patterns within categories

### Type Safety
All utilities should be fully typed with TypeScript:

```typescript
export const formatDate = (date: Date | string, format: string = 'MM/dd/yyyy'): string => {
  // Implementation
};
```

## Integration with Application

### Import Patterns
Utilities should be imported using absolute paths with the `@` alias:

```typescript
import { formatDate, validateEmail } from '@/utils/formatters';
import { apiRequest } from '@/utils/api';
```

### Testing
Each utility function should have corresponding unit tests in the `test/utils` directory:

```typescript
// test/utils/formatters.test.ts
import { formatDate } from '@/utils/formatters';

describe('formatDate', () => {
  it('should format date correctly', () => {
    // Test implementation
  });
});
```

## Dependencies

Utilities may leverage the following project dependencies:
- **Axios**: For API-related utilities
- **Mantine**: For UI and theme-related utilities
- **React Router**: For navigation and URL utilities
- **Redux Toolkit**: For state-related utilities

## Best Practices

1. **Pure Functions**: Utilities should be pure functions when possible, avoiding side effects
2. **Error Handling**: Include proper error handling and fallback values
3. **Documentation**: Each utility should have JSDoc comments explaining purpose and usage
4. **Performance**: Consider performance implications, especially for frequently used utilities
5. **Reusability**: Design utilities to be generic and reusable across different contexts

## Examples

### API Utility Example
```typescript
// utils/api/request-helpers.ts
import axiosInstance from '@/services/axios';

export const handleApiError = (error: any): string => {
  if (error.response?.data?.message) {
    return error.response.data.message;
  }
  return 'An unexpected error occurred';
};
```

### Form Utility Example
```typescript
// utils/forms/validation.ts
import * as Yup from 'yup';

export const emailValidation = Yup.string()
  .email('Invalid email format')
  .required('Email is required');

export const passwordValidation = Yup.string()
  .min(8, 'Password must be at least 8 characters')
  .required('Password is required');
```

## Related Documentation

- [Services Documentation](../services/README.md) - API integration patterns
- [Features Documentation](../features/README.md) - Feature-specific utilities
- [Test Utils Documentation](../test/utils/README.md) - Testing utility functions
- [App Configuration](../app/README.md) - Application-level utilities

## Contributing

When adding new utilities:
1. Follow the established patterns and naming conventions
2. Include comprehensive TypeScript types
3. Add unit tests for all functions
4. Update this README with new utility categories
5. Consider if the utility belongs in a feature-specific location instead