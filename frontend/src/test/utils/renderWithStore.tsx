// src/test/utils/renderWithStore.tsx
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { render } from '@testing-library/react';
import authReducer from '@/features/auth/authSlice';

export const renderWithStore = (ui: React.ReactElement, { preloadedState = {}, store = configureStore({
  reducer: {
    auth: authReducer,
  },
  preloadedState,
}) } = {}) => {
  return {
    store,
    ...render(<Provider store={store}>{ui}</Provider>),
  };
};
