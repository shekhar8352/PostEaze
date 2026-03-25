import type { ReactNode } from 'react';
import styles from './AuthPageShell.module.css';

export type AuthPageVariant = 'login' | 'register';

const copy: Record<
  AuthPageVariant,
  { kicker: string; headline: string; subline: string; meta: string }
> = {
  login: {
    kicker: 'Session',
    headline: 'Pick up where you left off.',
    subline:
      'One calm place to line up posts, sync channels, and ship on time—without the tab chaos.',
    meta: 'Secure session · Email verification keeps accounts real.',
  },
  register: {
    kicker: 'Onboarding',
    headline: 'Claim your publishing lane.',
    subline:
      'Create your workspace, verify your email once, and start scheduling with clarity.',
    meta: 'We never post without you · OAuth optional.',
  },
};

type AuthPageShellProps = {
  variant: AuthPageVariant;
  children: ReactNode;
};

export function AuthPageShell({ variant, children }: AuthPageShellProps) {
  const c = copy[variant];

  return (
    <div className={styles.root}>
      <div className={styles.backdrop} aria-hidden />
      <div className={styles.grain} aria-hidden />
      <div className={styles.grid}>
        <aside className={styles.brand}>
          <div>
            <div className={styles.wordmark}>
              Post<span>Eaze</span>
            </div>
            <p className={styles.kicker}>{c.kicker}</p>
            <h1 className={styles.headline}>{c.headline}</h1>
            <p className={styles.subline}>{c.subline}</p>
            <div className={styles.stripe} aria-hidden />
          </div>
          <p className={styles.meta}>
            <strong>PostEaze</strong> — schedule social content with intent, not noise.
          </p>
        </aside>
        <main className={styles.main}>
          <div className={styles.card}>{children}</div>
        </main>
      </div>
    </div>
  );
}
