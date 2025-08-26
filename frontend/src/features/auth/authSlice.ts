import { createSlice } from "@reduxjs/toolkit";

interface AuthState {
  loading: boolean;
  error: string | null;
  user: any;
}

const initialState: AuthState = {
  loading: false,
  error: null,
  user: null,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {},
  extraReducers: () => {},
});

export default authSlice.reducer;
