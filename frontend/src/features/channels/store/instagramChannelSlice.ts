import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface InstagramChannelState {
  selectedChannelId: string | null;
  isCreateModalOpen: boolean;
  isEditModalOpen: boolean;
  filters: {
    status: 'all' | 'active' | 'inactive';
    searchQuery: string;
  };
}

const initialState: InstagramChannelState = {
  selectedChannelId: null,
  isCreateModalOpen: false,
  isEditModalOpen: false,
  filters: {
    status: 'all',
    searchQuery: '',
  },
};

const instagramChannelSlice = createSlice({
  name: 'instagramChannel',
  initialState,
  reducers: {
    // Modal Actions
    openCreateModal: (state) => {
      state.isCreateModalOpen = true;
    },
    closeCreateModal: (state) => {
      state.isCreateModalOpen = false;
    },
    openEditModal: (state, action: PayloadAction<string>) => {
      state.isEditModalOpen = true;
      state.selectedChannelId = action.payload;
    },
    closeEditModal: (state) => {
      state.isEditModalOpen = false;
      state.selectedChannelId = null;
    },
    
    // Selection Actions
    setSelectedChannel: (state, action: PayloadAction<string | null>) => {
      state.selectedChannelId = action.payload;
    },
    clearSelectedChannel: (state) => {
      state.selectedChannelId = null;
    },
    
    // Filter Actions
    setStatusFilter: (state, action: PayloadAction<'all' | 'active' | 'inactive'>) => {
      state.filters.status = action.payload;
    },
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.filters.searchQuery = action.payload;
    },
    clearFilters: (state) => {
      state.filters = initialState.filters;
    },
    
    // Reset State
    resetInstagramChannelState: () => initialState,
  },
});

export const {
  openCreateModal,
  closeCreateModal,
  openEditModal,
  closeEditModal,
  setSelectedChannel,
  clearSelectedChannel,
  setStatusFilter,
  setSearchQuery,
  clearFilters,
  resetInstagramChannelState,
} = instagramChannelSlice.actions;

export default instagramChannelSlice.reducer;
