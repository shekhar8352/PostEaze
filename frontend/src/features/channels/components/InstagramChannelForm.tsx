import { useState } from 'react';
import { Formik, Form, Field } from 'formik';
import {
    TextInput,
    Button,
    Stack,
    Alert,
    Box,
    Divider,
    Text,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { Icons } from '@/app/theme';
import { useInstagramOAuth } from '../hooks/useInstagramOAuth';
import { useCreateInstagramChannel } from '../services/instagramChannelQueries';
import { instagramChannelSchema } from '../validation/instagramChannel.schema';
import type { InstagramChannelFormData } from '../types/instagram.types';

interface InstagramChannelFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
}

export const InstagramChannelForm = ({
    onSuccess,
    onCancel,
}: InstagramChannelFormProps) => {
    const [authCode, setAuthCode] = useState<string>('');
    const { openOAuthPopup, isLoading: isOAuthLoading } = useInstagramOAuth();
    const createChannel = useCreateInstagramChannel();

    const initialValues: InstagramChannelFormData = {
        channelName: '',
        email: '',
        website: '',
    };

    const handleOAuthConnect = async () => {
        try {
            const code = await openOAuthPopup();
            setAuthCode(code);

            notifications.show({
                title: 'Instagram Connected!',
                message: 'Successfully connected to Instagram',
                color: 'green',
                icon: <Icons.CheckCircle size={18} />,
            });
        } catch (error) {
            notifications.show({
                title: 'Connection Failed',
                message: error instanceof Error ? error.message : 'Failed to connect to Instagram',
                color: 'red',
                icon: <Icons.XCircle size={18} />,
            });
        }
    };

    const handleSubmit = async (values: InstagramChannelFormData) => {
        if (!authCode) {
            notifications.show({
                title: 'Instagram Not Connected',
                message: 'Please connect your Instagram account first',
                color: 'orange',
                icon: <Icons.AlertCircle size={18} />,
            });
            return;
        }

        try {
            await createChannel.mutateAsync({
                channelName: values.channelName,
                email: values.email,
                website: values.website || undefined,
                authCode,
            });

            notifications.show({
                title: 'Channel Created!',
                message: 'Instagram channel has been successfully created',
                color: 'green',
                icon: <Icons.CheckCircle size={18} />,
            });

            onSuccess?.();
        } catch (error) {
            notifications.show({
                title: 'Creation Failed',
                message: error instanceof Error ? error.message : 'Failed to create channel',
                color: 'red',
                icon: <Icons.XCircle size={18} />,
            });
        }
    };

    return (
        <Formik
            initialValues={initialValues}
            validationSchema={instagramChannelSchema}
            onSubmit={handleSubmit}
        >
            {({ errors, touched, isSubmitting, setFieldValue, setFieldTouched }) => (
                <Form>
                    <Stack gap="md">
                        {/* Channel Details Section */}
                        <div>
                            <Text size="sm" fw={600} mb="xs">
                                Channel Details
                            </Text>
                            <Stack gap="sm">
                                <Box>
                                    <Field name="channelName">
                                        {({ field }: any) => (
                                            <TextInput
                                                {...field}
                                                label="Channel Name"
                                                placeholder="My Instagram Channel"
                                                required
                                                error={touched.channelName && errors.channelName ? errors.channelName : null}
                                                onChange={(e) => {
                                                    setFieldValue('channelName', e.target.value);
                                                    setFieldTouched('channelName', true);
                                                }}
                                            />
                                        )}
                                    </Field>
                                </Box>

                                <Box>
                                    <Field name="email">
                                        {({ field }: any) => (
                                            <TextInput
                                                {...field}
                                                label="Email"
                                                placeholder="email@example.com"
                                                type="email"
                                                required
                                                error={touched.email && errors.email ? errors.email : null}
                                                onChange={(e) => {
                                                    setFieldValue('email', e.target.value);
                                                    setFieldTouched('email', true);
                                                }}
                                            />
                                        )}
                                    </Field>
                                </Box>

                                <Box>
                                    <Field name="website">
                                        {({ field }: any) => (
                                            <TextInput
                                                {...field}
                                                label="Website (Optional)"
                                                placeholder="https://example.com"
                                                error={touched.website && errors.website ? errors.website : null}
                                                onChange={(e) => {
                                                    setFieldValue('website', e.target.value);
                                                    setFieldTouched('website', true);
                                                }}
                                            />
                                        )}
                                    </Field>
                                </Box>
                            </Stack>
                        </div>

                        <Divider />

                        {/* Instagram Connection Section */}
                        <div>
                            <Text size="sm" fw={600} mb="xs">
                                Instagram Authorization
                            </Text>

                            {authCode ? (
                                <Alert
                                    icon={<Icons.CheckCircle size={18} />}
                                    color="green"
                                    variant="light"
                                >
                                    Instagram account connected successfully!
                                </Alert>
                            ) : (
                                <Alert
                                    icon={<Icons.InfoCircle size={18} />}
                                    color="blue"
                                    variant="light"
                                >
                                    Click the button below to connect your Instagram account
                                </Alert>
                            )}

                            <Button
                                leftSection={<Icons.Instagram size={20} />}
                                fullWidth
                                mt="sm"
                                onClick={handleOAuthConnect}
                                loading={isOAuthLoading}
                                disabled={!!authCode}
                                style={{
                                    background: authCode
                                        ? undefined
                                        : 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
                                }}
                                variant={authCode ? 'light' : 'filled'}
                                color={authCode ? 'green' : undefined}
                            >
                                {authCode ? 'Connected to Instagram' : 'Connect Instagram Account'}
                            </Button>
                        </div>

                        {/* Form Actions */}
                        <Stack gap="sm" mt="md">
                            <Button
                                type="submit"
                                leftSection={<Icons.Plus size={18} />}
                                loading={createChannel.isPending || isSubmitting}
                                disabled={!authCode}
                                fullWidth
                            >
                                Create Channel
                            </Button>

                            {onCancel && (
                                <Button
                                    variant="subtle"
                                    onClick={onCancel}
                                    disabled={createChannel.isPending || isSubmitting}
                                    fullWidth
                                >
                                    Cancel
                                </Button>
                            )}
                        </Stack>
                    </Stack>
                </Form>
            )}
        </Formik>
    );
};
