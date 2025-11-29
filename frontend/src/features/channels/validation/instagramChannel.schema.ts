import * as Yup from 'yup';

// Instagram Channel Form Validation Schema
export const instagramChannelSchema = Yup.object().shape({
  channelName: Yup.string()
    .min(3, 'Channel name must be at least 3 characters')
    .max(50, 'Channel name must not exceed 50 characters')
    .required('Channel name is required'),
  
  email: Yup.string()
    .email('Invalid email address')
    .required('Email is required'),
  
  website: Yup.string()
    .url('Must be a valid URL (include http:// or https://)')
    .optional(),
});

// Update Channel Validation Schema
export const updateInstagramChannelSchema = Yup.object().shape({
  channelName: Yup.string()
    .min(3, 'Channel name must be at least 3 characters')
    .max(50, 'Channel name must not exceed 50 characters')
    .optional(),
  
  email: Yup.string()
    .email('Invalid email address')
    .optional(),
  
  website: Yup.string()
    .url('Must be a valid URL (include http:// or https://)')
    .optional(),
});
