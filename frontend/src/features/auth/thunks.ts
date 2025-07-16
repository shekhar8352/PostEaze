import { createAsyncThunk } from '@reduxjs/toolkit';
import { signup } from './authApi';

export const signupUser = createAsyncThunk('auth/signup', signup);