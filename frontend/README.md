# PostEaze Frontend

The PostEaze frontend is a modern React application built with TypeScript, providing a responsive and intuitive user interface for social media management. The application uses a feature-based architecture with Redux Toolkit for state management and Mantine for UI components.

## Architecture

### Tech Stack

- **React 19** - Modern React with latest features and performance improvements
- **TypeScript** - Type-safe development with full IntelliSense support
- **Vite** - Fast build tool with hot module replacement (HMR)
- **Mantine** - Comprehensive React components library with built-in theming
- **Redux Toolkit** - Predictable state management with modern Redux patterns
- **React Router** - Declarative routing for single-page application navigation
- **Axios** - HTTP client for API communication with the Go backend
- **Formik + Yup** - Form handling and validation
- **Vitest** - Fast unit testing framework

### Application Structure

```
src/
├── app/           # Redux store configuration and app-level setup
├── features/      # Feature-based modules (auth, landing, etc.)
├── routes/        # Application routing configuration
├── services/      # API services and HTTP client setup
├── utils/         # Shared utility functions
├── assets/        # Static assets (images, icons, etc.)
├── test/          # Testing utilities and setup
├── App.tsx        # Root application component
└── main.tsx       # Application entry point
```

## Key Features

### State Management
- **Redux Toolkit** for predictable state management
- Feature-based slice organization
- Async thunks for API integration
- Type-safe hooks for React-Redux integration

### UI Framework
- **Mantine** components with custom theming
- Responsive design with mobile-first approach
- Consistent design system across all features
- Built-in accessibility features

### Routing
- **React Router v7** for client-side routing
- Feature-based route organization
- Protected routes for authenticated areas
- Lazy loading for code splitting

### Form Handling
- **Formik** for form state management
- **Yup** for schema validation
- Reusable form components
- Error handling and user feedback

## Getting Started

### Prerequisites
- Node.js 18+ and npm/yarn
- PostEaze backend running on port 8080

### Development Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Configuration

The application connects to the backend API. Ensure the backend is running and accessible at the configured endpoint.

## Development Guidelines

### Feature Organization
Each feature follows a consistent structure:
- Components and pages
- Redux slice for state management
- API service functions
- Route definitions
- Tests

### Code Style
- TypeScript strict mode enabled
- ESLint configuration for code quality
- Consistent naming conventions
- Component composition patterns

### Testing Strategy
- Unit tests with Vitest and React Testing Library
- Component testing with user interaction simulation
- Redux store testing with mock store
- API service testing with mocked responses

## Related Documentation

- [Source Code Structure](src/README.md) - Detailed source code organization
- [Features Architecture](src/features/README.md) - Feature-based development patterns
- [State Management](src/app/README.md) - Redux store configuration
- [API Services](src/services/README.md) - Backend integration patterns
- [Testing Setup](src/test/README.md) - Testing utilities and configuration
