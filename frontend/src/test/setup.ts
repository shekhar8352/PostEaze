import '@testing-library/jest-dom';

// Mock window.matchMedia which is required by Mantine components
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => {},
  }),
});

// Mock ResizeObserver which might also be needed by Mantine
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// Mock HTMLElement.scrollIntoView which might be called by Mantine components
Element.prototype.scrollIntoView = () => {};