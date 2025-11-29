import { combineReducers } from "@reduxjs/toolkit";
import authReducer from "@/features/auth/authSlice";
import instagramChannelReducer from "@/features/channels/store/instagramChannelSlice";

const rootReducer = combineReducers({
  auth: authReducer,
  instagramChannel: instagramChannelReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
export default rootReducer;
