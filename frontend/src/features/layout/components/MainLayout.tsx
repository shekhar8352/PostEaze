import { AppShell } from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { Outlet } from 'react-router-dom';
import { MainSidebar } from './MainSidebar';
import { AppHeader } from './AppHeader';
import layoutStyles from './MainLayout.module.css';

export const MainLayout = () => {
    const [mobileOpened, { toggle: toggleMobile }] = useDisclosure();
    const [desktopOpened, { toggle: toggleDesktop }] = useDisclosure(true);

    return (
        <AppShell
            header={{ height: 60 }}
            navbar={{
                width: 280,
                breakpoint: 'sm',
                collapsed: { mobile: !mobileOpened, desktop: !desktopOpened },
            }}
            padding="md"
            classNames={{ main: layoutStyles.main }}
        >
            <AppShell.Header>
                <AppHeader
                    mobileOpened={mobileOpened}
                    toggleMobile={toggleMobile}
                    desktopOpened={desktopOpened}
                    toggleDesktop={toggleDesktop}
                />
            </AppShell.Header>

            <AppShell.Navbar>
                <MainSidebar />
            </AppShell.Navbar>

            <AppShell.Main>
                <Outlet />
            </AppShell.Main>
        </AppShell>
    );
};
